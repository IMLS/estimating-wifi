package api

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"gsa.gov/18f/session-counter/model"
)

func Mac_to_mfg(cfg model.Config, mac string) string {
	// FIXME: error handling
	db, _ := sql.Open("sqlite3", cfg.Manufacturers.Db)
	// We need to try the longest to the shortest MAC address
	// in order to match.
	// Start with aa:bb:cc:dd:ee
	lengths := []int{14, 11, 8}

	for _, length := range lengths {
		substr := mac[0:length]
		// FIXME: error handling
		q := fmt.Sprintf("SELECT id FROM oui WHERE mac LIKE %s", "'%"+substr+"'")
		// fmt.Printf("query: %s\n", q)
		rows, _ := db.Query(q)
		var id string

		for rows.Next() {
			_ = rows.Scan(&id)
			if id != "" {
				// fmt.Printf("Found mfg: %s\n", id)
				return id
			}
		}
	}

	return "unknown"
}

func Report_mfg(mfg string, count int) {
	//resp, err := http.Post("https://jsonplaceholder.typicode.com/posts/1")

}
