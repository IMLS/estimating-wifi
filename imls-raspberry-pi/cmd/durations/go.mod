module github.com/18f/durations

replace gsa.gov/18f/version v0.0.0 => ../../internal/version

replace gsa.gov/18f/analysis v0.0.0 => ../../internal/analysis

replace gsa.gov/18f/config v0.0.0 => ../../internal/config

replace gsa.gov/18f/cryptopasta v0.0.0 => ../../internal/cryptopasta

go 1.16

require (
	github.com/barkimedes/go-deepcopy v0.0.0-20200817023428-a044a1957ca4 // indirect
	github.com/fogleman/gg v1.3.0
	github.com/jmoiron/sqlx v1.3.4
	github.com/jszwec/csvutil v1.5.0
	github.com/mattn/go-sqlite3 v1.14.7
	github.com/mxk/go-sqlite v0.0.0-20140611214908-167da9432e1f // indirect
	gsa.gov/18f/analysis v0.0.0
	gsa.gov/18f/config v0.0.0

)
