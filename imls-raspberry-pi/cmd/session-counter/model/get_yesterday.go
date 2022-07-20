package model

import (
	"time"

	"gsa.gov/18f/cmd/session-counter/state"
	"gsa.gov/18f/internal/interfaces"
)

func GetYesterday(cfg interfaces.Config) time.Time {
	offset := -24
	yesterday := state.GetClock().Now().In(time.Local).Add(time.Duration(offset) * time.Hour)
	return yesterday
}
