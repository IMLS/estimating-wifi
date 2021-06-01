module gsa.gov/18f/imls-data-convert

replace gsa.gov/18f/version v0.0.0 => ../../internal/version

replace gsa.gov/18f/analysis v0.0.0 => ../../internal/analysis

go 1.16

require (
	github.com/briandowns/spinner v1.12.0
	github.com/jszwec/csvutil v1.5.0
	github.com/mattn/go-sqlite3 v1.14.7
	gsa.gov/18f/analysis v0.0.0
	gsa.gov/18f/version v0.0.0
)
