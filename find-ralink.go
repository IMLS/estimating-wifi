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

	"github.com/jedib0t/go-pretty/v6/text"
	"gsa.gov/18f/find-ralink/constants"
)

const (
	LOOKING_FOR_USB = iota
	READING_HASH    = iota
	DONE_READING    = iota
)

type RAlink struct {
	exists        bool
	physicalId    int
	description   string
	busInfo       string
	logicalName   string
	serial        string
	mac           string
	configuration string
}

func getRAlinkDevice() RAlink {
	wlan := RAlink{}
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
	hash := make(map[string]string)
	usbSecRe := regexp.MustCompile(`^\s+\*-usb`)
	newSecRe := regexp.MustCompile(`^\s+\*-.*`)
	hashRe := regexp.MustCompile(`^\s+(.*?): (.*)`)
	state := LOOKING_FOR_USB

	for scanner.Scan() {
		line := scanner.Text()
		switch state {
		case LOOKING_FOR_USB:
			match := usbSecRe.MatchString(line)
			if match {
				// fmt.Println("-> READING_HASH")
				state = READING_HASH
			}
		case READING_HASH:
			// fmt.Printf("checking: [ %v ]\n", line)
			newSecMatch := newSecRe.MatchString(line)
			hashMatch := hashRe.MatchString(line)
			hashPieces := hashRe.FindStringSubmatch(line)

			if newSecMatch {
				// fmt.Println("-> DONE_READING")
				state = DONE_READING
			} else if hashMatch {
				// fmt.Printf("%v <- %v\n", hashPieces[1], hashPieces[2])
				// 0 is the full string, 1 the first group, 2 the second.
				hash[hashPieces[1]] = hashPieces[2]
			}
		case DONE_READING:
			// SKIP
		}
	}

	v, _ := regexp.MatchString("Ralink", hash["vendor"])
	if v {
		wlan.exists = true
	}

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
	return wlan
}

// https://stackoverflow.com/questions/18930910/access-struct-property-by-name
func getField(v *RAlink, field string) reflect.Value {
	r := reflect.ValueOf(v)
	f := reflect.Indirect(r).FieldByName(field)
	return f
}

func main() {
	fieldPtr := flag.String("descriptor", "logicalName", "Descriptor to extract from device.")
	existsPtr := flag.Bool("exists", false, "Ask if a device exists. Returns `true` or `false`.")
	flag.Parse()

	if os.Getenv("USER") != "root" {
		fmt.Println(text.FgRed.Sprint("find-ralink *really* needs to be run as root."))
	}

	// Essentially a shortcut...
	// Overrides the --descriptor field.
	if *existsPtr {
		*fieldPtr = "exists"
	}

	device := getRAlinkDevice()

	if device.exists {
		res := getField(&device, *fieldPtr)
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
