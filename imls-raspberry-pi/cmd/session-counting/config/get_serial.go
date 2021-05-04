package config

import (
	"bufio"
	"log"
	"os"
	"regexp"
)

const FakeSerial = "CESTNEPASUNESERIE"
const FakeSerialCheck = "PAS"

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
	lines := cpuinfoLines()
	serial := FakeSerial
	re := regexp.MustCompile(`Serial\s+:\s+([a-f0-9]+)`)
	for _, line := range lines {
		// log.Println("line", line)
		matched := re.FindStringSubmatch(line)
		if len(matched) > 0 {
			// log.Println("matched", matched)
			serial = string(matched[1])
		}
	}

	return serial
}
