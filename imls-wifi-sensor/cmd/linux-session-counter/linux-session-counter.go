package main

import (
	"fmt"
	"sync"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"gsa.gov/18f/internal/config"
	"gsa.gov/18f/internal/session-counter-helper/session_counter"
	"gsa.gov/18f/internal/session-counter-helper/state"
	zls "gsa.gov/18f/internal/session-counter-helper/zero-log-sentry"
	"gsa.gov/18f/internal/version"
)

var (
	cfgFile string
	mode    string = "prod"
)

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
	go session_counter.Run2(mode)

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
	//Sets default mode to "prod"
	//rootCmd.PersistentFlags().StringVar(&mode, "mode", "prod", "Mode to run the program in")
	//viper.BindPFlag("mode.run", rootCmd.PersistentFlags().Lookup("mode"))
	rootCmd.AddCommand(versionCmd)
	rootCmd.Execute()
}
