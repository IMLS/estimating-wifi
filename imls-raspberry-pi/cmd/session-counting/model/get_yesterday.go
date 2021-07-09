package model

import (
	"fmt"
	"time"

	"gsa.gov/18f/logwrapper"
)

func GetYesterdaySessionId() string {
	lw := logwrapper.NewLogger(nil)
	yesterday := GetYesterday()
	yestersession := fmt.Sprintf("%v%02d%02d", yesterday.Year(), yesterday.Month(), yesterday.Day())
	lw.Debug("considering yesterday to be [", yestersession, "]")
	return yestersession
}

func GetYesterday() time.Time {
	yesterday := time.Now().Add(-24 * time.Hour)
	return yesterday
}
