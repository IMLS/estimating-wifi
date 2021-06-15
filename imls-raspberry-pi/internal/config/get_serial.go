package config

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

	if err != nil {
		log.Println("error opening /proc/cpuinfo")
		log.Println(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Println("error reading /proc/cpuino")
		log.Println(err)
	}

	return lines
}

func GetSerial() string {
	serial := FakeSerial
	// Try and pull from the cache, so we don't keep opening up a /proc filesystem...
	if val, ok := cache["serial"]; ok {
		serial = val
	} else {
		if runtime.GOOS == "linux" && runtime.GOARCH == "arm" {
			lines := cpuinfoLines()
			re := regexp.MustCompile(`Serial\s+:\s+([a-f0-9]+)`)
			for _, line := range lines {
				// log.Println("line", line)
				matched := re.FindStringSubmatch(line)
				if len(matched) > 0 {
					// log.Println("matched", matched)
					serial = string(matched[1])
					cache["serial"] = serial
				}
			}
		} else {
			if !serialWarnGiven {
				log.Println("Not running on an RPi. Cannot grab serial number. Exiting.")
				serialWarnGiven = true
			}

		}
	}

	return serial
}
