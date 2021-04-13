package constants

const MACLENGTH = 17
const VERSION = "v0.0.2"

const SEARCHES_PATH = "/etc/session-counter/searches.json"

const (
	LOOKING_FOR_SECTION_HEADING = iota
	READING_ENTRY               = iota
	DONE_WITH_SECTION           = iota
)

const LSHW_EXE = "/usr/bin/lshw"
