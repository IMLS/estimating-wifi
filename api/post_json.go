package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"gsa.gov/18f/session-counter/config"
	"gsa.gov/18f/session-counter/model"
)

func postJSON(svr *config.Server, tok *model.Auth, uri string, data map[string]string) error {
	log.Println("storing JSON to", uri)
	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
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
