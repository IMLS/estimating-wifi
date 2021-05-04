VERSION = $(shell git describe --abbrev=0)
LDFLAGS = "-X gsa.gov/18f/version.Semver=$(VERSION)"

all: input-configuration session-counter wifi-hardware-search-cli

versioning:
	@echo $(VERSION) > version.txt

deps: versioning
	go mod download

input-configuration: deps
	CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 GOPATH=$$PWD/release go install -ldflags $(LDFLAGS) gsa.gov/18f/input-initial-configuration

session-counter: deps
	CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 GOPATH=$$PWD/release go install -ldflags $(LDFLAGS) gsa.gov/18f/session-counter

wifi-hardware-search-cli: deps
	CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 GOPATH=$$PWD/release go install -ldflags $(LDFLAGS) gsa.gov/18f/find-ralink

clean:
	chmod +w -R release/{bin,pkg}*
	rm -rf release/{bin,pkg}

test:
	go test ./internal/wifi-hardware-search/lshw/*.go
	go test ./internal/wifi-hardware-search/search/*.go
	go test ./cmd/input-configuration/*.go
	go test ./cmd/session-counting/*.go
	go test ./cmd/session-counting/api/*.go
	go test ./cmd/session-counting/config/*.go
