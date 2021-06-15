module gsa.gov/18f/config

replace gsa.gov/18f/cryptopasta v0.0.0 => ../cryptopasta

replace gsa.gov/18f/config v0.0.0 => ../config

go 1.16

require (
	gopkg.in/yaml.v2 v2.4.0
	gsa.gov/18f/cryptopasta v0.0.0
)
