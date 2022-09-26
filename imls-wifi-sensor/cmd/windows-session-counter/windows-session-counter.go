package main

import (
	"sync"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/toast.v1"
	"gsa.gov/18f/internal/config"
	"gsa.gov/18f/internal/session-counter-helper/session_counter"
	"gsa.gov/18f/internal/session-counter-helper/state"
	zls "gsa.gov/18f/internal/session-counter-helper/zero-log-sentry"
)

var (
	cfgFile string
	mode    string
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

	toastNotifSuccessfulInstall()
	// Run the network
	var wg sync.WaitGroup
	wg.Add(1)
	go session_counter.Run2()
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

func toastNotifSuccessfulInstall() {
	notification := toast.Notification{
		AppID:   "estimating-wifi",
		Title:   "session-counter",
		Message: "IMLS session counter is running in the background. Open Task Manager using Ctrl + Alt + Delete to make sure it's running.",
	}
	err := notification.Push()
	if err != nil {
		log.Fatal()
	}
}

func main() {
	rootCmd.PersistentFlags().StringVar(&cfgFile,
		"config",
		"session-counter.ini",
		"config file (default is session-counter.ini in /etc/imls, %PROGRAMDATA%\\IMLS, or current directory")
	rootCmd.PersistentFlags().StringVar(&mode, "mode", "prod", "Mode to run the program in")
	viper.BindPFlag("mode.run", rootCmd.PersistentFlags().Lookup("mode"))
	rootCmd.Execute()
}
