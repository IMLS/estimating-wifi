module gsa.gov/18f/tempdb

go 1.16

replace gsa.gov/18f/wifi-hardware-search v0.0.0 => ../../internal/wifi-hardware-search

replace gsa.gov/18f/version v0.0.0 => ../../internal/version

replace gsa.gov/18f/config v0.0.0 => ../../internal/config

replace gsa.gov/18f/http v0.0.0 => ../../internal/http

replace gsa.gov/18f/cryptopasta v0.0.0 => ../../internal/cryptopasta

replace gsa.gov/18f/analysis v0.0.0 => ../../internal/analysis

replace gsa.gov/18f/logwrapper v0.0.0 => ../../internal/logwrapper

require (
	github.com/jmoiron/sqlx v1.3.4
	github.com/mattn/go-sqlite3 v1.14.7
	gsa.gov/18f/analysis v0.0.0
	gsa.gov/18f/config v0.0.0
	gsa.gov/18f/logwrapper v0.0.0
)
