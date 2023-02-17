// go:build windows
package main

import (
	"os"

	"gsa.gov/18f/internal/wifi-hardware-search/netadapter"
	"gsa.gov/18f/internal/wifi-hardware-search/search"
)

func main() {
	device := search.SearchForMatchingDevice()
	if device.Exists {
		netadapter.RestartNetAdapter(device.Logicalname)
		os.Exit(0)
	} else {
		os.Exit(1)
	}
}
