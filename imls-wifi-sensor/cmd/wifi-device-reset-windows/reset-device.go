// go:build windows
package main

import (
	"fmt"
	"os"

	"golang.org/x/sys/windows"
	"gsa.gov/18f/internal/wifi-hardware-search/search"
)

func main() {
	title := windows.StringToUTF16Ptr("Test adapter reset")
	device := search.SearchForMatchingDevice()
	if device.Exists {
		message := windows.StringToUTF16Ptr(fmt.Sprintf("Found a compatible wifi device: %s (%s) [%s]",
			device.Logicalname,
			device.Description,
			device.Vendor))
		windows.MessageBox(0, message, title, windows.MB_OK)
		os.Exit(0)
	} else {
		errorMessage := windows.StringToUTF16Ptr("No device found to reset")
		windows.MessageBox(0, errorMessage, title, windows.MB_ICONWARNING)
		os.Exit(1)
	}
}
