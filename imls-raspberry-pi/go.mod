module gsa.gov/18f/rpi-binaries

go 1.16

replace gsa.gov/18f/input-initial-configuration v0.0.0 => ./cmd/input-configuration
replace gsa.gov/18f/session-counter v0.0.0 => ./cmd/session-counting
replace gsa.gov/18f/find-ralink v0.0.0 => ./cmd/wifi-hardware-search-cli

require (
        gsa.gov/18f/input-initial-configuration v0.0.0
	gsa.gov/18f/session-counter v0.0.0
	gsa.gov/18f/find-ralink v0.0.0
)
