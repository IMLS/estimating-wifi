package search

import (
	"embed"
	"encoding/json"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
	"gsa.gov/18f/internal/config"
	"gsa.gov/18f/internal/wifi-hardware-search/lshw"
	"gsa.gov/18f/internal/wifi-hardware-search/models"
	"gsa.gov/18f/internal/wifi-hardware-search/netadapter"
)

// This is used for truncating longer MAC addresses
// into a standard/32-bit form.
const MACLENGTH = 17

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
		log.Fatal().
			Err(err).
			Msg("could not unmarshal search strings from embedded data")
	}

	return searches
}

func SetMonitorMode(dev *models.Device) {
	cmds := make([]*exec.Cmd, 0)
	if runtime.GOOS == "windows" {
		cmds = append(cmds, exec.Command(config.GetWlanHelperPath(), dev.Logicalname, "mode", "monitor"))
	} else {
		cmds = append(cmds, exec.Command(config.GetIpPath(), "link", "set", dev.Logicalname, "down"))
		cmds = append(cmds, exec.Command(config.GetIwPath(), dev.Logicalname, "set", "monitor", "none"))
		cmds = append(cmds, exec.Command(config.GetIpPath(), "link", "set", dev.Logicalname, "up"))
	}
	// Run the commands to set the adapter into monitor mode.
	for _, c := range cmds {
		err := c.Start()
		if err != nil {
			log.Error().
				Err(err).
				Str("command", c.String()).
				Msg("command did not execute")
		}
		err = c.Wait()
		if err != nil {
			log.Fatal().
				Err(err).
				Str("command", c.String()).
				Msg("command failed")
		}
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
		findMatchingDevice(dev)
		if dev.Exists {
			break
		}
	}
	return dev
}

func SearchForMatchingDeviceWithQuery(field string, query string) *models.Device {
	dev := new(models.Device)
	dev.Exists = false
	s := &models.Search{Field: field, Query: query}
	dev.Search = s
	findMatchingDevice(dev)
	return dev
}

func osFindMatchingDevice(wlan *models.Device) []map[string]string {
	if runtime.GOOS == "windows" {
		return netadapter.GetDeviceHash(wlan)
	} else {
		// GetDeviceHash calls out to `lshw`.
		return lshw.GetDeviceHash(wlan)
	}
}

// PURPOSE
// Takes a Device structure and, using the Search fields of that structure,
// attempts to find a matching WLAN device.
func findMatchingDevice(wlan *models.Device) {
	devices := osFindMatchingDevice(wlan)

	// We start by assuming that we have not found the device.
	found := false

	// Now, go through the devices and find the one that matches our criteria.
	for _, hash := range devices {
		log.Debug().Msg("--------")
		for k, v := range hash {
			log.Debug().Str("key", k).Str("value", v).Msg("looking at field")
		}

		// The default is to search all the fields
		if wlan.Search.Field == "ALL" {

			for k := range hash {
				// Lowercase everything for purposes of pattern matching.
				v, _ := regexp.MatchString(strings.ToLower(wlan.Search.Query), strings.ToLower(hash[k]))
				log.Debug().
					Str("query", wlan.Search.Query).
					Str("field", wlan.Search.Field).
					Msg("searching all fields")
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
			log.Debug().
				Str("query", wlan.Search.Query).
				Str("field", wlan.Search.Field).
				Msg("searching one field")
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

			if len(hash["serial"]) >= MACLENGTH {
				wlan.Mac = strings.ToLower(hash["serial"][0:MACLENGTH])
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
