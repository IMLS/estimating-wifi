package model

import (
	"fmt"
	"time"

	"gsa.gov/18f/config"
	"gsa.gov/18f/logwrapper"
)

func GetYesterdaySessionId(cfg *config.Config) string {
	lw := logwrapper.NewLogger(nil)
	yesterday := GetYesterday(cfg)
	yestersession := fmt.Sprintf("%v%02d%02d", yesterday.Year(), yesterday.Month(), yesterday.Day())
	lw.Debug("considering yesterday to be [", yestersession, "]")
	return yestersession
}

func GetYesterday(cfg *config.Config) time.Time {
	offset := -24
	// Mocking the clock... now, time should work correctly.
	// if cfg.IsTestMode() {
	// 	offset = -1
	// }
	yesterday := cfg.Clock.Now().Add(time.Duration(offset) * time.Hour)
	return yesterday
}
