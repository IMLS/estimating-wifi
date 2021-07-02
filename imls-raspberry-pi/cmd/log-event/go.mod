module gsa.gov/18f/log-event

go 1.16

replace gsa.gov/18f/config v0.0.0 => ../../internal/config

replace gsa.gov/18f/logwrapper v0.0.0 => ../../internal/logwrapper

replace gsa.gov/18f/version v0.0.0 => ../../internal/version

replace gsa.gov/18f/cryptopasta v0.0.0 => ../../internal/cryptopasta

replace gsa.gov/18f/http v0.0.0 => ../../internal/http

require (
	gsa.gov/18f/config v0.0.0
	gsa.gov/18f/logwrapper v0.0.0
	gsa.gov/18f/version v0.0.0
)
