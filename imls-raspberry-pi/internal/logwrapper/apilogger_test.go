package logwrapper

import (
	"testing"

	"gsa.gov/18f/config"
)

func TestSimple(t *testing.T) {
	lw := NewLogger(nil)
	lw.Error("HI")
}

func TestFileLogger(t *testing.T) {
	cfg := config.NewConfig()
	cfg.ReadConfig("./file-config.yaml")
	lw := UnsafeNewLogger(cfg)
	lw.SetLogLevel(INFO)
	lw.Info("Hi")
}

func TestApiLogger(t *testing.T) {
	cfg := config.NewConfig()
	cfg.ReadConfig("./api-config.yaml")
	lw := UnsafeNewLogger(cfg)
	lw.SetLogLevel(INFO)
	lw.Info("Hi")
}
