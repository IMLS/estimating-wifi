module gsa.gov/18f/logwrapper

go 1.16

replace gsa.gov/18f/config v0.0.0 => ../../internal/config

replace gsa.gov/18f/cryptopasta v0.0.0 => ../../internal/cryptopasta

require (
	github.com/sirupsen/logrus v1.8.1
	gsa.gov/18f/config v0.0.0
)
