// go:build windows
package main

import (
	"fmt"
	"os"
	"strconv"

	"golang.org/x/sys/windows"
	"gsa.gov/18f/internal/wifi-hardware-search/search"
)

func main() {
	title := windows.StringToUTF16Ptr("IMLS: compatible wifi device query")

	// Gives the user 5 tries to insert USB wifi adapter
	for i := 1; i <= 5; i++ {
		device := search.SearchForMatchingDevice()
		if device.Exists {
			message := windows.StringToUTF16Ptr(fmt.Sprintf("Found a compatible wifi device: %s (%s) [%s]",
				device.Logicalname,
				device.Description,
				device.Vendor))
			windows.MessageBox(0, message, title, windows.MB_OK)
			os.Exit(0)
		} else {
			numTriesLeft := 5 - i
			messageConcat := "No compatible wifi device was found. You have " + strconv.Itoa(numTriesLeft) + " tries left."
			message := windows.StringToUTF16Ptr(messageConcat)
			windows.MessageBox(0, message, title, windows.MB_ICONWARNING)
		}
	}
	errorMessage := windows.StringToUTF16Ptr("All tries used. Please submit a ticket on the repository. Aborting install.")
	windows.MessageBox(0, errorMessage, title, windows.MB_ICONWARNING)
	os.Exit(1)
}
