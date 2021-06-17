module gsa.gov/18f/imls-raspberry-pi

go 1.16

replace gsa.gov/18f/wifi-hardware-search-cli v0.0.0 => ./cmd/wifi-hardware-search-cli

replace gsa.gov/18f/input-initial-configuration v0.0.0 => ./cmd/input-configuration

replace gsa.gov/18f/session-counter v0.0.0 => ./cmd/session-counting

replace gsa.gov/18f/log-event v0.0.0 => ./cmd/log-event

replace gsa.gov/18f/version v0.0.0 => ./internal/version

replace gsa.gov/18f/config v0.0.0 => ./internal/config

replace gsa.gov/18f/http v0.0.0 => ./internal/http

replace gsa.gov/18f/cryptopasta v0.0.0 => ./internal/cryptopasta

replace gsa.gov/18f/wifi-hardware-search v0.0.0 => ./internal/wifi-hardware-search

replace gsa.gov/18f/analysis v0.0.0 => ./internal/analysis

require (
	github.com/acarl005/stripansi v0.0.0-20180116102854-5a71ef0e047d // indirect
	github.com/fatih/color v1.10.0 // indirect
	github.com/jszwec/csvutil v1.5.0 // indirect
	gsa.gov/18f/input-initial-configuration v0.0.0 // indirect
	gsa.gov/18f/log-event v0.0.0 // indirect
	gsa.gov/18f/session-counter v0.0.0 // indirect
	gsa.gov/18f/version v0.0.0 // indirect
	gsa.gov/18f/wifi-hardware-search v0.0.0 // indirect
	gsa.gov/18f/wifi-hardware-search-cli v0.0.0 // indirect
	gsa.gov/18f/analysis v0.0.0 // indirect
)
