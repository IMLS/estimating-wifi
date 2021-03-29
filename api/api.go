package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"gsa.gov/18f/session-counter/config"
	"gsa.gov/18f/session-counter/model"
)

// FUNC bToMb
// Internal. Converts bytes to MB for reporting telemetry.
func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

// FUNC post_ram_usage
// Part of telemetry. Posts RAM usage.
func Report_telemetry(cfg *config.Server, tok *model.Auth) (err error) {
	var uri string = (cfg.Host + "/items/memory_usage")

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
	var uri string = (cfg.Host + "/items/people2")

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
	// log.Printf("Using access token: %v\n", tok.Token)
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
