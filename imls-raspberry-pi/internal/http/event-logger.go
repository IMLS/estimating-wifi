package http

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"gsa.gov/18f/config"
)

func NewEventLogger(cfg *config.Config) *EventLogger {
	el := new(EventLogger)
	el.Cfg = cfg
	return el
}

func (el *EventLogger) Log(tag string, info map[string]string) (int, error) {
	uri := FormatUri(el.Cfg.Umbrella.Scheme, el.Cfg.Umbrella.Host, el.Cfg.Umbrella.Data)

	log.Println("event log uri:", uri)
	tok, _ := config.ReadAuth()

	// ALWAYS log a hash table. Makes processing easier later.
	var asJson []byte
	var err error

	if info == nil {
		asJson, _ = json.Marshal(map[string]string{})
	} else {
		asJson, err = json.Marshal(info)
		if err != nil {
			asJson, _ = json.Marshal(map[string]string{
				"msg":   "could not marshal info for tag",
				"error": fmt.Sprint(err),
			})
		}
	}
	//asB64 := b64.URLEncoding.EncodeToString([]byte(asJson))

	data := map[string]string{
		"pi_serial":   config.GetSerial(),
		"fcfs_seq_id": tok.FCFSId,
		"device_tag":  tok.DeviceTag,
		"session_id":  el.Cfg.SessionId,
		"localtime":   time.Now().Format(time.RFC3339),
		"tag":         tag,
		"info":        string(asJson),
	}
	ndx, err := PostJSON(uri, []map[string]string{data})
	return ndx, err

}
