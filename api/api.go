package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"gsa.gov/18f/session-counter/constants"
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
				fmt.Printf("Found mfg: %s\n", id)
				return id
			}
		}
	}

	return "unknown"
}

func get_token(cfg model.Config) model.Token {
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
	fmt.Println(res)
	return res
}

func post_manufactuerer_count(cfg model.Config, e model.Entry) {
	var uri string = (cfg.Server.Scheme + "://" +
		cfg.Server.Host + "/items/" +
		cfg.Server.Collection)

	var tok model.Token = get_token(cfg)

	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	fmt.Println("Reporting ", e.MAC, e.Mfg, e.Count)
	fmt.Println("mfg type: ", reflect.TypeOf(e.Mfg))

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
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Print("received data: ")
	fmt.Println(string(body[:]))
}

func Report_mfg(cfg model.Config, e model.Entry) {
	post_manufactuerer_count(cfg, e)
}
