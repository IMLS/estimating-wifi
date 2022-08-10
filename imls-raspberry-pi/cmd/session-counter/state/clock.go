package state

import (
	"github.com/benbjohnson/clock"
)

var clockSingleton *clock.Clock = nil

func GetClock() clock.Clock {
	if clockSingleton == nil {
		c := clock.New()
		clockSingleton = &c
	}
	return *clockSingleton
}

func SetClock(clock clock.Clock) {
	clockSingleton = &clock
}
