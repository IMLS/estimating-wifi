module gsa.gov/18f/logwrapper

go 1.16

replace gsa.gov/18f/config v0.0.0 => ../../internal/config

replace gsa.gov/18f/cryptopasta v0.0.0 => ../../internal/cryptopasta

replace gsa.gov/18f/http v0.0.0 => ../../internal/http

replace gsa.gov/18f/wifi-hardware-search v0.0.0 => ../../internal/wifi-hardware-search

require (
	github.com/newrelic/go-agent/v3 v3.14.0
	github.com/newrelic/go-agent/v3/integrations/nrlogrus v1.0.1
	github.com/sirupsen/logrus v1.8.1
	gsa.gov/18f/config v0.0.0
	gsa.gov/18f/http v0.0.0
)
