package api

import (
	"encoding/json"
	"log"
	"time"

	"gsa.gov/18f/session-counter/config"
)

// PROC We need some state.
type EventLogger struct {
	Cfg    *config.Config
	Server *config.Server
}

func (el *EventLogger) init(cfg *config.Config, svr *config.Server) {
	el.Server = svr
	el.Cfg = cfg
}

func NewEventLogger(cfg *config.Config, svr *config.Server) *EventLogger {
	el := new(EventLogger)
	el.init(cfg, svr)
	return el
}

// "event_id":    strconv.Itoa(session_id),
// 		"device_uuid": config.GetSerial(),
// 		"lib_user":    tok.User,
// 		"localtime":   time.Now().Format(time.RFC3339),
// 		// FIXME: The server needs to auto-set this
// 		"servertime": time.Now().Format(time.RFC3339),
// 		"session_id": cfg.SessionId,
// 		"device_id":  uid,
// 		"last_seen":  strconv.Itoa(count),

func (el *EventLogger) Log(tag string, info map[string]string) int {
	var uri string = (el.Server.Host + el.Server.Eventpath)
	log.Println("event log uri:", uri)
	tok, _ := GetToken(el.Server)

	// ALWAYS log a hash table. Makes processing easier later.
	var asJson []byte
	if info == nil {
		asJson, _ = json.Marshal(map[string]string{})
	} else {
		asJson, _ = json.Marshal(info)
	}
	//asB64 := b64.URLEncoding.EncodeToString([]byte(asJson))

	data := map[string]string{
		"device_uuid": config.GetSerial(),
		"lib_user":    tok.User,
		"session_id":  el.Cfg.SessionId,
		"localtime":   time.Now().Format(time.RFC3339),
		"servertime":  time.Now().Format(time.RFC3339),
		"tag":         tag,
		"info":        string(asJson),
	}
	ndx, _ := postJSON(el.Server, tok, uri, data)
	return ndx

}
