module gsa.gov/18f/session-counter

go 1.16

replace gsa.gov/18f/wifi-hardware-search v0.0.0 => ../../internal/wifi-hardware-search

replace gsa.gov/18f/version v0.0.0 => ../../internal/version

replace gsa.gov/18f/config v0.0.0 => ../../internal/config

replace gsa.gov/18f/http v0.0.0 => ../../internal/http

replace gsa.gov/18f/cryptopasta v0.0.0 => ../../internal/cryptopasta

replace gsa.gov/18f/analysis v0.0.0 => ../../internal/analysis

replace gsa.gov/18f/logwrapper v0.0.0 => ../../internal/logwrapper

require (
	github.com/barkimedes/go-deepcopy v0.0.0-20200817023428-a044a1957ca4 // indirect
	github.com/coreos/go-systemd v0.0.0-20191104093116-d3cd4ed1dbcf
	github.com/jmoiron/sqlx v1.3.4
	github.com/mattn/go-sqlite3 v1.14.6
	github.com/robfig/cron/v3 v3.0.0
	golang.org/x/tools v0.1.4 // indirect
	gsa.gov/18f/analysis v0.0.0
	gsa.gov/18f/config v0.0.0
	gsa.gov/18f/http v0.0.0
	gsa.gov/18f/logwrapper v0.0.0
	gsa.gov/18f/version v0.0.0
	gsa.gov/18f/wifi-hardware-search v0.0.0
)
