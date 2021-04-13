package lshw

import (
	"bufio"
	"fmt"
	"log"
	"os/exec"
	"regexp"

	"gsa.gov/18f/find-ralink/config"
	"gsa.gov/18f/find-ralink/constants"
	"gsa.gov/18f/find-ralink/models"
)

func GetDeviceHash(wlan *models.Device) []map[string]string {
	wlan.Exists = false

	cmd := exec.Command(constants.LSHW_EXE, "-class", "network")
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
	state := constants.LOOKING_FOR_USB

	// Build up an array of hashes. Instead of looking for the device here,
	// we'll instead collect all the devices into hashes, and hold them for a moment.
	devices := make([]map[string]string, 0)

	for scanner.Scan() {
		line := scanner.Text()
		switch state {
		case constants.LOOKING_FOR_USB:
			match := usbSecRe.MatchString(line)
			if match {
				if config.Verbose {
					fmt.Println("-> READING_HASH")
				}
				// Create a new hash.
				hash = make(map[string]string)
				state = constants.READING_HASH
			}
		case constants.READING_HASH:
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
				state = constants.DONE_READING
			} else if hashMatch {
				// fmt.Printf("%v <- %v\n", hashPieces[1], hashPieces[2])
				// 0 is the full string, 1 the first group, 2 the second.
				hash[hashPieces[1]] = hashPieces[2]
			}
		case constants.DONE_READING:
			state = constants.LOOKING_FOR_USB
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
