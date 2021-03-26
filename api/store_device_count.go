package api

import (
	"strconv"
	"time"

	"gsa.gov/18f/session-counter/config"
	"gsa.gov/18f/session-counter/model"
)

// FUNC StoreDeviceCount
// Stores the device count JSON to directus and reval.
func StoreDeviceCount(cfg *config.Config, svr *config.Server, tok *model.Auth, session_id int, uid string, count int) error {
	var uri string = (svr.Host + svr.Datapath)

	data := map[string]string{
		// FIXME: This needs to be captured first and passed in.
		"event_id":    strconv.Itoa(session_id),
		"device_uuid": config.GetSerial(),
		"lib_user":    tok.User,
		"localtime":   time.Now().Format(time.RFC3339),
		// FIXME: The server needs to auto-set this
		"servertime": time.Now().Format(time.RFC3339),
		"session_id": cfg.SessionId,
		"device_id":  uid,
		"last_seen":  strconv.Itoa(count),
	}

	err := postJSON(svr, tok, uri, data)
	return err
}
