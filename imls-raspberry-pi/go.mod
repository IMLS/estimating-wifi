module gsa.gov/18f/imls-raspberry-pi

go 1.16

replace gsa.gov/18f/find-ralink v0.0.0 => ./cmd/wifi-hardware-search-cli

replace gsa.gov/18f/input-initial-configuration v0.0.0 => ./cmd/input-configuration

replace gsa.gov/18f/session-counter v0.0.0 => ./cmd/session-counting

replace gsa.gov/18f/version v0.0.0 => ./internal/version

replace gsa.gov/18f/wifi-hardware-search v0.0.0 => ./internal/wifi-hardware-search

require (
	gopkg.in/yaml.v2 v2.4.0
	gsa.gov/18f/session-counter v0.0.0
)
