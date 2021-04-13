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
)

func findMatchingDevice(wlan *models.Device) {

	devices := lshw.GetDeviceHash(wlan)
	found := false

	// Now, go through the devices and find the one that matches our criteria.
	// Either we're looking for a vendor or for something in the config string.
	// A more general search could be implemented, but I'll keep it limited to prevent
	// spurious matches that were unexpected/surprising in practice. (I hope.)
	for _, hash := range devices {

		if config.Verbose {
			fmt.Println("---------")
			for k, v := range hash {
				fmt.Println(k, "<-", v)
			}
		}

		// The default is to search all the fields
		// lowercase pattern matching
		if wlan.Search.Field == "ALL" {
			for k := range hash {
				v, _ := regexp.MatchString(strings.ToLower(wlan.Search.Query), strings.ToLower(hash[k]))
				if v {
					wlan.Exists = true
				}
				// Stop searching if we find something.
				if wlan.Exists {
					break
				}
			}
		} else {
			// Otherwise, search the field specified.
			if config.Verbose {
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
func getField(v *models.Device, field string) reflect.Value {
	r := reflect.ValueOf(v)
	f := reflect.Indirect(r).FieldByName(field)
	return f
}

func main() {
	// FLAGS
	// This thing has a whole mess of flags.
	// The default values help us make sure that sensible things happen if options
	// are not explicitly declared.
	verbosePtr := flag.Bool("verbose", false, "Verbose output.")
	discoverPtr := flag.Bool("discover", false, "Attempt to discover the device.")
	searchPtr := flag.String("search", "ralink", "Search string to use in hardware listing. Must use with `field`.")
	fieldPtr := flag.String("field", "ALL", "Field to search.")
	extractPtr := flag.String("extract", "", "Field to extract from device data.")
	existsPtr := flag.Bool("exists", false, "Ask if a device exists. Returns `true` or `false`.")
	versionPtr := flag.Bool("version", false, "Get the software version and exit.")
	flag.Parse()

	// Using a "global" indicator of verboseness.
	// Not sure if there is a more Go-ish way to do this.
	config.Verbose = *verbosePtr

	// VERSION
	// If they just want the version, print it and exit.
	if *versionPtr {
		fmt.Println("Version", constants.VERSION)
		os.Exit(0)
	}

	// ROOT
	// We can't do this without root. Some things... might work.
	// Print a big red warning.
	if os.Getenv("USER") != "root" {
		fmt.Println(text.FgRed.Sprint("find-ralink *really* needs to be run as root."))
	}

	// EXTRACT
	// If they didn't tell us what to extract, then
	// extract the field they're searching for. That is, we use the
	// extractPtr value to tell us what from the device description to return
	// at the end of the seach.
	if *extractPtr == "" {
		// If it is the default "ALL" value, choose the logicalname
		if *fieldPtr == "ALL" {
			*extractPtr = "logicalname"
		} else {
			*extractPtr = *fieldPtr
		}
	}

	// We populate this via calls.
	device := new(models.Device)
	device.Extract = *extractPtr

	// If we ask for auto-discovery, try it and exit.
	if *discoverPtr {
		for _, s := range config.GetSearches() {
			device.Search = &s
			findMatchingDevice(device)
			if device.Exists {
				break
			}
		}
	} else {
		s := &config.Search{Field: *fieldPtr, Query: *searchPtr}
		device.Search = s
		findMatchingDevice(device)
	}

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
