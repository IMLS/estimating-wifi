package model

import (
	"fmt"
	"time"

	"gsa.gov/18f/logwrapper"
)

func GetYesterdaySessionId() string {
	lw := logwrapper.NewLogger(nil)
	yesterday := time.Now().Add(-1 * time.Hour)
	yestersession := fmt.Sprintf("%v%02d%02d", yesterday.Year(), yesterday.Month(), yesterday.Day())
	lw.Debug("considering yesterday to be [", yestersession, "]")
	return yestersession
}
