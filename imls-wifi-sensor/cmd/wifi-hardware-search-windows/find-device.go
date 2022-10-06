// go:build windows
package main

import (
	"fmt"

	"golang.org/x/sys/windows"
	"gsa.gov/18f/internal/wifi-hardware-search/search"
)

func main() {
	device := search.SearchForMatchingDevice()
	title := windows.StringToUTF16Ptr("IMLS: compatible wifi device query")
	if device.Exists {
		message := windows.StringToUTF16Ptr(fmt.Sprintf("found a compatible wifi device: %s (%s) [%s]",
			device.Logicalname,
			device.Description,
			device.Vendor))
		windows.MessageBox(0, message, title, windows.MB_OK)
	} else {
		message := windows.StringToUTF16Ptr("no compatible wifi device was found")
		windows.MessageBox(0, message, title, windows.MB_ICONWARNING)
	}
}
