package model

import (
	"time"

	"gsa.gov/18f/internal/session-counter-helper/state"
)

func GetYesterday() time.Time {
	offset := -24
	yesterday := state.GetClock().Now().In(time.Local).Add(time.Duration(offset) * time.Hour)
	return yesterday
}
