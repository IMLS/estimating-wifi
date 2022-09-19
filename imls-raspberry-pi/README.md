# Wifi Sensor

## Running session-counter
Assumes you are a developer with git and Golang installed

- Open your terminal to this repository
- cd to imls-raspberry-pi/cmd/session-counter
- Run `go run session-counter.go`
- To run `developer mode` instead of production, run `go run session-counter.go --mode dev`
    - Instead of running Wireshark, this runs fakeWireshark that creates fake MAC addresses for testing purposes
    - This is helpful for testing purposes when a developer does not have access to real sensors and devices