package lshw

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"
	"testing"
)

var t1 string = `
  *-usb
       description: Ethernet interface
       vendor: Standard Microsystems Corp.
       physical id: 1
       bus info: usb@1:1.1.1
       logical name: eth0
       version: 3.00
       serial: b8:27:eb:6e:aa:ad
       size: 1Gbit/s
       capacity: 1Gbit/s
       capabilities: usb-2.10 ethernet physical tp mii 10bt 10bt-fd 100bt 100bt-fd 1000bt-fd autonegotiation
       configuration: autonegotiation=on broadcast=yes driver=lan78xx driverversion=5.10.17-v7+ duplex=full ip=10.0.0.121 link=yes maxpower=2mA multicast=yes port=MII speed=1Gbit/s
  *-network:0
       description: Wireless interface
       physical id: 2
       logical name: wlan0
       serial: b8:27:eb:3b:ff:f8
       capabilities: ethernet physical wireless
       configuration: broadcast=yes driver=brcmfmac driverversion=7.45.229 firmware=01-2dbd9d2e ip=10.0.0.122 multicast=yes wireless=IEEE 802.11
  *-network:1
       description: Wireless interface
       physical id: 3
       bus info: usb@1:1.2
       logical name: wlan1
       serial: 00:c0:ca:89:19:60
       capabilities: ethernet physical wireless
       configuration: broadcast=yes driver=rt2800usb driverversion=5.10.17-v7+ firmware=0.36 ip=10.0.0.114 link=yes multicast=yes
wireless=IEEE 802.11
`

var t2 string = `
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
configuration: broadcast=yes driver=rt2800usb driverversion=5.10.17-v7l+ firmware=0.36 link=yes maxpower=450mA multicast=yes speed=480Mbit/s wireless=IEEE 802.11
*-network:0
description: Ethernet interface
physical id: 1
logical name: eth0
serial: dc:a6:32:d5:85:01
size: 1Gbit/s
capacity: 1Gbit/s
capabilities: ethernet physical tp mii 10bt 10bt-fd 100bt 100bt-fd 1000bt 1000bt-fd autonegotiation
configuration: autonegotiation=on broadcast=yes driver=bcmgenet driverversion=5.10.17-v7l+ duplex=full ip=192.168.1.166 link=yes multicast=yes port=MII speed=1Gbit/s
*-network:1
description: Wireless interface
physical id: 2
logical name: wlan0
serial: dc:a6:32:d5:85:03
capabilities: ethernet physical wireless
configuration: broadcast=yes driver=brcmfmac driverversion=7.45.229 firmware=01-2dbd9d2e multicast=yes wireless=IEEE 802.11
`

var t3 string = `*-usb                     
description: Wireless interface
vendor: Ralink
`

type hashtest struct {
	lshw    string
	hashndx int
	field   string
	regexp  string
}

func testHelper(lshwout string) []map[string]string {

	lines := make([]string, 0)
	scanner := bufio.NewScanner(strings.NewReader(lshwout))
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}
	result := ParseLSHW(lines)
	return result
}
func TestParse(t *testing.T) {

	var tests = []hashtest{
		{lshw: t3, hashndx: 0, field: "description", regexp: "Wireless"},
		{lshw: t2, hashndx: 1, field: "description", regexp: "Ethernet"},
	}

	for _, test := range tests {
		arrh := testHelper(test.lshw)
		s := arrh[test.hashndx][test.field]
		if ok, _ := regexp.MatchString(test.regexp, s); !ok {
			fmt.Println("test", test)
			t.Fatalf("Failed to match [%v] in [%v]", test.regexp, s)
		}

	}

}
