// go:build windows
package main

import (
	"time"

	"golang.org/x/sys/windows"
)

func main() {
	for {
		message := windows.StringToUTF16Ptr("still running!")
		windows.MessageBox(0, message, windows.StringToUTF16Ptr("test"), windows.MB_OK)
		time.Sleep(30 * time.Second)
	}
}
