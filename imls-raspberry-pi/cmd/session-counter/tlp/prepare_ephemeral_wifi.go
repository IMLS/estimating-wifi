package tlp

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"gsa.gov/18f/internal/config"
	"gsa.gov/18f/internal/logwrapper"
	"gsa.gov/18f/internal/state"
	"gsa.gov/18f/internal/structs"
)

// PrepEphemeralWifi converts raw data to a map[string]string. This makes it
// ready for storage locally (SQLite) or via an API (where everything becomes
// text anyway).
func PrepEphemeralWifi(ka *Keepalive, cfg *config.Config, kb *KillBroker,
	inHash <-chan map[string]int, outArr chan<- []structs.WifiEvent) {
	lw := logwrapper.NewLogger(nil)
	lw.Debug("Starting PrepEphemeralWifi")
	var ping, pong chan interface{} = nil, nil
	var chKill chan interface{} = nil
	if kb != nil {
		chKill = kb.Subscribe()
	} else {
		ping, pong = ka.Subscribe("PrepEphemeralWifi", 30)
	}

	for {
		select {
		case <-ping:
			pong <- "PrepEphemeralWifi"
		case <-chKill:
			lw.Debug("exiting PrepEphemeralWifi")
			return

		// Block waiting to read the incoming hash
		case h := <-inHash:
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
			reportArr := make([]structs.WifiEvent, 0)

			for anondevice := range h {
				mfg, _ := strconv.Atoi(strings.Split(anondevice, ":")[0])
				pid, _ := strconv.Atoi(strings.Split(anondevice, ":")[1])

				data := structs.WifiEvent{
					SessionID:         fmt.Sprint(state.GetCurrentSessionID(cfg)),
					Localtime:         cfg.Clock.Now().Format(time.RFC3339),
					FCFSSeqID:         cfg.Auth.FCFSId,
					DeviceTag:         cfg.Auth.DeviceTag,
					ManufacturerIndex: mfg,
					PatronIndex:       pid,
				}

				reportArr = append(reportArr, data)
			}

			outArr <- reportArr

		}
	}

}
