package lshw

import (
	"bufio"
	"fmt"
	"log"
	"os/exec"
	"regexp"

	"gsa.gov/18f/find-ralink/config"
	"gsa.gov/18f/find-ralink/models"
)

// PURPOSE
// This function calls out to `lshw` and
// then passes it off for parsing into a hashmap.
func GetDeviceHash(wlan *models.Device) []map[string]string {
	wlan.Exists = false

	cmd := exec.Command(config.LSHW_EXE, "-class", "network")
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

// PURPOSE
// We believe our slice construction had a pass-by-reference issue.
// This makes a fresh copy of a hashmap so that we can insert the
// devices that we find into a slice, and pass that slice-of-maps
// back for use in the main().
func deepCopy(h map[string]string) map[string]string {
	nh := make(map[string]string)
	for k, v := range h {
		nh[k] = v
	}

	return nh
}

// PURPOSE
// This function takes an array of strings (representing the output of `lshw`)
// and parses them into a list of hashes. Each map represents a piece of hardware
// attached to the machine. The keys are the descriptors provided by
// `lshw`, and the values are... the values reported by `lshw`.
func ParseLSHW(string_array []string) []map[string]string {
	sectionHeading := regexp.MustCompile(`^\s*\*-(usb|network)((?:\:\d))?\s*`)
	entryPattern := regexp.MustCompile(`\s*([a-z ]+):\s+(.*)\s*`)

	// Build up an array of hashes. Instead of looking for the device here,
	// we'll instead collect all the devices into hashes, and hold them for a moment.
	devices := make([]map[string]string, 0)
	// Make sure the hash is in scope for the loop
	hash := make(map[string]string)
	// Start looking for a section heading.
	// state := constants.LOOKING_FOR_SECTION_HEADING

	for _, line := range string_array {
		newSecMatch := sectionHeading.MatchString(line)
		hashMatch := entryPattern.MatchString(line)
		if *config.Verbose {
			fmt.Printf("RE   line [%v] nsm [%v] hm [%v]\n", line, newSecMatch, hashMatch)
		}

		if newSecMatch {
			// If we find a new section, and we have something in the hash,
			// copy it and store it to be passed back.
			if len(hash) > 0 {
				devices = append(devices, deepCopy(hash))
				hash = make(map[string]string)
			}
		} else if hashMatch {
			// If we are in the middle of a hash, then see if we can
			// pull it apart and keep the pieces.
			hashPieces := entryPattern.FindStringSubmatch(line)
			// 0 is the full string, 1 the first group, 2 the second.
			key := hashPieces[1]
			value := hashPieces[2]
			hash[key] = value
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
