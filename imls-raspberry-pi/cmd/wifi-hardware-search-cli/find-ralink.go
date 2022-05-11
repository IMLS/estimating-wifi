package main

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/jedib0t/go-pretty/v6/text"
	"gsa.gov/18f/internal/version"
	"gsa.gov/18f/internal/wifi-hardware-search/config"
	"gsa.gov/18f/internal/wifi-hardware-search/models"
	"gsa.gov/18f/internal/wifi-hardware-search/search"
)

func findMatchingDevice(wlan *models.Device) {
	search.FindMatchingDevice(wlan)
}

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

func main() {
	// FLAGS
	// This thing has a whole mess of flags.
	// The default values help us make sure that sensible things happen if options
	// are not explicitly declared.
	verbosePtr := flag.Bool("verbose", false, "Verbose output.")
	// The default behavior is for the tool to --discover a compatible adapter.
	discoverPtr := flag.Bool("discover", true, "Attempt to discover the device. Default.")
	// If either --search or --field are used, we disable --discover as a mode.
	// We don't require both to be set, but do provide sensible defaults in the event that the user
	// forgets. This may not be preferable to throwing an error, but the goal is to make sure the tool
	// runs, even if it doesn't work. We'd rather not have the ansible script crashing.
	searchPtr := flag.String("search", "", "Search string to use in hardware listing. Must use with `field`.")
	fieldPtr := flag.String("field", "", "Field to search.")
	// By default, we're always looking for the logical name of the device. But, we might want other fields.
	extractPtr := flag.String("extract", "logical name", "Field to extract from device data.")
	// Instead of a exit(-1), we might want to have this print true/false, which is important for
	// ansible integration.
	existsPtr := flag.Bool("exists", false, "Ask if a device exists. Returns `true` or `false`.")
	// It is possible, but unlikely, that the location of lshw will change.
	lshwPtr := flag.String("lshw-path", config.GetLSHWLocation(), "Path to the `lshw` binary.")
	// In case we care about versioning.
	versionPtr := flag.Bool("version", false, "Get the software version and exit.")
	flag.Parse()

	// Using a "global" indicator of verboseness.
	// Not sure if there is a more Go-ish way to do this.
	config.Verbose = verbosePtr

	// Override configuration if needed, in case things move/change names.
	config.SetLSHWLocation(*lshwPtr)

	// ROOT
	// We can't do this without root. Some things... might work.
	// Print a big red warning.
	if os.Getuid() != 0 {
		fmt.Println(text.FgRed.Sprint("wifi-hardware-search-cli *really* needs to be run as root."))
	}

	// VERSION
	// If they just want the version, print it and exit.
	if *versionPtr {
		fmt.Println(version.GetVersion())
		os.Exit(0)
	}

	// We populate this via calls.
	device := new(models.Device)
	device.Extract = *extractPtr

	// If either --field or --search are used, then we need to do two things
	//  1. disable --discovery
	//  2. make sure there are sensible defaults for both flags, because they operate together.
	if *fieldPtr != "" || *searchPtr != "" {
		*discoverPtr = false
		if *fieldPtr == "" {
			*fieldPtr = "ALL"
		}
		if *searchPtr == "" {
			*searchPtr = "ralink"
		}
	}

	// If we ask for auto-discovery, try it and exit.
	if *discoverPtr {
		// We either have searches in /etc, or a few held in reserve
		// that are compiled in via `embed`. GetSearches pulls the "live"
		// searches if it can, and falls back to the embedded if needed.
		// It goes through each one-by-one.
		for _, s := range search.GetSearches() {
			if *config.Verbose {
				fmt.Println("search", s)
			}
			device.Search = &s
			// findMatchingDevice populates device.Exists if something is found.
			findMatchingDevice(device)
			// We can stop searching if we find something.
			if device.Exists {
				break
			}
		}
	} else {
		// The alternative is we're not doing an exhaustive/discovery search.
		// Therefore, we should use the field/search ptrs.
		s := &models.Search{Field: *fieldPtr, Query: *searchPtr}
		device.Search = s
		findMatchingDevice(device)
	}

	// If we asked for a true/false value, print that.
	if *existsPtr {
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
		res := getField(device, *extractPtr)
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
