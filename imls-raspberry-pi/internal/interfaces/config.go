package interfaces

type Config interface {
	GetSerial() string

	GetFCFSSeqID() string
	GetDeviceTag() string
	GetAPIKey() string

	// TODO: change to generic set method?
	SetFCFSSeqID(string)
	SetDeviceTag(string)
	SetAPIKey(string)
	SetStorageMode(string)
	SetRunMode(string)
	SetUniquenessWindow(int)

	SetManufacturersPath(string)
	SetDurationsPath(string)
	SetQueuesPath(string)

	SetRootPath(string)
	SetImagesPath(string)

	GetLogLevel() string
	GetLoggers() []string
	// TODO: pull this out
	Log() Logger

	GetEventsURI() string
	GetDurationsURI() string

	IncrementSessionID() int64
	GetCurrentSessionID() int64

	IsStoringToAPI() bool
	IsStoringLocally() bool
	IsProductionMode() bool
	IsDeveloperMode() bool
	IsTestMode() bool

	GetManufacturersDatabase() Database
	GetDurationsDatabase() Database
	GetQueuesDatabase() Database

	GetWiresharkPath() string
	GetWiresharkDuration() int

	GetWWWRoot() string
	GetWWWImages() string

	GetMinimumMinutes() int
	GetMaximumMinutes() int
	GetUniquenessWindow() int
	GetResetCron() string
}

type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	SetLogLevel(level string)
}
