package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"gsa.gov/18f/cmd/session-counter/tlp"
	"gsa.gov/18f/internal/state"
	"gsa.gov/18f/internal/version"
	"gsa.gov/18f/internal/wifi-hardware-search/search"
)

var (
	cfgFile  string
	logLevel string
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
	cfg := state.GetConfig()
	sq := state.NewQueue("sent")
	iq := state.NewQueue("images")
	durationsdb := cfg.GetDurationsDatabase()
	c := cron.New()

	go runEvery("*/1 * * * *", c,
		func() {
			log.Debug().Msg("RUNNING SIMPLESHARK")
			tlp.SimpleShark(
				search.SetMonitorMode,
				search.SearchForMatchingDevice,
				tlp.TSharkRunner)
		})

	go runEvery(cfg.GetResetCron(), c,
		func() {
			log.Info().
				Str("time", fmt.Sprintf("%v", state.GetClock().Now().In(time.Local))).
				Msg("RUNNING PROCESSDATA")
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

func launchTLP() {
	state.SetConfigAtPath(cfgFile)
	cfg := state.GetConfig()

	log.Info().
		Int64("session_id", cfg.GetCurrentSessionID()).
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

func initLogs() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	switch lvl := logLevel; lvl {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}

func main() {
	rootCmd.PersistentFlags().StringVar(&cfgFile,
		"config",
		"config.sqlite3",
		"config file (default is config.sqlite3 in current directory")
	rootCmd.PersistentFlags().StringVar(&logLevel,
		"logging",
		"info",
		"logging level (debug, info, warn, error)")
	cobra.OnInitialize(initLogs)
	rootCmd.AddCommand(versionCmd)
	rootCmd.Execute()
}
