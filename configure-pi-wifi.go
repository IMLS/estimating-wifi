package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os/exec"
	"reflect"
	"regexp"
	"strconv"
)

const (
	LOOKING_FOR_USB = iota
	READING_HASH    = iota
	DONE_READING    = iota
)

type RAlink struct {
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

	// output, err := ioutil.ReadAll(stdout)
	// if err != nil {
	// 	log.Println("cpw: unable to read all of stdout from `lshw`")
	// 	log.Fatal(err)
	// }

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
				//fmt.Println("-> READING_HASH")
				state = READING_HASH
			}
		case READING_HASH:
			//fmt.Printf("checking: [ %v ]", line)
			newSecMatch := newSecRe.MatchString(line)
			hashMatch := hashRe.MatchString(line)
			hashPieces := hashRe.FindStringSubmatch(line)

			if newSecMatch {
				//fmt.Println("-> DONE_READING")
				state = DONE_READING
			} else if hashMatch {
				// fmt.Printf("%v <- %v\n", hashPieces[1], hashPieces[2])
				hash[hashPieces[1]] = hashPieces[2]
			}
		case DONE_READING:
			// SKIP
		}
	}

	wlan.physicalId, _ = strconv.Atoi(hash["physical id"])
	wlan.description = hash["description"]
	wlan.busInfo = hash["bus info"]
	wlan.logicalName = hash["logical name"]
	wlan.serial = hash["serial"]
	if len(hash["serial"]) >= 17 {
		wlan.mac = hash["serial"][0:17]
	} else {
		wlan.mac = hash["serial"]
	}

	wlan.configuration = hash["configuration"]
	return wlan
}

// https://stackoverflow.com/questions/18930910/access-struct-property-by-name
func getField(v *RAlink, field string) string {
	r := reflect.ValueOf(v)
	f := reflect.Indirect(r).FieldByName(field)
	return f.String()
}

func main() {
	fieldPtr := flag.String("descriptor", "logicalName", "Descriptor to extract from device.")
	flag.Parse()

	device := getRAlinkDevice()
	res := getField(&device, *fieldPtr)
	fmt.Println(res)
}
