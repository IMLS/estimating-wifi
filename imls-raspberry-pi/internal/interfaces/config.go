package interfaces

import "github.com/benbjohnson/clock"

type Config interface {
	GetSerial() string
	GetFCFSSeqID() string
	GetDeviceTag() string
	GetAPIKey() string

	GetLogLevel() string
	GetLoggers() []string
	Log() Logger

	GetEventsURI() string
	GetDurationsURI() string

	InitializeSessionID()
	IncrementSessionID() int
	GetCurrentSessionID() int
	GetPreviousSessionID() int

	IsStoringToAPI() bool
	IsStoringLocally() bool
	IsProductionMode() bool
	IsDeveloperMode() bool
	IsTestMode() bool

	GetManufacturerDatabase() Database

	GetClock() clock.Clock
	GetMinimumMinutes() int
	GetMaximumMinutes() int
}

type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	SetLogLevel(level string)
}
