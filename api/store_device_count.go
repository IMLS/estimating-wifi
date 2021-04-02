package api

import (
	"strconv"
	"time"

	"gsa.gov/18f/session-counter/config"
)

// FUNC StoreDeviceCount
// Stores the device count JSON via Umbrella
func StoreDeviceCount(cfg *config.Config, tok *config.AuthConfig, session_id int, uid string) error {
	uri := FormatUri(cfg.Umbrella.Scheme, cfg.Umbrella.Host, cfg.Umbrella.Data)

	data := map[string]string{
		// FIXME: This needs to be captured first and passed in.
		"event_id":    strconv.Itoa(session_id),
		"device_uuid": config.GetSerial(),
		"lib_user":    tok.Email,
		"localtime":   time.Now().Format(time.RFC3339),
		// FIXME: The server needs to auto-set this
		"servertime": time.Now().Format(time.RFC3339),
		"session_id": cfg.SessionId,
		"device_id":  uid,
	}

	postJSON(cfg, tok, uri, []map[string]string{data})
	return nil
}

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
	for k := range h {
		data := map[string]string{
			// FIXME: This needs to be captured first and passed in.
			"event_id":    strconv.Itoa(session_id),
			"device_uuid": config.GetSerial(),
			"lib_user":    tok.Email,
			"localtime":   time.Now().Format(time.RFC3339),
			// FIXME: The server needs to auto-set this
			"servertime": time.Now().Format(time.RFC3339),
			"session_id": cfg.SessionId,
			"device_id":  k,
		}
		reportArr = append(reportArr, data)
	}

	postJSON(cfg, tok, uri, reportArr)
	return nil

}
