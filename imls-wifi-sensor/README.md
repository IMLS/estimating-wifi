# Wifi Sensor

Assumes you are a developer with git and Golang installed

## Changing the Time to Send Data
Modify session-counter.ini in imls-wifi-sensor/cmd/linux-session-counter and/or imls-wifi-sensor/cmd/windows-session-counter

- If in prod:
  - [cron]
    reset=0 0 * * *
  - This sends data to the database every 24 hours
- If in dev:
  - [cron]
    reset=*/5 * * * *
  - This sends data to the database every 5 minutes


## Running session-counter

- Open your terminal to this repository
- cd to imls-raspberry-pi/cmd/session-counter
- Run `go run session-counter.go`
- To run `developer mode` instead of production, run `go run session-counter.go --mode dev`
  - Instead of running Wireshark, this runs fakeWireshark that creates fake MAC addresses for testing purposes
  - This is helpful for testing purposes when a developer does not have access to real sensors and devices