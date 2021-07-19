package model

import (
	"time"

	"gsa.gov/18f/internal/interfaces"
)

func GetYesterday(cfg interfaces.Config) time.Time {
	offset := -24
	yesterday := cfg.GetClock().Now().Add(time.Duration(offset) * time.Hour)
	return yesterday
}
