// go:build windows
package main

import (
	"gsa.gov/18f/internal/wifi-hardware-search/netadapter"
)

func main() {
	netadapter.RestartNetAdapter()
}
