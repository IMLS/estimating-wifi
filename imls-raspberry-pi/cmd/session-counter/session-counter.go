package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/robfig/cron/v3"
	"gsa.gov/18f/cmd/session-counter/tlp"
	"gsa.gov/18f/internal/state"
	"gsa.gov/18f/internal/structs"
	"gsa.gov/18f/internal/version"
	"gsa.gov/18f/internal/wifi-hardware-search/models"
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

	state.NewConfigFromPath(*configPathPtr)

}

func runEvery(crontab string, kb *tlp.KillBroker, fun func() bool) {
	cfg := state.GetConfig()
	c := cron.New()
	_, err := c.AddFunc(crontab, func() {
		if !fun() {
			cfg.Log().Error(fun, " returned false...")
		}
	})
	if err != nil {
		cfg.Log().Error("cron: could not set up crontab entry")
		cfg.Log().Fatal(err.Error())
	}
	c.Start()

	ch := kb.Subscribe()
	go kb.Start()
	for {
		<-ch
		c.Stop()
		return
	}
}

func fakeShark2(dev string) []string {
	return []string{"DE:AD:BE:EF:00:00", "BE:EF:00:00:00:00"}
}

func run2() {
	cfg := state.GetConfig()
	// The MAC database MUST be ephemeral. Put it in RAM.
	mac_db := state.NewSqliteDB(":memory:")
	mac_db.CreateTableFromStruct(structs.EphemeralDuration{})

	kb := tlp.NewKillBroker()
	go runEvery("*/1 * * * *", kb,
		func() bool {
			cfg.Log().Debug("RUNNING SIMPLESHARK")
			return tlp.SimpleShark(mac_db,
				// search.SetMonitorMode,
				func(d *models.Device) {},
				// search.SearchForMatchingDevice,
				func() *models.Device { return &models.Device{Exists: true, Logicalname: "fakewan0"} },
				// tlp.TSharkRunner
				fakeShark2)
		})

	sq := state.NewQueue("sent")
	iq := state.NewQueue("images")
	go runEvery("*/2 * * * *", kb,
		func() bool {
			cfg.Log().Debug("RUNNING PROCESSDATA")
			// Copy ephemeral durations over to the durations table
			tlp.ProcessData(mac_db, sq, iq)
			// Clear out the ephemeral data for the next day of monitoring
			mac_db.GetPtr().Exec("DELETE FROM ephemeraldurations")
			// Draw images of the data
			tlp.WriteImages(cfg.GetDurationsDatabase())
			return true
		})
}

func main() {
	// DO NOT USE LOGGING YET
	initConfigFromFlags()
	cfg := state.GetConfig()
	state.InitConfig()
	// NOW YOU MAY USE LOGGING.

	cfg.Log().Debug("session id at startup is ", cfg.GetCurrentSessionID())

	// Run the network
	var wg sync.WaitGroup
	wg.Add(1)
	go run2()
	// Wait forever.
	wg.Wait()
}
