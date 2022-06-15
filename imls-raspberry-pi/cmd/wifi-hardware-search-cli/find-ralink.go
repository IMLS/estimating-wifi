package main

import (
	"os"
	"reflect"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"gsa.gov/18f/internal/state"
	"gsa.gov/18f/internal/version"
	"gsa.gov/18f/internal/wifi-hardware-search/models"
	"gsa.gov/18f/internal/wifi-hardware-search/search"
	"gsa.gov/18f/internal/zero-log-sentry"
)

var (
	cfgFile    string
	discover   bool
	lshwSearch string
	lshwField  string
	extract    string
	exists     bool
)

// https://stackoverflow.com/questions/18930910/access-struct-property-by-name
// PURPOSE
// This reflects on the Device structure and attempts to pull values out by
// name (passed in as a string). The alternative would be some kind of case statement,
// I think, where we pattern match on what is passed in, and extract the correct
// field as a result. The case statement might be safer, but for now, we'll do this
// fancy thing. If it becomes a problem in practice, we can always do if/else if/else...
func getField(v *models.Device, field string) reflect.Value {
	r := reflect.ValueOf(v)
	// Replace all strings, and titlecase the word. This matches
	// the way the fields are named in the data structure.
	adjusted := strings.Title(strings.ReplaceAll(field, " ", ""))
	f := reflect.Indirect(r).FieldByName(adjusted)
	return f
}

func launch() {
	if os.Getuid() != 0 {
		log.Warn().Msg("wifi-hardware-search-cli needs to be run as root.")
	}

	device := new(models.Device)
	device.Extract = extract

	// If either --field or --search are used, then we need to do two things
	//  1. disable --discovery
	//  2. make sure there are sensible defaults for both flags, because they operate together.
	if lshwField != "" || lshwSearch != "" {
		discover = false
		if lshwField == "" {
			lshwField = "ALL"
		}
		if lshwSearch == "" {
			lshwSearch = "ralink"
		}
	}

	// If we ask for auto-discovery, try it and exit.
	if discover {
		// We either have searches in /etc, or a few held in reserve
		// that are compiled in via `embed`. GetSearches pulls the "live"
		// searches if it can, and falls back to the embedded if needed.
		// It goes through each one-by-one.
		for _, s := range search.GetSearches() {
			log.Debug().
				Str("field", s.Field).
				Str("query", s.Query).
				Msg("searching")
			device.Search = &s
			// FindMatchingDevice populates device.Exists if something is found.
			search.FindMatchingDevice(device)
			// We can stop searching if we find something.
			if device.Exists {
				break
			}
		}
	} else {
		// The alternative is we're not doing an exhaustive/discovery search.
		// Therefore, we should use the field/search params
		s := &models.Search{Field: lshwField, Query: lshwSearch}
		device.Search = s
		search.FindMatchingDevice(device)
	}

	// If we asked for a true/false value, print that.
	if exists {
		// If we're explicitly asking to see if it exists, say no, and
		// return a zero error code.
		if device.Exists {
			log.Info().Msg("device exists")
			os.Exit(0)
		} else {
			log.Info().Msg("device does not exist")
			os.Exit(1)
		}
	} else if device.Exists {
		// Otherwise, we're going to use reflection to look into the device structure
		// and try and extract the field they asked for. If they named the field incorrectly,
		// this will fail in a pretty gruesome way.
		res := getField(device, extract)
		log.Info().Str("result", res.String()).Msg("obtained field")
		os.Exit(0)
	} else {
		// Otherwise... this is bad. Exit with error.
		log.Fatal().Msg("Device not found")
		os.Exit(-1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "wifi-hardware-search",
	Short: "A helper tool for finding wifi usb adapters.",
	Long:  `wifi-hardware-search returns the first supported wifi usb adapter.`,
	Run: func(cmd *cobra.Command, args []string) {
		launch()
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "wifi-hardware-search version",
	Long:  `Print the version number of wifi-hardware-search`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info().Str("version", version.GetVersion()).Msg("fin.")
	},
}

func initialize() {
	state.SetConfigAtPath(cfgFile)
	dsn := state.GetSentryDSN()
	if dsn != "" {
		zls.SetupZeroLogSentry("wifi-hardware-search-cli", dsn)
	}
}

func main() {
	rootCmd.PersistentFlags().StringVar(&cfgFile,
		"config",
		"session-counter.ini",
		"config file (default is session-counter.ini in /etc/imls, %PROGRAMDATA%\\IMLS, or current directory")
	rootCmd.PersistentFlags().BoolVar(&discover,
		"discover",
		true,
		"attempt to discover the device")
	rootCmd.PersistentFlags().StringVar(&lshwSearch,
		"search",
		"",
		"search string to use in hardware listing. must be used with `field`")
	rootCmd.PersistentFlags().StringVar(&lshwField,
		"field",
		"",
		"field to search for")
	rootCmd.PersistentFlags().StringVar(&extract,
		"extract",
		"logical name",
		"field to extract from device data")
	rootCmd.PersistentFlags().BoolVar(&exists,
		"exists",
		false,
		"returns true or false if a device exists")

	cobra.OnInitialize(initialize)
	rootCmd.AddCommand(versionCmd)
	rootCmd.Execute()
}
