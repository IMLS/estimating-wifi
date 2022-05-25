package state

import (
	"bufio"
	"os"
	"regexp"
	"runtime"

	"github.com/rs/zerolog/log"
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

func getCachedSerial() string {
	serial := FakeSerial
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
				log.Warn().Msg("Not running on a Raspberry Pi. Cannot grab serial number.")
				serialWarnGiven = true
			}

		}
	}
	return serial
}
