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

func Mac_to_mfg(cfg model.Config, mac string) string {
	db, err := sql.Open("sqlite3", cfg.Manufacturers.Db)
	if err != nil {
		log.Fatal("Failed to open manufacturer database.")
	}
	// Close the DB at the end of the function.
	defer db.Close()

	// We need to try the longest to the shortest MAC address
	// in order to match.
	// Start with aa:bb:cc:dd:ee
	lengths := []int{14, 11, 8}

	for _, length := range lengths {
		if len(mac) >= length {
			substr := mac[0:length]
			// FIXME: error handling
			q := fmt.Sprintf("SELECT id FROM oui WHERE mac LIKE %s", "'%"+substr+"'")

			rows, err := db.Query(q)
			if err != nil {
				log.Fatalf("Manufactuerer query failed: %s", q)
			}
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

	return "unknown"
}

func Get_token(cfg model.Config) model.Token {
	user := os.Getenv(constants.EnvUsername)
	pass := os.Getenv(constants.EnvPassword)
	var uri string = (cfg.Server.Scheme + "://" +
		cfg.Server.Host + "/auth/login")

	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	reqBody, _ := json.Marshal(map[string]string{
		"email":    user,
		"password": pass,
	})

	req, _ := http.NewRequest("POST", uri, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-type", "application/json")

	resp, _ := client.Do(req)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	res := model.Token{}
	json.Unmarshal(body, &res)

	return res
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

func post_ram_usage(cfg model.Config, tok model.Token) {
	var uri string = (cfg.Server.Scheme + "://" +
		cfg.Server.Host + "/items/" +
		"memory_usage")

	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	//fmt.Println("Reporting ", e.MAC, e.Mfg, e.Count)
	//fmt.Println("mfg type: ", reflect.TypeOf(e.Mfg))
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	all := fmt.Sprintf("alloc[%v], total[%v], sys[%v], numgc[%v]",
		bToMb(m.Alloc), bToMb(m.TotalAlloc),
		bToMb(m.Sys), m.NumGC)

	reqBody, _ := json.Marshal(map[string]string{
		"bytes": strconv.Itoa(int(m.Alloc)),
		"notes": all,
	})

	req, _ := http.NewRequest("POST", uri, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tok.Data.AccessToken))

	resp, _ := client.Do(req)
	defer resp.Body.Close()

}

func post_manufactuerer_count(cfg model.Config, tok model.Token, e model.Entry) {
	var uri string = (cfg.Server.Scheme + "://" +
		cfg.Server.Host + "/items/" +
		cfg.Server.Collection)

	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	//fmt.Println("Reporting ", e.MAC, e.Mfg, e.Count)
	//fmt.Println("mfg type: ", reflect.TypeOf(e.Mfg))

	reqBody, _ := json.Marshal(map[string]string{
		"mfgs":               e.Mfg,
		"mac":                e.MAC[0:8],
		"count":              strconv.Itoa(e.Count),
		"mfgl":               "not implemented",
		"libid":              "not implemented",
		"local_date_created": fmt.Sprintf(time.Now().Format(time.RFC3339)),
	})

	req, _ := http.NewRequest("POST", uri, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tok.Data.AccessToken))

	resp, _ := client.Do(req)
	defer resp.Body.Close()

	// Hopefully it doesn't matter if we do anything
	// with the body. However, this is a big FIXME.
	// body, _ := ioutil.ReadAll(resp.Body)

}

func Report_telemetry(cfg model.Config, tok model.Token) {
	post_ram_usage(cfg, tok)
}

func Report_mfg(cfg model.Config, tok model.Token, e model.Entry) {
	post_manufactuerer_count(cfg, tok, e)
}
