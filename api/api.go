package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"gsa.gov/18f/session-counter/constants"
	"gsa.gov/18f/session-counter/model"
)

// FUNC Mac_to_mfg
// Looks up a MAC address in the manufactuerer's database.
// Returns a valid name or "unknown" if the name cannot be found.
func Mac_to_mfg(cfg model.Config, mac string) string {
	db, err := sql.Open("sqlite3", cfg.Manufacturers.Db)
	if err != nil {
		log.Fatal("Failed to open manufacturer database.")
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
			if err != nil {
				log.Printf("Manufactuerer query failed: %s", q)
			} else {
				var id string

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

// FUNC Get_token
// Fetches a token from Directus for authenticating
// subsequent interactions with the service.
// Requires environment variables to be set
func Get_token(cfg model.Config) model.Token {
	user := os.Getenv(constants.EnvUsername)
	pass := os.Getenv(constants.EnvPassword)
	var uri string = (cfg.Server.Scheme + "://" +
		cfg.Server.Host + "/auth/login")

	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	reqBody, err := json.Marshal(map[string]string{
		"email":    user,
		"password": pass,
	})

	if err != nil {
		log.Fatal("Could not authenticate to Directus.")
	}

	req, err := http.NewRequest("POST", uri, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-type", "application/json")
	if err != nil {
		log.Fatal("Unable to construct URI for authentication.")
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error in client request to Directus /auth.")
	}
	// Closes the connection at function exit.
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Unable to read body of response from Directus /auth.")
	}
	res := model.Token{}
	json.Unmarshal(body, &res)

	return res
}

// FUNC bToMb
// Internal. Converts bytes to MB for reporting telemetry.
func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

// FUNC post_ram_usage
// Part of telemetry. Posts RAM usage.
func post_ram_usage(cfg model.Config, tok model.Token) {
	var uri string = (cfg.Server.Scheme + "://" +
		cfg.Server.Host + "/items/" +
		"memory_usage")

	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	all := fmt.Sprintf("alloc[%v], total[%v], sys[%v], numgc[%v]",
		bToMb(m.Alloc), bToMb(m.TotalAlloc),
		bToMb(m.Sys), m.NumGC)

	reqBody, err := json.Marshal(map[string]string{
		"bytes": strconv.Itoa(int(m.Alloc)),
		"notes": all,
	})

	if err != nil {
		log.Fatal("Failed to marshal RAM data to JSON.")
	}

	req, err := http.NewRequest("POST", uri, bytes.NewBuffer(reqBody))
	if err != nil {
		log.Fatal("Unable to generate request for RAM usage.")
	}

	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tok.Data.AccessToken))

	resp, err := client.Do(req)

	if err != nil {
		log.Fatal("Error in client POST for RAM usage.")
	}
	defer resp.Body.Close()

}

// FUNC post_manufactuerer_count
// Posts the manufactuerer count to Directus.
func post_manufactuerer_count(cfg model.Config, tok model.Token, e model.Entry) {
	var uri string = (cfg.Server.Scheme + "://" +
		cfg.Server.Host + "/items/" +
		cfg.Server.Collection)

	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	reqBody, err := json.Marshal(map[string]string{
		"mfgs":               e.Mfg,
		"mac":                e.MAC[0:8],
		"count":              strconv.Itoa(e.Count),
		"mfgl":               "not implemented",
		"libid":              "not implemented",
		"local_date_created": fmt.Sprintf(time.Now().Format(time.RFC3339)),
	})
	if err != nil {
		log.Fatal("Unable to marshal post of mfg count to JSON.")
	}

	req, err := http.NewRequest("POST", uri, bytes.NewBuffer(reqBody))
	if err != nil {
		log.Fatal("Unable to construct request for manufactuerer POST.")
	}
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tok.Data.AccessToken))

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Failure in client manufactuerer POST to Directus.")
	}
	// Close the body at function exit.
	defer resp.Body.Close()

	// We could process the result, but why?
}

func Report_telemetry(cfg model.Config, tok model.Token) {
	post_ram_usage(cfg, tok)
}

func Report_mfg(cfg model.Config, tok model.Token, e model.Entry) {
	post_manufactuerer_count(cfg, tok, e)
}
