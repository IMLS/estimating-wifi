package search

import (
	"embed"
	"encoding/json"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"gsa.gov/18f/internal/wifi-hardware-search/config"
	"gsa.gov/18f/internal/wifi-hardware-search/lshw"
	"gsa.gov/18f/internal/wifi-hardware-search/models"
)

// The text file is embedded at compile time.
// https://pkg.go.dev/embed#FS.ReadFile
//go:embed searches.json
var f embed.FS

// PURPOSE
// GetSearches attempts to read in the JSON document from the filesystem
// and use that, or it attempts to use the embedded version. The embedded version
// is used as a fallback in the case that we cannot find a (presumably tweaked/custom)
// set of searches in /etc...
func GetSearches() []models.Search {
	searches := make([]models.Search, 0)

	// Use the embedded file, which has a limited set of search terms.
	data, _ := f.ReadFile("searches.json")
	err := json.Unmarshal(data, &searches)

	if err != nil {
		log.Fatal("could not unmarshal search strings from embedded data.")
	}

	return searches
}

func SetMonitorMode(dev *models.Device) {
	cmds := make([]*exec.Cmd, 0)
	cmds = append(cmds, exec.Command("/usr/sbin/ip", "link", "set", dev.Logicalname, "down"))
	cmds = append(cmds, exec.Command("/usr/sbin/iw", dev.Logicalname, "set", "monitor", "none"))
	cmds = append(cmds, exec.Command("/usr/sbin/ip", "link", "set", dev.Logicalname, "up"))

	// Run the commands to set the adapter into monitor mode.
	for _, c := range cmds {
		c.Start()
		c.Wait()
	}
}

// PURPOSE
// Find any matching device. Returns the device structure
func SearchForMatchingDevice() *models.Device {
	dev := new(models.Device)
	dev.Exists = false
	for _, s := range GetSearches() {
		dev.Search = &s
		// findMatchingDevice populates device. Exits if something is found.
		FindMatchingDevice(dev)
		if dev.Exists {
			break
		}
	}
	return dev
}

// PURPOSE
// Takes a Device structure and, using the Search fields of that structure,
// attempts to find a matching WLAN device.
func FindMatchingDevice(wlan *models.Device) {
	// GetDeviceHash calls out to `lshw`.
	devices := lshw.GetDeviceHash(wlan)

	// We start by assuming that we have not found the device.
	found := false

	// Now, go through the devices and find the one that matches our criteria.
	for _, hash := range devices {

		// The default is to search all the fields
		if wlan.Search.Field == "ALL" {

			for k := range hash {
				// Lowercase everything for purposes of pattern matching.
				v, _ := regexp.MatchString(strings.ToLower(wlan.Search.Query), strings.ToLower(hash[k]))

				if v {
					// If we find it, set the fact that it exists. This will be picked up
					// back out in main() for the final act of printing a message to the user.
					wlan.Exists = true
				}
				// Stop searching if we find something.
				if wlan.Exists {
					break
				}
			}
		} else {
			// If we aren't doing a full search, then this is the alternative: check just
			// one field. It will still be a lowercase search, but it will be against one field only.
			v, _ := regexp.MatchString(strings.ToLower(wlan.Search.Query), strings.ToLower(hash[wlan.Search.Field]))
			if v {
				wlan.Exists = true
			}
		}

		// If we find something, proceed. But only keep the first thing we find.
		// Back in 'main', we'll handle the case where wlan.exists is false.
		if wlan.Exists && !found {
			found = true
			wlan.Vendor = strings.ToLower(hash["vendor"])
			wlan.Physicalid, _ = strconv.Atoi(hash["physical id"])
			wlan.Description = strings.ToLower(hash["description"])
			wlan.Businfo = strings.ToLower(hash["bus info"])
			wlan.Logicalname = strings.ToLower(hash["logical name"])
			wlan.Serial = strings.ToLower(hash["serial"])

			if len(hash["serial"]) >= config.MACLENGTH {
				wlan.Mac = strings.ToLower(hash["serial"][0:config.MACLENGTH])
			} else {
				wlan.Mac = strings.ToLower(hash["serial"])
			}
			wlan.Configuration = strings.ToLower(hash["configuration"])

			// Once we populate something in this loop, break out.
			// This will return us to the caller with a populated structure.
			break
		}
	}
}
