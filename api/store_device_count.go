package api

import (
	"strconv"
	"strings"
	"time"

	"gsa.gov/18f/session-counter/config"
)

// This posts an array of data to ReVal.
// Filter out things that were last seen more than zero minutes ago.
func StoreDevicesCount(cfg *config.Config, tok *config.AuthConfig, session_id int, h map[string]int) error {
	uri := FormatUri(cfg.Umbrella.Scheme, cfg.Umbrella.Host, cfg.Umbrella.Data)

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

	// Now, bundle that as an array of hashmaps.
	reportArr := make([]map[string]string, 0)
	for anondevice := range h {

		data := map[string]string{
			// event_id
			// An "event" was registered before the data is inserted. This is
			// essentially a FK into the events table.
			"event_id": strconv.Itoa(session_id),
			// The session id is a unique ID that is generated at powerup.
			"session_id": cfg.SessionId,
			// The time on the device.
			"localtime": time.Now().Format(time.RFC3339),
			// The serial number of the Pi.
			"pi_serial": config.GetSerial(),
			// The FCFS Seq Id entered at setup time.
			"fcfs_seq_id": tok.FCFSId,
			// The tag entered at setup time.
			"device_tag": tok.DeviceTag,
			// The "anondevice" is now something like "0:32" or "26:384"
			// We split that into a manufacturer ID and a device ID.
			// The manufactuerer Ids are consistent for a session (a powerup cycle)
			// The patron id is tracked for 2 hours (or whatever the config says)
			"manufacturer_index": strings.Split(anondevice, ":")[0],
			"patron_index":       strings.Split(anondevice, ":")[1],
		}

		reportArr = append(reportArr, data)
	}

	postJSON(cfg, tok, uri, reportArr)
	return nil

}
