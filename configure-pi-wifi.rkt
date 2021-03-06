#lang racket
(require "cpw.rkt")

(cpw)

(module+ test
  (require rackunit)
  
  (define test-string
    #<<TEST
  *-usb                     
       description: Wireless interface
       product: 802.11 n WLAN
       vendor: Ralink
       physical id: 1
       bus info: usb@1:1.1
       logical name: wlan1
       version: 1.01
       serial: 9c:ef:d5:fc:98:b7:00:00:00:00:00:00:00:00
       capabilities: usb-2.00 logical wireless
       configuration: broadcast=yes driver=rt2800usb driverversion=5.10.11-v7l+ firmware=0.36 link=yes maxpower=450mA multicast=yes speed=480Mbit/s wireless=IEEE 802.11
  *-network:0
       description: Ethernet interface
       physical id: 1
       logical name: eth0
       serial: dc:a6:32:d5:85:01
       capacity: 1Gbit/s
       capabilities: ethernet physical tp mii 10bt 10bt-fd 100bt 100bt-fd 1000bt 1000bt-fd autonegotiation
       configuration: autonegotiation=on broadcast=yes driver=bcmgenet driverversion=5.10.11-v7l+ link=no multicast=yes port=MII
  *-network:1
       description: Wireless interface
       physical id: 2
       logical name: wlan0
       serial: dc:a6:32:d5:85:03
       capabilities: ethernet physical wireless
       configuration: broadcast=yes driver=brcmfmac driverversion=7.45.229 firmware=01-2dbd9d2e ip=192.168.1.168 multicast=yes wireless=IEEE 802.11
  *-network:2
       description: Ethernet interface
       physical id: 3
       logical name: br-b3b3c2929cb6
       serial: 02:42:12:c7:b2:cf
       capabilities: ethernet physical
       configuration: autonegotiation=off broadcast=yes driver=bridge driverversion=2.3 firmware=N/A ip=172.18.0.1 link=no multicast=yes
  *-network:3
       description: Ethernet interface
       physical id: 4
       logical name: docker0
       serial: 02:42:de:7d:24:51
       capabilities: ethernet physical
       configuration: autonegotiation=off broadcast=yes driver=bridge driverversion=2.3 firmware=N/A ip=172.17.0.1 link=no multicast=yes
TEST
    )

  (check-pred hash? (find-ralink (regexp-split #px"\n" test-string)))
  )