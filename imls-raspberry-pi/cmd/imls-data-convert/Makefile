.PHONY: test

test:
	go test *.go

build: outputs mac_x86 linux_x86

outputs:
	mkdir -p build

# TODO
# https://victoria.dev/blog/how-to-create-a-self-documenting-makefile/

mac_x86:
	 GOOS=darwin GOARCH=amd64 go build -o build/mac-idc
mac_arm:
	 GOOS=darwin GOARCH=arm64 go build -o build/armmac-idc
windows_x86:
	GOOS=windows GOARCH=amd64 go build -o build/win-idc.exe
linux_x86:
	GOOS=linux GOARCH=amd64 go build -o build/linux-idc
linux_rpi:
	GOOS=linux GOARCH=arm GOARM=7 go build -o build/rpi-idc

crossbuild: outputs mac_x86 mac_arm windows_x86 linux_x86 linux_rpi

all: test build
