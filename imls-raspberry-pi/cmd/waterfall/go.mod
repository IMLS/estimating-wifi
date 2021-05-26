module github.com/18f/waterfall

replace gsa.gov/18f/version v0.0.0 => ../../internal/version

replace gsa.gov/18f/imls-data-convert v0.0.0 => ../imls-data-convert

go 1.16

require (
	github.com/fogleman/gg v1.3.0 // indirect
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/jszwec/csvutil v1.5.0
	golang.org/x/image v0.0.0-20210504121937-7319ad40d33e // indirect
	gsa.gov/18f/imls-data-convert v0.0.0
)
