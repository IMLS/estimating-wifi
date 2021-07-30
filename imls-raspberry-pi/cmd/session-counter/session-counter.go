package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
	"gsa.gov/18f/cmd/session-counter/tlp"
	"gsa.gov/18f/internal/state"
	"gsa.gov/18f/internal/version"
	"gsa.gov/18f/internal/wifi-hardware-search/search"
)

func initConfigFromFlags() {
	versionPtr := flag.Bool("version", false, "Get the software version and exit.")
	configPathPtr := flag.String("config", "", "Path to config.yaml. REQUIRED.")
	flag.Parse()

	// If they just want the version, print and exit.
	if *versionPtr {
		fmt.Println(version.GetVersion())
		os.Exit(0)
	}

	// Make sure a config is passed.
	if *configPathPtr == "" {
		log.Fatal("The flag --config MUST be provided.")
		os.Exit(1)
	}

	if _, err := os.Stat(*configPathPtr); os.IsNotExist(err) {
		log.Println("Looked for config at ", *configPathPtr)
		log.Fatal("Cannot find config file. Exiting.")
	}

	state.SetConfigAtPath(*configPathPtr)

}

func runEvery(crontab string, c *cron.Cron, fun func()) {
	cfg := state.GetConfig()
	id, err := c.AddFunc(crontab, fun)
	cfg.Log().Debug("launched crontab ", crontab, " with id ", id)
	if err != nil {
		cfg.Log().Error("cron: could not set up crontab entry")
		cfg.Log().Fatal(err.Error())
	}
}

func run2() {
	cfg := state.GetConfig()
	sq := state.NewQueue("sent")
	iq := state.NewQueue("images")
	durationsdb := cfg.GetDurationsDatabase()
	c := cron.New()

	go runEvery("*/1 * * * *", c,
		func() {
			cfg.Log().Debug("RUNNING SIMPLESHARK")
			tlp.SimpleShark(
				search.SetMonitorMode,
				search.SearchForMatchingDevice,
				tlp.TSharkRunner)
		})

	go runEvery(cfg.GetResetCron(), c,
		func() {
			cfg.Log().Info("RUNNING PROCESSDATA at ", state.GetClock().Now().In(time.Local))
			// Copy ephemeral durations over to the durations table
			tlp.ProcessData(durationsdb, sq, iq)
			// Draw images of the data
			tlp.WriteImages(durationsdb)
			// Try sending the data
			tlp.SimpleSend(durationsdb)
			// Increment the session counter
			cfg.IncrementSessionID()
			// Clear out the ephemeral data for the next day of monitoring
			state.ClearEphemeralDB()
		})

	// Start the cron jobs...
	c.Start()
}

func main() {
	// DO NOT USE LOGGING YET
	initConfigFromFlags()
	cfg := state.GetConfig()
	// NOW YOU MAY USE LOGGING.

	cfg.Log().Info("Startup session id ", cfg.GetCurrentSessionID())

	// Run the network
	var wg sync.WaitGroup
	wg.Add(1)
	go run2()
	// Wait forever.
	wg.Wait()
}
