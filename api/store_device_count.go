package api

import (
	"strconv"
	"time"

	"gsa.gov/18f/session-counter/config"
)

// FUNC StoreDeviceCount
// Stores the device count JSON via Umbrella
func StoreDeviceCount(cfg *config.Config, tok *config.AuthConfig, session_id int, uid string, count int) error {
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
		"last_seen":  strconv.Itoa(count),
	}

	postJSON(cfg, tok, uri, data)
	return nil
}
