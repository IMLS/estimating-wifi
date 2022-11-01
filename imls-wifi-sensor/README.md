# Wifi Sensor

Assumes you are a developer with git and Golang installed

## Changing the Time to Send Data
Modify session-counter.ini in imls-wifi-sensor/cmd/linux-session-counter and/or imls-wifi-sensor/cmd/windows-session-counter

- If in prod:
  - [cron]
    `reset=0 0 ** *`
  - This sends data to the database every 24 hours
- If in dev:
  - [cron]
    `reset=*/5* ** *`
  - This sends data to the database every 5 minutes


## Running session-counter

- Open your terminal to this repository
- cd to imls-raspberry-pi/cmd/session-counter
- Run `go run session-counter.go`
- To run `developer mode` instead of production, run `go run session-counter.go --mode dev`
  - Instead of running Wireshark, this runs fakeWireshark that creates fake MAC addresses for testing purposes
  - This is helpful for testing purposes when a developer does not have access to real sensors and devices
- To run with real hardware on linux, you'll need some software installed:
    - `sudo add-apt-repository -y ppa:wireshark-dev/stable` 
    - `sudo apt install -y iw tshark`. 


### Possible errors

You may get a soft block error for your wifi device.

```
9:43AM FTL command failed error="exit status 2" command="/usr/sbin/ip link set wlx9cefd5fa48b4 up"
```

```
~/git/estimating-wifi/imls-wifi-sensor$ sudo rfkill list
0: phy0: Wireless LAN
        Soft blocked: yes
        Hard blocked: no
```

You could [follow these instructions](https://askubuntu.com/questions/62166/siocsifflags-operation-not-possible-due-to-rf-kill) to unblock depending on your host machine. (These instructions are for Ubuntu.)