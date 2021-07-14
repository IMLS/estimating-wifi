package logwrapper

import (
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
	data := map[string]interface{}{
		"pi_serial":   a.cfg.GetSerial(),
		"fcfs_seq_id": a.cfg.Auth.FCFSId,
		"device_tag":  a.cfg.Auth.DeviceTag,
		"session_id":  a.cfg.SessionId,
		"localtime":   time.Now().Format(time.RFC3339),
		"tag":         a.l.GetLogLevelName(),
		"info":        string(p),
	}

	err = http.PostJSON(a.cfg, a.cfg.GetLoggingUri(), []map[string]interface{}{data})
	if err != nil {
		log.Println("could not log to API")
		log.Println(err.Error())
	}

	return len(p), nil
}
