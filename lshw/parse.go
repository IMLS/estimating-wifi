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

	arr := make([]string, 0)
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		arr = append(arr, line)
	}

	return ParseLSHW(arr)
}

func deepCopy(h map[string]string) map[string]string {
	nh := make(map[string]string)
	for k, v := range h {
		nh[k] = v
	}

	return nh
}

func ParseLSHW(string_array []string) []map[string]string {
	hash := make(map[string]string)
	sectionHeading := regexp.MustCompile(`^\s*\*-(usb|network)((?:\:\d))?\s*`)
	entryPattern := regexp.MustCompile(`\s*([a-z ]+):\s+(.*)\s*`)

	// Start looking for a section heading.
	state := constants.LOOKING_FOR_SECTION_HEADING

	// Build up an array of hashes. Instead of looking for the device here,
	// we'll instead collect all the devices into hashes, and hold them for a moment.
	devices := make([]map[string]string, 0)

	for _, line := range string_array {
		switch state {
		case constants.LOOKING_FOR_SECTION_HEADING:
			// See if we can find a section heading.
			match := sectionHeading.MatchString(line)
			//fmt.Printf("LFSH line [%v] match [%v]\n", line, match)

			// If we do, change state.
			if match {
				// Create a new hash.
				hash = make(map[string]string)
				state = constants.READING_ENTRY
			}
		// Now we're in an lshw entry.
		case constants.READING_ENTRY:
			newSecMatch := sectionHeading.MatchString(line)
			hashMatch := entryPattern.MatchString(line)
			hashPieces := entryPattern.FindStringSubmatch(line)
			//fmt.Printf("RE   line [%v] nsm [%v] hm [%v]\n", line, newSecMatch, hashMatch)

			if newSecMatch {
				state = constants.READING_ENTRY
				devices = append(devices, deepCopy(hash))
				// Create a new hash to continue reading into
				hash = make(map[string]string)
			} else if hashMatch {
				// 0 is the full string, 1 the first group, 2 the second.
				key := hashPieces[1]
				value := hashPieces[2]
				hash[key] = value
			}
		}
	}

	// Don't lose the last hash!
	devices = append(devices, hash)
	if *config.Verbose {
		fmt.Println("found", len(devices), "devices")
		fmt.Println("devices\n", devices)
	}

	return devices
}
