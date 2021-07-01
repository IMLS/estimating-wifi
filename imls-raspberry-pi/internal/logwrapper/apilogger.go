package logwrapper

import (
	"fmt"
	"log"
	"time"

	"gsa.gov/18f/config"
	"gsa.gov/18f/http"
)

type ApiLogger struct {
	l   *StandardLogger
	cfg *config.Config
}

// type Writer interface {
//     Write(p []byte) (n int, err error)
// }

func NewApiLogger(cfg *config.Config) (api *ApiLogger) {
	api = &ApiLogger{cfg: cfg}
	return api
}

func (a *ApiLogger) Write(p []byte) (n int, err error) {
	fmt.Printf("API: %v\n", string(p))

	data := map[string]string{
		"pi_serial":   config.GetSerial(),
		"fcfs_seq_id": a.cfg.Auth.FCFSId,
		"device_tag":  a.cfg.Auth.DeviceTag,
		"session_id":  a.cfg.SessionId,
		"localtime":   time.Now().Format(time.RFC3339),
		"tag":         a.l.GetLogLevelName(),
		"info":        string(p),
	}

	_, err = http.PostJSON(a.cfg, a.cfg.GetLoggingUri(), []map[string]string{data})
	if err != nil {
		log.Println("could not log to API")
	}

	return len(p), nil
}
