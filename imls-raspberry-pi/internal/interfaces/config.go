package interfaces

type Config interface {
	GetSerial() string
	GetFCFSSeqId() string
	GetDeviceTag() string
	GetAPIKey() string

	GetLogLevel() string
	GetLoggers() []string

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
}
