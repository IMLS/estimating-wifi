package model

import (
	"time"

	"gsa.gov/18f/internal/interfaces"
	"gsa.gov/18f/internal/state"
)

func GetYesterday(cfg interfaces.Config) time.Time {
	offset := -24
	yesterday := state.GetClock().Now().Add(time.Duration(offset) * time.Hour)
	return yesterday
}
