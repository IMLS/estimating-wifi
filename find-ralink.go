package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/jedib0t/go-pretty/v6/text"
	"gsa.gov/18f/find-ralink/config"
	"gsa.gov/18f/find-ralink/constants"
)

const (
	LOOKING_FOR_USB = iota
	READING_HASH    = iota
	DONE_READING    = iota
)

type Device struct {
	exists        bool
	search        config.Search
	physicalid    int
	description   string
	businfo       string
	logicalname   string
	serial        string
	mac           string
	configuration string
	vendor        string
	extract       string
}

func getDeviceHash(wlan *Device) []map[string]string {
	wlan.exists = false

	cmd := exec.Command("/usr/bin/lshw", "-class", "network")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Println("cpw: cannot get stdout from lshw")
		log.Fatal(err)
	}

	if err := cmd.Start(); err != nil {
		log.Println("cpw: cannot start `lshw` command")
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(stdout)
	hash := make(map[string]string, 0)
	usbSecRe := regexp.MustCompile(`^\s+\*-(usb|network).*`)
	newSecRe := regexp.MustCompile(`^\s+\*-.*`)
	hashRe := regexp.MustCompile(`^\s+(.*?): (.*)`)
	state := LOOKING_FOR_USB

	// Build up an array of hashes. Instead of looking for the device here,
	// we'll instead collect all the devices into hashes, and hold them for a moment.
	devices := make([]map[string]string, 0)

	for scanner.Scan() {
		line := scanner.Text()
		switch state {
		case LOOKING_FOR_USB:
			match := usbSecRe.MatchString(line)
			if match {
				if config.Verbose {
					fmt.Println("-> READING_HASH")
				}
				// Create a new hash.
				hash = make(map[string]string)
				state = READING_HASH
			}
		case READING_HASH:
			if config.Verbose {
				fmt.Printf("checking: [ %v ]\n", line)
			}
			newSecMatch := newSecRe.MatchString(line)
			hashMatch := hashRe.MatchString(line)
			hashPieces := hashRe.FindStringSubmatch(line)

			if newSecMatch {
				if config.Verbose {
					fmt.Println("-> DONE_READING")
				}
				state = DONE_READING
			} else if hashMatch {
				// fmt.Printf("%v <- %v\n", hashPieces[1], hashPieces[2])
				// 0 is the full string, 1 the first group, 2 the second.
				hash[hashPieces[1]] = hashPieces[2]
			}
		case DONE_READING:
			state = LOOKING_FOR_USB
			devices = append(devices, hash)
			if config.Verbose {
				fmt.Println("devices len", len(devices))
			}
		}
	}

	// Don't lose the last hash!
	devices = append(devices, hash)

	return devices
}

func findMatchingDevice(wlan *Device) {

	devices := getDeviceHash(wlan)
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
		if wlan.search.Field == "ALL" {
			for k := range hash {
				v, _ := regexp.MatchString(strings.ToLower(wlan.search.Query), strings.ToLower(hash[k]))
				if v {
					wlan.exists = true
				}
				// Stop searching if we find something.
				if wlan.exists {
					break
				}
			}
		} else {
			// Otherwise, search the field specified.
			if config.Verbose {
				fmt.Println("query", wlan.search.Query, "field", wlan.search.Field)
			}
			v, _ := regexp.MatchString(strings.ToLower(wlan.search.Query), strings.ToLower(hash[wlan.search.Field]))
			if v {
				wlan.exists = true
			}
		}

		// If we find something, proceed. But only keep the first thing we find.
		// Back in 'main', we'll handle the case where wlan.exists is false.
		if wlan.exists && !found {
			found = true
			wlan.vendor = strings.ToLower(hash["vendor"])
			wlan.physicalid, _ = strconv.Atoi(hash["physical id"])
			wlan.description = strings.ToLower(hash["description"])
			wlan.businfo = strings.ToLower(hash["bus info"])
			wlan.logicalname = strings.ToLower(hash["logical name"])
			wlan.serial = strings.ToLower(hash["serial"])

			if len(hash["serial"]) >= constants.MACLENGTH {
				wlan.mac = strings.ToLower(hash["serial"][0:constants.MACLENGTH])
			} else {
				wlan.mac = strings.ToLower(hash["serial"])
			}
			wlan.configuration = strings.ToLower(hash["configuration"])

			// Once we populate something in this loop, break out.
			// This will return us to the caller with a populated structure.
			break
		}
	}
}

// https://stackoverflow.com/questions/18930910/access-struct-property-by-name
func getField(v *Device, field string) reflect.Value {
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
	device := new(Device)
	device.extract = *extractPtr

	// If we ask for auto-discovery, try it and exit.
	if *discoverPtr {
		for _, s := range config.GetSearches() {
			device.search = s
			findMatchingDevice(device)
			if device.exists {
				break
			}
		}
	} else {
		s := &config.Search{Field: *fieldPtr, Query: *searchPtr}
		device.search = *s
		findMatchingDevice(device)
	}

	if *existsPtr {
		// If we're explicitly asking to see if it exists, say no, and
		// return a zero error code.
		if device.exists {
			fmt.Println("true")
		} else {
			fmt.Println("false")
		}
		os.Exit(0)
	} else if device.exists {
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
