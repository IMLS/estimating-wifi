package main

import (
	"fmt"
	"sync"
	"time"

	cron "github.com/robfig/cron/v3"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"gsa.gov/18f/cmd/session-counter/state"
	"gsa.gov/18f/cmd/session-counter/tlp"
	zls "gsa.gov/18f/cmd/session-counter/zero-log-sentry"
	"gsa.gov/18f/internal/config"
	"gsa.gov/18f/internal/version"
	"gsa.gov/18f/internal/wifi-hardware-search/search"

	_ "net/http/pprof"
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
	durationsdb := state.NewDurationsDB()
	c := cron.New()

	go runEvery("*/1 * * * *", c,
		func() {
			log.Debug().Msg("RUNNING SIMPLESHARK")
			tlp.SimpleShark(
				search.SetMonitorMode,
				search.SearchForMatchingDevice,
				tlp.TSharkRunner)
		})

	go runEvery(config.GetResetCron(), c,
		func() {
			log.Info().
				Str("time", fmt.Sprintf("%v", state.GetClock().Now().In(time.Local))).
				Msg("RUNNING PROCESSDATA")
			// Copy ephemeral durations over to the durations table
			tlp.ProcessData(durationsdb, sq)
			// Draw images of the data
			// tlp.WriteImages(durationsdb)
			// Try sending the data
			tlp.SimpleSend(durationsdb)
			// Increment the session counter
			state.IncrementSessionID()
			// Clear out the ephemeral data for the next day of monitoring
			state.ClearEphemeralDB()
			durationsdb.ClearDurationsDB()

		})

	// Start the cron jobs...
	c.Start()
}

func launchTLP() {
	// if viper.GetBool("WITH_PROFILE") {
	// 	go http.ListenAndServe("localhost:8080", nil)
	// 	log.Info().
	// 		Str("time", fmt.Sprintf("%v", state.GetClock().Now().In(time.Local))).
	// 		Msg("Launching pprof server server")
	// }

	config.SetConfigAtPath(cfgFile)
	dsn := config.GetSentryDSN()
	if dsn != "" {
		zls.SetupZeroLogSentry("session-counter", dsn)
		zls.SetTags(map[string]string{
			"tag":     config.GetDeviceTag(),
			"fscs_id": config.GetFSCSID(),
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
