#!/bin/sh

#VERSION := $(shell git describe --tags --abbrev=0)
LDFLAGS = "-X gsa.gov/18f/version.Semver=v3.5"
ENVVARS = GOOS=linux GOARCH=arm GOARM=7
${ENVVARS} go install -ldflags $(LDFLAGS) gsa.gov/18f/cmd/linux-session-counter
${ENVVARS} go install -ldflags $(LDFLAGS) gsa.gov/18f/cmd/wifi-hardware-search-cli
