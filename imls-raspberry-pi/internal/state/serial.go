package state

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"runtime"
)

const FakeSerial = "CESTNEPASUNESERIE"
const FakeSerialCheck = "PAS"

var serialWarnGiven = false

// Create a cache, so repeated calls to get the serial don't
// open up endless file sockets...
var cache map[string]string = make(map[string]string)

func cpuinfoLines() (lines []string) {
	file, err := os.Open("/proc/cpuinfo")
	// If we can't find/read cpuinfo, return an empty list.
	if err != nil {
		return make([]string, 0)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return make([]string, 0)
	}
	return lines
}

func decodeSerial() {
	serial := FakeSerial
	if the_config.Serial != "" {
		// optionally override with a pre-defined serial. note
		// that this is only for non-arm usage.
		serial = the_config.Serial
	}
	// Try and pull from the cache, so we don't keep opening up a /proc filesystem...
	if val, ok := cache["serial"]; ok {
		serial = val
	} else {
		if runtime.GOOS == "linux" && runtime.GOARCH == "arm" {
			lines := cpuinfoLines()
			re := regexp.MustCompile(`Serial\s+:\s+([a-f0-9]+)`)
			for _, line := range lines {
				matched := re.FindStringSubmatch(line)
				if len(matched) > 0 {
					serial = string(matched[1])
					cache["serial"] = serial
				}
			}
		} else {
			if !serialWarnGiven {
				log.Println("Not running on an RPi. Cannot grab serial number.")
				serialWarnGiven = true
			}

		}
	}
	the_config.Serial = serial
}

func GetSerial() string {
	return the_config.Serial
}
