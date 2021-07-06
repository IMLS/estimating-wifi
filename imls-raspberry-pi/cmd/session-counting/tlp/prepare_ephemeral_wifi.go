package tlp

import (
	"strconv"
	"strings"
	"time"

	"gsa.gov/18f/config"
	"gsa.gov/18f/logwrapper"
)

// Converts raw data to a map[string]string
// This makexs it ready for storage locally (SQLite) or
// via an API (where everything becomes text anyway).
func PrepEphemeralWifi(ka *Keepalive, cfg *config.Config, kb *Broker,
	in_hash <-chan map[string]int, out_arr chan<- []map[string]interface{}) {
	lw := logwrapper.NewLogger(nil)
	lw.Debug("Starting PrepEphemeralWifi")
	var ping, pong chan interface{} = nil, nil
	var ch_kill chan interface{} = nil
	if kb != nil {
		ch_kill = kb.Subscribe()
	} else {
		ping, pong = ka.Subscribe("PrepEphemeralWifi", 30)
	}

	event_ndx := 0

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
			lw.Debug("event ndx: %v", event_ndx)
			lw.Length("macs-to-store", remove)
			// Now, bundle that as an array of hashmaps.
			reportArr := make([]map[string]interface{}, 0)
			for anondevice := range h {
				data := map[string]interface{}{
					// event_id
					// An "event" was registered before the data is inserted. This is
					// essentially a FK into the events table.
					"event_id": strconv.Itoa(event_ndx),
					// The session id is a unique ID that is generated at powerup.
					"session_id": cfg.SessionId,
					// The time on the device.
					"localtime": time.Now().Format(time.RFC3339),
					// The serial number of the Pi.
					"pi_serial": config.GetSerial(),
					// The FCFS Seq Id entered at setup time.
					"fcfs_seq_id": cfg.Auth.FCFSId,
					// The tag entered at setup time.
					"device_tag": cfg.Auth.DeviceTag,
					// The "anondevice" is now something like "0:32" or "26:384"
					// We split that into a manufacturer ID and a device ID.
					// The manufacturer Ids are consistent for a session (a powerup cycle)
					// The patron id is tracked for 2 hours (or whatever the config says)
					"manufacturer_index": strings.Split(anondevice, ":")[0],
					"patron_index":       strings.Split(anondevice, ":")[1],
				}

				reportArr = append(reportArr, data)
			}

			event_ndx += 1
			out_arr <- reportArr

		}
	}

}
