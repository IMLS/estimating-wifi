module gsa.gov/18f/config

replace gsa.gov/18f/cryptopasta v0.0.0 => ../cryptopasta

replace gsa.gov/18f/config v0.0.0 => ../config

replace gsa.gov/18f/logwrapper v0.0.0 => ../logwrapper

go 1.16

require (
	gopkg.in/yaml.v2 v2.4.0
	gsa.gov/18f/cryptopasta v0.0.0
)
