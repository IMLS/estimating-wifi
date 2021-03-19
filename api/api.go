package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"gsa.gov/18f/session-counter/config"
	"gsa.gov/18f/session-counter/model"
)

// FUNC Mac_to_mfg
// Looks up a MAC address in the manufactuerer's database.
// Returns a valid name or "unknown" if the name cannot be found.
func Mac_to_mfg(cfg *config.Config, mac string) string {
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

func unmarshalResponse(cfg *config.Server, body []byte) (tok *model.Auth, err error) {
	a := new(model.Auth)
	switch cfg.Name {
	case "directus":
		res := new(model.DirectusToken)
		err := json.Unmarshal(body, &res)

		if err != nil {
			log.Println(err)
			return nil, errors.New("api: error unmarshalling directus token")
		}
		a.User = model.GetUser(res)
		a.Token = model.GetToken(res)

		return a, nil
	case "reval":
		res := new(model.RevalToken)
		err := json.Unmarshal(body, &res)

		if err != nil {
			log.Println(err)
			return nil, errors.New("api: error unmarshalling reval token")
		}
		a.User = model.GetUser(res)
		a.Token = model.GetToken(res)
		return a, nil
	default:
		return nil, errors.New("api: no parser found for token response")
	}

	// Never get here
	// log.Fatal("Should never get here. API :micdrop:")
	// return nil, errors.New("api: catastrophic fail.")
}

// FUNC GetToken
// Fetches a token from Directus for authenticating
// subsequent interactions with the service.
// Requires environment variables to be set
func GetToken(cfg *config.Server) (tok *model.Auth, err error) {
	var uri string = ("https://" + cfg.Host + cfg.Authpath)

	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	auth, err := config.ReadAuth()
	if err != nil {
		log.Println("api: could not read auth from filesystem")
	}

	reqBody, err := json.Marshal(map[string]string{
		cfg.User: auth.User,
		cfg.Pass: auth.Token,
	})

	if err != nil {
		return nil, fmt.Errorf("api: could not marshal POST body for %v", cfg.Name)
	}

	req, err := http.NewRequest("POST", uri, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-type", "application/json")
	if err != nil {
		return nil, errors.New("api: unable to construct URI for authentication")
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("api: error in client request to %v %v", cfg.Name, cfg.Authpath)
	}
	// Closes the connection at function exit.
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("api: unable to read body of response from %v %v", cfg.Name, cfg.Authpath)
	}

	tok, err = unmarshalResponse(cfg, body)

	return tok, err
}

// FUNC bToMb
// Internal. Converts bytes to MB for reporting telemetry.
func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

// FUNC post_ram_usage
// Part of telemetry. Posts RAM usage.
func Report_telemetry(cfg *config.Server, tok *model.Auth) (err error) {
	var uri string = ("https://" + cfg.Host + "/items/memory_usage")

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
		return errors.New("api: ailed to marshal RAM data to JSON")
	}

	req, err := http.NewRequest("POST", uri, bytes.NewBuffer(reqBody))
	if err != nil {
		return errors.New("api: unable to generate request for RAM usage")
	}

	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tok.Token))

	resp, err := client.Do(req)

	if err != nil {
		return errors.New("api: error in client POST for RAM usage")
	}
	defer resp.Body.Close()
	return nil
}

// FUNC post_manufactuerer_count
// Posts the manufactuerer count to Directus.
func Report_mfg(cfg *config.Server, tok *model.Auth, e model.Entry) (err error) {
	var uri string = ("https://" + cfg.Host + "/items/people2")

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
		"local_date_created": time.Now().Format(time.RFC3339),
	})
	if err != nil {
		return errors.New("api: unable to marshal post of mfg count to JSON")
	}

	req, err := http.NewRequest("POST", uri, bytes.NewBuffer(reqBody))
	if err != nil {
		return errors.New("api: unable to construct request for manufactuerer POST")
	}
	req.Header.Set("Content-type", "application/json")
	// log.Printf("Using access token: %v\n", tok.Data.AccessToken)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tok.Token))

	// log.Printf("req:\n%v\n", req)
	resp, err := client.Do(req)
	// log.Printf("resp: %v\n", resp)
	// FIXME
	// If we fail to auth, it won't be a failed POST.
	// we'll get back a resp object with a 401
	/*
		2021/03/16 10:57:32 resp: &{
			401
			Unauthorized 401 HTTP/2.0 2 0
			map[
				Content-Length:[96]
				Content-Type:[application/json; charset=utf-8]
				Date:[Tue, 16 Mar 2021 14:57:32 GMT]
				Etag:[W/"60-SpvBqFAbsdy4SkXwsevzfPClFZA"]
				Strict-Transport-Security:[max-age=31536000]
				Vary:[Origin]
				X-Content-Type-Options:[nosniff]
				X-Frame-Options:[DENY]
				X-Powered-By:[Directus]
				X-Vcap-Request-Id:[f6273432-1e44-45fa-7950-54ca7d0eac47]
				X-Xss-Protection:[1; mode=block]
				] 0x2649b10 96 [] false false map[] 0x2528700 0x25e0120
			}
	*/
	if err != nil {
		log.Printf("err resp: %v\n", resp)
		return errors.New("api: failure in client manufactuerer POST to Directus")
	} else {
		if resp.StatusCode < 200 || resp.StatusCode > 299 {
			log.Printf("api: bad status on POST to: %v\n", uri)
			log.Printf("api: bad status on POST response: [ %v ]\n", resp.Status)
		}
	}
	// Close the body at function exit.
	defer resp.Body.Close()

	return nil
}
