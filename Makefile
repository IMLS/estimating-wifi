.PHONY: all

all: clean
	go build session-counter.go

clean:
	rm -f session-counter

install: all
	echo Copying executable into /usr/local/bin
	sudo cp session-counter /usr/local/bin/session-counter

test:
	go test *.go
	go test config/*.go
	go test api/*.go