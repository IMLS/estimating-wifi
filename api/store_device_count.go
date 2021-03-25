package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"gsa.gov/18f/session-counter/config"
	"gsa.gov/18f/session-counter/model"
)

// FUNC StoreDeviceCount
// Stores the device count JSON to directus and reval.
func StoreDeviceCount(cfg *config.Config, svr *config.Server, tok *model.Auth, uid string, count int) error {
	var uri string = ("https://" + svr.Host + svr.Postpath)
	log.Println("storing to", uri)
	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	data := map[string]string{
		"mfgs":               uid,
		"mac":                uid,
		"count":              strconv.Itoa(count),
		"mfgl":               "not implemented",
		"libid":              "not implemented",
		"local_date_created": time.Now().Format(time.RFC3339),
	}

	// Either
	// reqBody, err := json.Marshal(data)

	var reqBody []byte
	var err error
	switch svr.Name {
	case "directus":
		reqBody, err = json.Marshal(data)
	case "reval":
		source := map[string][]map[string]string{"source": {data}}
		reqBody, err = json.Marshal(source)
	}

	if err != nil {
		return errors.New("api: unable to marshal post of mfg count to JSON")
	}

	req, err := http.NewRequest("POST", uri, bytes.NewBuffer(reqBody))
	if err != nil {
		return errors.New("api: unable to construct request for manufactuerer POST")
	}

	req.Header.Set("Content-type", "application/json")
	log.Printf("Using access token: %v\n", tok.Token)
	req.Header.Set("Authorization", fmt.Sprintf("%s %s", svr.Bearer, tok.Token))

	log.Printf("req:\n%v\n", req)
	resp, err := client.Do(req)
	log.Printf("resp: %v\n", resp)
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
		return fmt.Errorf("api: failure in client manufactuerer POST to %v", svr.Name)
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
