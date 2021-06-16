module gsa.gov/18f/input-initial-configuration

go 1.16

replace gsa.gov/18f/config v0.0.0 => ../../internal/config

replace gsa.gov/18f/cryptopasta v0.0.0 => ../../internal/cryptopasta

replace gsa.gov/18f/version v0.0.0 => ../../internal/version

require (
	gopkg.in/yaml.v2 v2.4.0
	github.com/acarl005/stripansi v0.0.0-20180116102854-5a71ef0e047d
	github.com/fatih/color v1.10.0
	golang.org/x/sys v0.0.0-20201119102817-f84b799fce68 // indirect
	gsa.gov/18f/config v0.0.0
	gsa.gov/18f/version v0.0.0
)
