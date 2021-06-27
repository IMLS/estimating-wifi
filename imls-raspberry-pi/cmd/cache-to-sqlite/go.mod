module gsa.gov/18f/cache-to-sqlite

replace gsa.gov/18f/version v0.0.0 => ../../internal/version

replace gsa.gov/18f/analysis v0.0.0 => ../../internal/analysis

replace gsa.gov/18f/config v0.0.0 => ../../internal/config

replace gsa.gov/18f/cryptopasta v0.0.0 => ../../internal/cryptopasta

go 1.16

require (
	github.com/briandowns/spinner v1.12.0
	github.com/kr/pretty v0.1.0 // indirect
	github.com/mattn/go-sqlite3 v1.14.7
	golang.org/x/sys v0.0.0-20200223170610-d5e6a3e2c0ae // indirect
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127 // indirect
	gsa.gov/18f/analysis v0.0.0
	gsa.gov/18f/version v0.0.0
)
