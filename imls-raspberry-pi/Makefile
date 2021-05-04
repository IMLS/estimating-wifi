all: input-configuration session-counter wifi-hardware-search-cli

deps:
	go mod download

input-configuration: deps
	GOPATH=$$PWD/release go install gsa.gov/18f/input-initial-configuration

session-counter: deps
	GOPATH=$$PWD/release go install gsa.gov/18f/session-counter

wifi-hardware-search-cli: deps
	GOPATH=$$PWD/release go install gsa.gov/18f/find-ralink

# not working at present.
test:
	go test ./internal/wifi-hardware-search/lshw/*.go
	go test ./internal/wifi-hardware-search/search/*.go
