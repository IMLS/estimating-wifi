// Package api provides auxilary functions for API integration
package api

import (
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"gsa.gov/18f/internal/interfaces"
)

// FUNC Mac_to_mfg
// Looks up a MAC address in the manufacturer's database.
// Returns a valid name or "unknown" if the name cannot be found.

// We hit this *all the time*. Perhaps we can speed it up?
var cache map[string]string = make(map[string]string)

func MacToMfg(cfg interfaces.Config, mac string) string {
	db := cfg.GetManufacturersDatabase()
	// We need to try the longest to the shortest MAC address
	// in order to match.
	// Start with aa:bb:cc:dd:ee
	// ... then   aa:bb:cc:dd
	// ... then   aa:bb:cc
	lengths := []int{14, 11, 8}
	for _, length := range lengths {
		// If we're given a short MAC address, don't
		// try and slice more of the string than exists.
		if len(mac) >= length {
			substr := mac[0:length]
			// Can we pull this out of the "memoized" cache?
			if v, ok := cache[substr]; ok {
				return v
			} else {
				q := fmt.Sprintf("SELECT id FROM oui WHERE mac LIKE %s", "'"+substr+"%'")
				rows, err := db.Query(q)
				if err != nil {
					cfg.Log().Debug("manufacturer not found: ", q)
					cfg.Log().Debug(err.Error())
				} else {
					var id string

					defer rows.Close()

					for rows.Next() {
						err = rows.Scan(&id)
						if err != nil {
							cfg.Log().Fatal("failed in DB result row scanning")
						}
						if id != "" {
							cache[substr] = id
							return id
						}
					}
				}
			}

		}
	}

	// If we got here, then it doesn't matter which subset we use...
	// we don't recognize this MAC address.
	for _, length := range lengths {
		if len(mac) >= length {
			substr := mac[0:length]
			cache[substr] = "unknown"
		}
	}

	return "unknown"
}
