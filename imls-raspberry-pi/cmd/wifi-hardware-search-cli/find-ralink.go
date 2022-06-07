package main

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"gsa.gov/18f/internal/state"
	"gsa.gov/18f/internal/version"
	"gsa.gov/18f/internal/wifi-hardware-search/models"
	"gsa.gov/18f/internal/wifi-hardware-search/search"
)

var (
	cfgFile    string
	logLevel   string
	verbose    bool
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
		fmt.Println(text.FgRed.Sprint("wifi-hardware-search-cli *really* needs to be run as root."))
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
			if verbose {
				fmt.Println("search", s)
			}
			device.Search = &s
			// FindMatchingDevice populates device.Exists if something is found.
			search.FindMatchingDevice(device, verbose)
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
		search.FindMatchingDevice(device, verbose)
	}

	// If we asked for a true/false value, print that.
	if exists {
		// If we're explicitly asking to see if it exists, say no, and
		// return a zero error code.
		if device.Exists {
			fmt.Println("true")
		} else {
			fmt.Println("false")
		}
		os.Exit(0)
	} else if device.Exists {
		// Otherwise, we're going to use reflection to look into the device structure
		// and try and extract the field they asked for. If they named the field incorrectly,
		// this will fail in a pretty gruesome way.
		res := getField(device, extract)
		if reflect.TypeOf(res).Kind() == reflect.Bool {
			fmt.Println(res.Interface())
		} else {
			fmt.Println(res)
		}
		os.Exit(0)
	} else {
		// Otherwise... this is bad. Exit with error.
		fmt.Println("Device not found")
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
		fmt.Println(version.GetVersion())
	},
}

func initialize() {
	state.SetConfigAtPath(cfgFile)
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
		"session-counter.ini",
		"config file (default is session-counter.ini in /etc/imls, %PROGRAMDATA%\\IMLS, or current directory")
	rootCmd.PersistentFlags().StringVar(&logLevel,
		"logging",
		"info",
		"logging level (debug, info, warn, error)")
	rootCmd.PersistentFlags().BoolVar(&verbose,
		"verbose",
		false,
		"verbose output")
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
