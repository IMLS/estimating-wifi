package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"gsa.gov/18f/cmd/session-counter/tlp"
	"gsa.gov/18f/internal/state"
	"gsa.gov/18f/internal/version"
	"gsa.gov/18f/internal/wifi-hardware-search/search"
	zls "gsa.gov/18f/internal/zero-log-sentry"
)

var (
	cfgFile string
)

func runEvery(crontab string, c *cron.Cron, fun func()) {
	id, err := c.AddFunc(crontab, fun)
	log.Debug().
		Str("crontab", crontab).
		Str("id", fmt.Sprintf("%v", id)).
		Msg("runEvery")
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("cron: could not set up crontab entry")
	}
}

func run2() {
	sq := state.NewQueue[int64]("sent")
	durationsdb := state.GetDurationsDatabase()
	c := cron.New()

	go runEvery("*/1 * * * *", c,
		func() {
			log.Debug().Msg("RUNNING SIMPLESHARK")
			tlp.SimpleShark(
				search.SetMonitorMode,
				search.SearchForMatchingDevice,
				tlp.TSharkRunner)
		})

	go runEvery(state.GetResetCron(), c,
		func() {
			log.Info().
				Str("time", fmt.Sprintf("%v", state.GetClock().Now().In(time.Local))).
				Msg("RUNNING PROCESSDATA")
			// Copy ephemeral durations over to the durations table
<<<<<<< HEAD
			tlp.ProcessData(durationsdb, sq, iq)
=======
			tlp.ProcessData(durationsdb, sq)
			// Draw images of the data
			tlp.WriteImages(durationsdb)
>>>>>>> 994993fa (TLP test runs; generics fixed.)
			// Try sending the data
			tlp.SimpleSend(durationsdb)
			// Increment the session counter
			state.IncrementSessionID()
			// Clear out the ephemeral data for the next day of monitoring
			state.ClearEphemeralDB()
		})

	// Start the cron jobs...
	c.Start()
}

func launchTLP() {
	state.SetConfigAtPath(cfgFile)
	dsn := state.GetSentryDSN()
	if dsn != "" {
		zls.SetupZeroLogSentry("session-counter", dsn)
		zls.SetTags(map[string]string{
			"tag":     state.GetDeviceTag(),
			"fcfs_id": state.GetFCFSSeqID(),
		})
	}

	log.Info().
		Int64("session_id", state.GetCurrentSessionID()).
		Msg("session id at launch")

	// Run the network
	var wg sync.WaitGroup
	wg.Add(1)
	go run2()

	// Wait forever.
	wg.Wait()
}

var rootCmd = &cobra.Command{
	Use:   "session-counter",
	Short: "A tool for monitoring wifi devices while preserving privacy.",
	Long: `session-counter watches to see what wifi devices are nearby while
carefully leaving out information that would impose on the privacy of people
around you.`,
	Run: func(cmd *cobra.Command, args []string) {
		launchTLP()
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "session-counter version",
	Long:  `Print the version number of session-counter`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(version.GetVersion())
	},
}

func main() {
	rootCmd.PersistentFlags().StringVar(&cfgFile,
		"config",
		"session-counter.ini",
		"config file (default is session-counter.ini in /etc/imls, %PROGRAMDATA%\\IMLS, or current directory")
	rootCmd.AddCommand(versionCmd)
	rootCmd.Execute()
}
