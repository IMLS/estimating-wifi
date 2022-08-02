package state

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"os"
	"regexp"
	"runtime"
	"strings"

	"github.com/rs/zerolog/log"
	"gsa.gov/18f/internal/wifi-hardware-search/netadapter"
)

const FakeSerial = "CESTNEPASUNESERIE"
const FakeSerialCheck = "PAS"

var (
	serialWarnGiven    = false
	GetSerialPSCommand = "(Get-CimInstance -Class Win32_ComputerSystemProduct).UUID"
)

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
		} else if runtime.GOOS == "windows" {
			ps := netadapter.New()
			lines := ps.Execute(GetSerialPSCommand)
			serial := strings.TrimSpace(string(lines)) // remove \r\n
			hash := sha256.Sum256([]byte(serial))
			cache["serial"] = fmt.Sprintf("%x", hash)
		} else {
			if !serialWarnGiven {
				log.Warn().Msg("Cannot grab serial number.")
				serialWarnGiven = true
			}

		}
	}
	return serial
}
