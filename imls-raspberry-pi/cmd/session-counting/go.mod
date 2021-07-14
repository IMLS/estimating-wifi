module gsa.gov/18f/session-counter

go 1.16

replace gsa.gov/18f/wifi-hardware-search v0.0.0 => ../../internal/wifi-hardware-search

replace gsa.gov/18f/version v0.0.0 => ../../internal/version

replace gsa.gov/18f/config v0.0.0 => ../../internal/config

replace gsa.gov/18f/http v0.0.0 => ../../internal/http

replace gsa.gov/18f/cryptopasta v0.0.0 => ../../internal/cryptopasta

replace gsa.gov/18f/analysis v0.0.0 => ../../internal/analysis

replace gsa.gov/18f/logwrapper v0.0.0 => ../../internal/logwrapper

replace gsa.gov/18f/structs v0.0.0 => ../../internal/structs

replace gsa.gov/18f/state v0.0.0 => ../../internal/state

require (
	github.com/benbjohnson/clock v1.1.0
	github.com/coreos/go-systemd v0.0.0-20191104093116-d3cd4ed1dbcf
	github.com/mattn/go-sqlite3 v1.14.7
	github.com/robfig/cron/v3 v3.0.0
	github.com/stretchr/testify v1.2.2
	golang.org/x/net v0.0.0-20210405180319-a5a99cb37ef4 // indirect
	golang.org/x/sys v0.0.0-20210510120138-977fb7262007 // indirect
	gsa.gov/18f/analysis v0.0.0
	gsa.gov/18f/config v0.0.0
	gsa.gov/18f/http v0.0.0
	gsa.gov/18f/logwrapper v0.0.0
	gsa.gov/18f/state v0.0.0
	gsa.gov/18f/structs v0.0.0
	gsa.gov/18f/version v0.0.0
	gsa.gov/18f/wifi-hardware-search v0.0.0
)
