// Package logwrapper wraps logging to various Write interfaces.
package logwrapper

import (
	"fmt"
	"log"
	"time"

	"gsa.gov/18f/internal/config"
	"gsa.gov/18f/internal/http"
)

type APILogger struct {
	l   *StandardLogger
	cfg *config.Config
}

func NewAPILogger(cfg *config.Config) (api *APILogger) {
	api = &APILogger{cfg: cfg}
	return api
}

func (a *APILogger) Write(p []byte) (n int, err error) {

	// work around a catch-22 where we are trying to log startup events without
	// a session.
	sessionID := a.cfg.SessionID
	// if a.cfg.SessionId != nil {
	// 	sessionID = a.cfg.SessionId.GetSessionId()
	// }

	data := map[string]interface{}{
		"pi_serial":   a.cfg.GetSerial(),
		"fcfs_seq_id": a.cfg.Auth.FCFSId,
		"device_tag":  a.cfg.Auth.DeviceTag,
		"session_id":  fmt.Sprint(sessionID),
		"localtime":   time.Now().Format(time.RFC3339),
		"tag":         a.l.GetLogLevelName(),
		"info":        string(p),
	}

	err = http.PostJSON(a.cfg, a.cfg.GetLoggingURI(), []map[string]interface{}{data})
	if err != nil {
		log.Println("could not log to API")
		log.Println(err.Error())
	}

	return len(p), nil
}
