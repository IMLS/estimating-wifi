// Package logwrapper wraps logging to various Write interfaces.
package logwrapper

import (
	"log"
	"time"

	"gsa.gov/18f/internal/http"
	"gsa.gov/18f/internal/interfaces"
)

type APILogger struct {
	l   *StandardLogger
	cfg interfaces.Config
}

func NewAPILogger(cfg interfaces.Config) (api *APILogger) {
	api = &APILogger{cfg: cfg}
	return api
}

func (a *APILogger) Write(p []byte) (n int, err error) {

	data := map[string]interface{}{
		"pi_serial":   a.cfg.GetSerial(),
		"fcfs_seq_id": a.cfg.GetFCFSSeqID(),
		"device_tag":  a.cfg.GetDeviceTag(),
		"session_id":  a.cfg.GetCurrentSessionID(),
		"localtime":   time.Now().Format(time.RFC3339),
		"tag":         a.l.GetLogLevelName(),
		"info":        string(p),
	}

	err = http.PostJSON(a.cfg, a.cfg.GetEventsURI(), []map[string]interface{}{data})
	if err != nil {
		log.Println("could not log to API")
		log.Println(err.Error())
	}

	return len(p), nil
}
