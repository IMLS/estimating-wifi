package tlp

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"strconv"
	"strings"

	"github.com/jadudm/imls-wifi-adapter-search/models"
	"github.com/jadudm/imls-wifi-adapter-search/search"
	"gsa.gov/18f/session-counter/config"
	"gsa.gov/18f/session-counter/constants"
)

func tshark(cfg *config.Config) []string {

	tsharkCmd := exec.Command(cfg.Wireshark.Path,
		"-a", fmt.Sprintf("duration:%d", cfg.Wireshark.Duration),
		"-I", "-i", cfg.Wireshark.Adapter,
		"-Tfields", "-e", "wlan.sa")

	tsharkOut, _ := tsharkCmd.StdoutPipe()
	tsharkCmd.Start()
	tsharkBytes, _ := ioutil.ReadAll(tsharkOut)
	tsharkCmd.Wait()
	macs := strings.Split(string(tsharkBytes), "\n")

	// Let's not worry about how many packets.
	// Just track the MAC addresses.
	// pkts := make(map[string]int)
	// for _, a_mac := range macs {
	// 	v, ok := pkts[a_mac]
	// 	if ok {
	// 		pkts[a_mac] = v + 1
	// 	} else {
	// 		pkts[a_mac] = 1
	// 	}
	// }

	return macs
}

/* PROCESS runWireshark
 * Runs a subprocess for a duration of OBSERVE_SECONDS.
 * Therefore, this process effectively blocks for that time.
 * Gathers a hashmap of [MAC -> count] values. This hashmap
 * is then communicated out.
 * Empty MAC addresses are filtered out.
 */
func RunWireshark(ka *Keepalive, cfg *config.Config, in <-chan bool, out chan []string) {
	log.Println("Starting runWireshark")
	// If we have to wait twice the monitor duration, something broke.
	ping, pong := ka.Subscribe("runWireshark", cfg.Wireshark.Duration*2)

	// Adapter count... every "ac" ticks, we look up the adapter.
	// (ac % 0) guarantees that we look it up the first time.
	ticker := 0
	minutes_interval, _ := strconv.Atoi(cfg.Wireshark.CheckWlan)

	for {
		select {

		case <-ping:
			// We ping faster than this process can reply. However, we have a long
			// enough timeout that we will *eventually* catch up with all of the pings.
			pong <- "wireshark"

		case <-in:
			// Look up the adapter. Use the find-ralink library.
			// The % will trigger first time through, which we want.
			var dev *models.Device = nil
			// If the config doesn't have this in it, we get a divide by zero.
			if (ticker % minutes_interval) == 0 {
				dev = search.SearchForMatchingDevice()
				log.Println(dev)
			}

			// Only do a reading and continue the pipeline
			// if we find an adapter.
			if dev != nil && dev.Exists {
				// Load the config for use.
				cfg.Wireshark.Adapter = dev.Logicalname
				// This will block for [cfg.Wireshark.Duration] seconds.
				macmap := tshark(cfg)
				// Mark and remove too-short MAC addresses
				// for removal from the tshark findings.
				var keepers []string
				// for `k, _ :=` is the same as `for k :=`
				for _, k := range macmap {
					if len(k) >= constants.MACLENGTH {
						keepers = append(keepers, k)
					}
				}
				// Report out the cleaned MACmap.
				out <- keepers
			} else {
				log.Println("No wifi device found. No scanning carried out.")
				// Report an empty array of keepers
				out <- make([]string, 0)
			}

			// Bump our ticker
			ticker += 1
		}
	}
}
