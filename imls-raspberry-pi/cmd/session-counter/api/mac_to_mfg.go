package api

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"gsa.gov/18f/internal/config"
	"gsa.gov/18f/internal/logwrapper"
)

func CheckMfgDatabaseExists(cfg *config.Config) {
	lw := logwrapper.NewLogger(nil)

	_, err := os.Stat(cfg.Manufacturers.Db)

	if os.IsNotExist(err) {
		lw.Fatal("cannot find mfg database: ", cfg.Manufacturers.Db)
	}

	db, dberr := sql.Open("sqlite3", cfg.Manufacturers.Db)
	if err != nil {
		lw.Info("failed to open manufacturer database: ", cfg.Manufacturers.Db)
		lw.Fatal(dberr.Error())
	}
	defer db.Close()
}

// FUNC Mac_to_mfg
// Looks up a MAC address in the manufactuerer's database.
// Returns a valid name or "unknown" if the name cannot be found.

// We hit this *all the time*. Perhaps we can speed it up?
var cache map[string]string = make(map[string]string)

func MacToMfg(cfg *config.Config, mac string) string {
	lw := logwrapper.NewLogger(nil)
	db, err := sql.Open("sqlite3", cfg.Manufacturers.Db)
	if err != nil {
		lw.Fatal("Failed to open manufacturer database: ", cfg.Manufacturers.Db)
	}
	// Close the DB at the end of the function.
	// If not, it's a resource leak.
	defer db.Close()

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
				// Close the rows down, too...
				// Another possible leak?
				if err != nil {
					lw.Debug("manufactuerer not found: ", q)
					lw.Debug(err.Error())
				} else {
					var id string

					defer rows.Close()

					for rows.Next() {
						err = rows.Scan(&id)
						if err != nil {
							lw.Fatal("failed in DB result row scanning")
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
