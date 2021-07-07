package tlp

import (
	"strconv"
	"strings"
	"time"

	"gsa.gov/18f/analysis"
	"gsa.gov/18f/config"
	"gsa.gov/18f/logwrapper"
)

// Converts raw data to a map[string]string
// This makexs it ready for storage locally (SQLite) or
// via an API (where everything becomes text anyway).
func PrepEphemeralWifi(ka *Keepalive, cfg *config.Config, kb *Broker,
	in_hash <-chan map[string]int, out_arr chan<- []analysis.WifiEvent) {
	lw := logwrapper.NewLogger(nil)
	lw.Debug("Starting PrepEphemeralWifi")
	var ping, pong chan interface{} = nil, nil
	var ch_kill chan interface{} = nil
	if kb != nil {
		ch_kill = kb.Subscribe()
	} else {
		ping, pong = ka.Subscribe("PrepEphemeralWifi", 30)
	}

	for {
		select {
		case <-ping:
			pong <- "PrepEphemeralWifi"
		case <-ch_kill:
			lw.Debug("exiting PrepEphemeralWifi")
			return

		// Block waiting to read the incoming hash
		case h := <-in_hash:
			// Remove all the UIDs that we saw more than 0 minutes ago
			var remove []string
			for k, v := range h {
				if v > 0 {
					remove = append(remove, k)
				}
			}
			for _, r := range remove {
				delete(h, r)
			}
			// lw.Length("macs-to-store", remove)
			// Now, bundle that as an array of hashmaps.
			reportArr := make([]analysis.WifiEvent, 0)

			for anondevice := range h {
				mfg, _ := strconv.Atoi(strings.Split(anondevice, ":")[0])
				pid, _ := strconv.Atoi(strings.Split(anondevice, ":")[1])

				data := analysis.WifiEvent{
					SessionId:         cfg.SessionId,
					Localtime:         time.Now().Format(time.RFC3339),
					FCFSSeqId:         cfg.Auth.FCFSId,
					DeviceTag:         cfg.Auth.DeviceTag,
					ManufacturerIndex: mfg,
					PatronIndex:       pid,
				}

				reportArr = append(reportArr, data)
			}

			out_arr <- reportArr

		}
	}

}
