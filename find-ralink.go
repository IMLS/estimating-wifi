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
	"gsa.gov/18f/find-ralink/constants"
)

const (
	LOOKING_FOR_USB = iota
	READING_HASH    = iota
	DONE_READING    = iota
)

type Device struct {
	exists        bool
	searchString  string
	physicalId    int
	description   string
	busInfo       string
	logicalName   string
	serial        string
	mac           string
	configuration string
}

func getRAlinkDevice(wlan *Device, verbose bool) {
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
				if verbose {
					fmt.Println("-> READING_HASH")
				}
				// Create a new hash.
				hash = make(map[string]string)
				state = READING_HASH
			}
		case READING_HASH:
			if verbose {
				fmt.Printf("checking: [ %v ]\n", line)
			}
			newSecMatch := newSecRe.MatchString(line)
			hashMatch := hashRe.MatchString(line)
			hashPieces := hashRe.FindStringSubmatch(line)

			if newSecMatch {
				if verbose {
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
			if verbose {
				fmt.Println("devices len", len(devices))
			}
		}
	}

	// Don't lose the last hash!
	devices = append(devices, hash)

	// Now, go through the devices and find the one that matches our criteria.
	// Either we're looking for a vendor or for something in the config string.
	// A more general search could be implemented, but I'll keep it limited to prevent
	// spurious matches that were unexpected/surprising in practice. (I hope.)
	for _, hash := range devices {

		if verbose {
			fmt.Println("---------")
			for k, v := range hash {
				fmt.Println(k, "<-", v)
			}
		}

		// NOTE: Do the searches case insensitive.
		v, _ := regexp.MatchString(strings.ToLower(wlan.searchString), strings.ToLower(hash["vendor"]))
		if v {
			wlan.exists = true
		}

		v, _ = regexp.MatchString(strings.ToLower(wlan.searchString), strings.ToLower(hash["configuration"]))
		if v {
			if verbose {
				fmt.Println("Found config matching", wlan.searchString)
			}
			wlan.exists = true
		}

		// If we find something, proceed.
		// Back in 'main', we'll handle the case where wlan.exists is false.
		if wlan.exists {
			wlan.physicalId, _ = strconv.Atoi(hash["physical id"])
			wlan.description = hash["description"]
			wlan.busInfo = hash["bus info"]
			wlan.logicalName = hash["logical name"]
			wlan.serial = hash["serial"]

			if len(hash["serial"]) >= constants.MACLENGTH {
				wlan.mac = hash["serial"][0:constants.MACLENGTH]
			} else {
				wlan.mac = hash["serial"]
			}
			wlan.configuration = hash["configuration"]

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
	verbosePtr := flag.Bool("verbose", false, "Verbose output.")
	mfgPtr := flag.String("search", "Ralink", "Search string to use in hardware listing.")
	fieldPtr := flag.String("descriptor", "logicalName", "Descriptor to extract from device.")
	existsPtr := flag.Bool("exists", false, "Ask if a device exists. Returns `true` or `false`.")
	versionPtr := flag.Bool("version", false, "Get the software version and exit.")

	flag.Parse()

	// If they just want the version, print and exit.
	if *versionPtr {
		fmt.Println("Version ", constants.VERSION)
		os.Exit(0)
	}

	if os.Getenv("USER") != "root" {
		fmt.Println(text.FgRed.Sprint("find-ralink *really* needs to be run as root."))
	}

	// Essentially a shortcut...
	// Overrides the --descriptor field.
	if *existsPtr {
		*fieldPtr = "exists"
	}

	device := new(Device)
	device.searchString = *mfgPtr
	getRAlinkDevice(device, *verbosePtr)

	if device.exists {
		res := getField(device, *fieldPtr)
		if reflect.TypeOf(res).Kind() == reflect.Bool {
			fmt.Println(res.Interface())
		} else {
			fmt.Println(res)
		}
		os.Exit(0)
	} else {
		// If we're explicitly asking to see if it exists, say no, and
		// return a zero error code.
		if *fieldPtr == "exists" {
			fmt.Println("false")
			os.Exit(0)
		} else {
			// Otherwise... this is bad. Exit with error.
			fmt.Println("Device not found")
			os.Exit(-1)
		}
	}
}
