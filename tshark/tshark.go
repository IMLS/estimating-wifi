package tshark

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"

	"gsa.gov/18f/session-counter/model"
)

func Tshark(cfg model.Config) map[string]int {

	tsharkCmd := exec.Command(cfg.Wireshark.Path,
		"-a", fmt.Sprintf("duration:%d", cfg.Wireshark.Duration),
		"-I", "-i", cfg.Wireshark.Adapter,
		"-Tfields", "-e", "wlan.sa")

	tsharkOut, _ := tsharkCmd.StdoutPipe()
	tsharkCmd.Start()
	tsharkBytes, _ := ioutil.ReadAll(tsharkOut)
	tsharkCmd.Wait()
	macs := strings.Split(string(tsharkBytes), "\n")

	pkts := make(map[string]int)

	for _, a_mac := range macs {
		v, ok := pkts[a_mac]
		if ok {
			pkts[a_mac] = v + 1
		} else {
			pkts[a_mac] = 1
		}
	}

	return pkts
}
