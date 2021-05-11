.PHONY: all clean test

test:
	go test *.go

build:
	rm -f log-event
	go build log-event.go
	
all: test build
