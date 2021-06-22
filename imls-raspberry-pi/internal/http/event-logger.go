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

func (el *EventLogger) LogJSON(tag string, json string) (int, error) {
	uri := FormatUri(el.Cfg.Umbrella.Scheme, el.Cfg.Umbrella.Host, el.Cfg.Umbrella.Logging)

	data := map[string]string{
		"pi_serial":   config.GetSerial(),
		"fcfs_seq_id": el.Cfg.Auth.FCFSId,
		"device_tag":  el.Cfg.Auth.DeviceTag,
		"session_id":  el.Cfg.SessionId,
		"localtime":   time.Now().Format(time.RFC3339),
		"tag":         tag,
		"info":        json,
	}

	ndx, err := PostJSON(el.Cfg, uri, []map[string]string{data})
	if el.Cfg.StorageMode == "api" {
		err = nil
	}
	return ndx, err

}

func (el *EventLogger) Log(tag string, info map[string]string) (int, error) {
	// ALWAYS log a hash table. Makes processing easier later.
	var asJson []byte
	var err error

	if info == nil {
		asJson, _ = json.Marshal(map[string]string{})
	} else {
		asJson, err = json.Marshal(info)
		if err != nil {
			log.Fatal("could not marshal info for tag", fmt.Sprint(err))
		}
	}
	//asB64 := b64.URLEncoding.EncodeToString([]byte(asJson))
	return el.LogJSON(tag, string(asJson))
}
