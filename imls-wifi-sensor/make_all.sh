#!/bin/sh

#VERSION := $(shell git describe --tags --abbrev=0)
#LDFLAGS="-X gsa.gov/18f/version.Semver=v3.5"
#ENVVARS='GOOS=linux GOARCH=arm GOARM=7'

#GOOS=linux GOARCH=arm GOARM=7 go install -ldflags -X "gsa.gov/18f/version.Semver=v3.5" gsa.gov/18f/cmd/linux-session-counter
#GOOS=linux GOARCH=arm GOARM=7 go install -ldflags -X "gsa.gov/18f/version.Semver=v3.5" gsa.gov/18f/cmd/wifi-hardware-search-cli
#GOOS=linux GOARCH=arm GOARM=7 go install gsa.gov/18f/cmd/linux-session-counter
#GOOS=linux GOARCH=arm GOARM=7 go install gsa.gov/18f/cmd/wifi-hardware-search-cli

go mod download
echo "FOOBAR"
#GOOS=linux GOARCH=arm GOARM=7 go install imls-wifi-sensor/cmd/linux-session-counter
#GOOS=linux GOARCH=arm GOARM=7 go install imls-wifi-sensor/cmd/wifi-hardware-search-cli
