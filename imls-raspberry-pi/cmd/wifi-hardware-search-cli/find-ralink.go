package main

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/jedib0t/go-pretty/v6/text"
	"gsa.gov/18f/find-ralink/config"
	"gsa.gov/18f/find-ralink/constants"
	"gsa.gov/18f/find-ralink/lshw"
	"gsa.gov/18f/find-ralink/models"
	"gsa.gov/18f/version"
)

func findMatchingDevice(wlan *models.Device) {
	// GetDeviceHash calls out to `lshw`.
	devices := lshw.GetDeviceHash(wlan)

	// We start by assuming that we have not found the device.
	found := false

	// Now, go through the devices and find the one that matches our criteria.
	for _, hash := range devices {

		if *config.Verbose {
			fmt.Println("---------")
			for k, v := range hash {
				fmt.Println(k, "<-", v)
			}
		}

		// The default is to search all the fields
		if wlan.Search.Field == "ALL" {

			for k := range hash {
				// Lowercase everything for purposes of pattern matching.
				v, _ := regexp.MatchString(strings.ToLower(wlan.Search.Query), strings.ToLower(hash[k]))
				if *config.Verbose {
					fmt.Println("query", wlan.Search.Query, "field", wlan.Search.Field)
				}
				if v {
					// If we find it, set the fact that it exists. This will be picked up
					// back out in main() for the final act of printing a message to the user.
					wlan.Exists = true
				}
				// Stop searching if we find something.
				if wlan.Exists {
					break
				}
			}
		} else {
			// If we aren't doing a full search, then this is the alternative: check just
			// one field. It will still be a lowercase search, but it will be against one field only.
			if *config.Verbose {
				fmt.Println("query", wlan.Search.Query, "field", wlan.Search.Field)
			}
			v, _ := regexp.MatchString(strings.ToLower(wlan.Search.Query), strings.ToLower(hash[wlan.Search.Field]))
			if v {
				wlan.Exists = true
			}
		}

		// If we find something, proceed. But only keep the first thing we find.
		// Back in 'main', we'll handle the case where wlan.exists is false.
		if wlan.Exists && !found {
			found = true
			wlan.Vendor = strings.ToLower(hash["vendor"])
			wlan.Physicalid, _ = strconv.Atoi(hash["physical id"])
			wlan.Description = strings.ToLower(hash["description"])
			wlan.Businfo = strings.ToLower(hash["bus info"])
			wlan.Logicalname = strings.ToLower(hash["logical name"])
			wlan.Serial = strings.ToLower(hash["serial"])

			if len(hash["serial"]) >= constants.MACLENGTH {
				wlan.Mac = strings.ToLower(hash["serial"][0:constants.MACLENGTH])
			} else {
				wlan.Mac = strings.ToLower(hash["serial"])
			}
			wlan.Configuration = strings.ToLower(hash["configuration"])

			// Once we populate something in this loop, break out.
			// This will return us to the caller with a populated structure.
			break
		}
	}
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
	lshwPtr := flag.String("lshw-path", config.LSHW_EXE, "Path to the `lshw` binary.")
	searchesPtr := flag.String("searchjson-path", config.SEARCHES_PATH, "Path to a JSON file containing default searches.")
	// In case we care about versioning.
	versionPtr := flag.Bool("version", false, "Get the software version and exit.")
	flag.Parse()

	// Using a "global" indicator of verboseness.
	// Not sure if there is a more Go-ish way to do this.
	config.Verbose = verbosePtr

	// Override configuration if needed, in case things move/change names.
	config.LSHW_EXE = *lshwPtr
	config.SEARCHES_PATH = *searchesPtr

	// ROOT
	// We can't do this without root. Some things... might work.
	// Print a big red warning.
	if os.Getenv("USER") != "root" {
		fmt.Println(text.FgRed.Sprint("find-ralink *really* needs to be run as root."))
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
		for _, s := range config.GetSearches() {
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
		s := &config.Search{Field: *fieldPtr, Query: *searchPtr}
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
