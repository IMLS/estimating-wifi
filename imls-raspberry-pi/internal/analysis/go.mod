module gsa.gov/18f/analysis

go 1.16

replace gsa.gov/18f/config v0.0.0 => ../config

replace gsa.gov/18f/cryptopasta v0.0.0 => ../cryptopasta

replace gsa.gov/18f/wifi-hardware-search v0.0.0 => ../wifi-hardware-search

replace gsa.gov/18f/logwrapper v0.0.0 => ../logwrapper

replace gsa.gov/18f/http v0.0.0 => ../http

require (
	github.com/fogleman/gg v1.3.0
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	golang.org/x/image v0.0.0-20210622092929-e6eecd499c2c // indirect
	gsa.gov/18f/config v0.0.0
	gsa.gov/18f/logwrapper v0.0.0
)
