// go:build windows
package main

import (
	"os"

	"golang.org/x/sys/windows"
	"gsa.gov/18f/internal/wifi-hardware-search/netadapter"
	"gsa.gov/18f/internal/wifi-hardware-search/search"
)

func main() {
	title := windows.StringToUTF16Ptr("Test adapter reset")
	device := search.SearchForMatchingDevice()
	if device.Exists {
		netadapter.RestartNetAdapter(device.Logicalname)
		windows.MessageBox(0, device.Logicalname, title, windows.MB_OK)
		os.Exit(0)
	} else {
		errorMessage := windows.StringToUTF16Ptr("No device found to reset")
		windows.MessageBox(0, errorMessage, title, windows.MB_ICONWARNING)
		os.Exit(1)
	}
}
