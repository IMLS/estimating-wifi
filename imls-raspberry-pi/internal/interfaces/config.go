package interfaces

import "github.com/benbjohnson/clock"

type Config interface {
	GetSerial() string
	GetFCFSSeqId() string
	GetDeviceTag() string
	GetAPIKey() string

	GetLogLevel() string
	GetLoggers() []string
	Log() Logger

	GetEventsUri() string
	GetDurationsUri() string

	InitializeSessionId()
	IncrementSessionId() int
	GetCurrentSessionId() int
	GetPreviousSessionId() int

	IsStoringToApi() bool
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
}
