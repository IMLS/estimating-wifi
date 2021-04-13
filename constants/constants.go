package constants

const MACLENGTH = 17
const VERSION = "v0.0.2"

const SEARCHES_PATH = "/etc/session-counter/searches.json"

const (
	LOOKING_FOR_USB = iota
	READING_HASH    = iota
	DONE_READING    = iota
)

const LSHW_EXE = "/usr/bin/lshw"
