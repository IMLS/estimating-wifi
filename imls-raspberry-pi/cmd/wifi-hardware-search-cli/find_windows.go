// go:build windows
package main

import (
	"fmt"

	"golang.org/x/sys/windows"
	"gsa.gov/18f/internal/wifi-hardware-search/models"
	"gsa.gov/18f/internal/wifi-hardware-search/search"
)

func main() {
	device := new(models.Device)
	search.FindMatchingDevice(device)
	title := "compatible device query"
	if device.Exists {
		message := fmt.Sprintf("found: %s (%s) [%s]",
			device.Logicalname,
			device.Description,
			device.Vendor)
		windows.MessageBox(0, title, message, windows.MB_OK)
	} else {
		windows.MessageBox(0, title, "not found", windows.MB_ICONWARNING)
	}
}
