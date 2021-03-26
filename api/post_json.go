package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"gsa.gov/18f/session-counter/config"
	"gsa.gov/18f/session-counter/model"
)

func postJSON(svr *config.Server, tok *model.Auth, uri string, data map[string]string) (int, error) {
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
		return -1, errors.New("api: unable to marshal post of mfg count to JSON")
	}

	req, err := http.NewRequest("POST", uri, bytes.NewBuffer(reqBody))
	if err != nil {
		return -1, errors.New("api: unable to construct request for manufactuerer POST")
	}

	req.Header.Set("Content-type", "application/json")
	log.Printf("Using access token: %v\n", tok.Token)
	req.Header.Set("Authorization", fmt.Sprintf("%s %s", svr.Bearer, tok.Token))

	log.Printf("req:\n%v\n", req)
	resp, err := client.Do(req)
	log.Printf("resp: %v\n", resp)

	magic_index := -1

	if err != nil {
		log.Printf("err resp: %v\n", resp)
		return -1, fmt.Errorf("api: failure in client manufactuerer POST to %v", svr.Name)
	} else {
		if resp.StatusCode < 200 || resp.StatusCode > 299 {
			log.Printf("api: bad status on POST to: %v\n", uri)
			log.Printf("api: bad status on POST response: [ %v ]\n", resp.Status)
		} else {
			var dat map[string]interface{}
			body, _ := ioutil.ReadAll(resp.Body)
			err := json.Unmarshal(body, &dat)
			if err != nil {
				return -1, fmt.Errorf("api: could not unmarshal response body")
			}
			// 2021/03/26 14:00:18 resp.Body {"data":{"magic_index":12,"device_uuid":"1000000089bbf88b","lib_user":"10x@gsa.gov","session_id":"effc67d0068b4e7f","localtime":"2021-03-26T18:00:17Z","servertime":"2021-03-26T18:00:17Z","tag":"nil","info":"{}"}}
			log.Println("resp.Body", string(body))
			if _, ok := dat["data"]; ok {
				inner := dat["data"].(map[string]interface{})
				if _, ok := inner["magic_index"]; ok {
					magic_index = int(inner["magic_index"].(float64))
				}
			}
		}
	}
	// Close the body at function exit.
	defer resp.Body.Close()

	return magic_index, nil
}
