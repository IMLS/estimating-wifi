VERSION := $(shell git describe --tags --abbrev=0)
LDFLAGS = "-X gsa.gov/18f/version.Semver=$(VERSION)"

ifeq ($(shell uname -m),armv7l)
RELEASEPATH = ../release/bin/*
ENVVARS = CGO_ENABLED=1 GOOS=linux GOARCH=arm GOARM=7
else
RELEASEPATH = ../release/bin/linux_arm/*
ENVVARS = CGO_ENABLED=1 GOOS=linux
endif

check-install:
	go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.41.1

check:
	golangci-lint run ./...

all: input-configuration session-counter wifi-hardware-search-cli log-event cache-to-sqlite waterfall

deps:
	go mod download

input-configuration: deps
	${ENVVARS} GOPATH=$$PWD/../release go install -ldflags $(LDFLAGS) gsa.gov/18f/cmd/input-initial-configuration

session-counter: deps
	${ENVVARS} GOPATH=$$PWD/../release go install -ldflags $(LDFLAGS) gsa.gov/18f/cmd/session-counter

wifi-hardware-search-cli: deps
	${ENVVARS} GOPATH=$$PWD/../release go install -ldflags $(LDFLAGS) gsa.gov/18f/cmd/wifi-hardware-search-cli

log-event: deps
	${ENVVARS} GOPATH=$$PWD/../release go install -ldflags $(LDFLAGS) gsa.gov/18f/cmd/log-event

cache-to-sqlite: deps
	${ENVVARS} GOPATH=$$PWD/../release go install -ldflags $(LDFLAGS) gsa.gov/18f/cmd/cache-to-sqlite

waterfall: deps
	${ENVVARS} GOPATH=$$PWD/../release go install -ldflags $(LDFLAGS) gsa.gov/18f/cmd/waterfall

clean:
	mkdir -p ../release/{bin,pkg}
	chmod -R +w ../release/{bin,pkg}*
	rm -rf ../release/{bin,pkg}

test:
	go test -coverprofile all.out -timeout 45m ./...
