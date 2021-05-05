module gsa.gov/18f/session-counter

go 1.16

replace gsa.gov/18f/wifi-hardware-search v0.0.0 => ../../internal/wifi-hardware-search

replace gsa.gov/18f/version v0.0.0 => ../../internal/version

require (
	github.com/coreos/go-systemd v0.0.0-20191104093116-d3cd4ed1dbcf
	github.com/mattn/go-sqlite3 v1.14.6
	gopkg.in/yaml.v2 v2.4.0
	gsa.gov/18f/version v0.0.0
	gsa.gov/18f/wifi-hardware-search v0.0.0
)
