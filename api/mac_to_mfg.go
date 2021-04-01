package api

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"gsa.gov/18f/session-counter/config"
)

func CheckMfgDatabaseExists(cfg *config.Config) {
	_, err := os.Stat(cfg.Manufacturers.Db)

	if os.IsNotExist(err) {
		log.Fatal("cannot find mfg database: ", cfg.Manufacturers.Db)
	}

	db, dberr := sql.Open("sqlite3", cfg.Manufacturers.Db)
	if err != nil {
		log.Println("Failed to open manufacturer database:", cfg.Manufacturers.Db)
		log.Fatal(dberr)
	}
	defer db.Close()
}

// FUNC Mac_to_mfg
// Looks up a MAC address in the manufactuerer's database.
// Returns a valid name or "unknown" if the name cannot be found.
func MacToMfg(cfg *config.Config, mac string) string {
	db, err := sql.Open("sqlite3", cfg.Manufacturers.Db)
	if err != nil {
		log.Fatal("Failed to open manufacturer database:", cfg.Manufacturers.Db)
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
			q := fmt.Sprintf("SELECT id FROM oui WHERE mac LIKE %s", "'"+substr+"%'")
			rows, err := db.Query(q)
			// Close the rows down, too...
			// Another possible leak?
			if err != nil {
				log.Println(err)
				log.Printf("manufactuerer not found: %s", q)
			} else {
				var id string

				defer rows.Close()

				for rows.Next() {
					err = rows.Scan(&id)
					if err != nil {
						log.Fatal("Failed in DB result row scanning.")
					}
					if id != "" {
						return id
					}
				}
			}
		}
	}

	return "unknown"
}
