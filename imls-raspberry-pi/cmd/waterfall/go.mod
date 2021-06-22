module github.com/18f/waterfall

replace gsa.gov/18f/version v0.0.0 => ../../internal/version

replace gsa.gov/18f/analysis v0.0.0 => ../../internal/analysis

replace gsa.gov/18f/config v0.0.0 => ../../internal/config

replace gsa.gov/18f/cryptopasta v0.0.0 => ../../internal/cryptopasta

go 1.16

require (
	github.com/fogleman/gg v1.3.0
	github.com/jmoiron/sqlx v1.3.4
	github.com/jszwec/csvutil v1.5.0
	github.com/mattn/go-sqlite3 v1.14.7
	gsa.gov/18f/analysis v0.0.0
	gsa.gov/18f/config v0.0.0

)
