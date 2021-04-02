package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"gsa.gov/18f/session-counter/config"
)

// {
// 	"tables": [
// 	  {
// 		"headers": [
// 		  "event_id",
// 		  "device_uuid",
// 		  "lib_user",
// 		  "localtime",
// 		  "servertime",
// 		  "session_id",
// 		  "device_id"
// 		], // End headers
// 		"whole_table_errors": [],
// 		"rows": [
// 		  {
// 			"row_number": 2,
// 			"errors": [],
// 			"data": {
// 			  "event_id": "-1",
// 			  "device_uuid": "1000000089bbf88b",
// 			  "lib_user": "matthew.jadud@gsa.gov",
// 			  "localtime": "2021-04-02T10:46:53-04:00",
// 			  "servertime": "2021-04-02T10:46:53-04:00",
// 			  "session_id": "9475068c05fea81f",
// 			  "device_id": "unknown:6"
// 			} //end data
// 		  } // end row
// 		], // end rows
// 		"valid_row_count": 1,
// 		"invalid_row_count": 0
// 	  } // end table
// 	],
// 	"valid": true
//   }
type RevalResponse struct {
	Tables []struct {
		Headers          []string      `json:"headers"`
		WholeTableErrors []interface{} `json:"whole_table_errors"`
		Rows             []struct {
			RowNumber int               `json:"row_number"`
			Errors    []interface{}     `json:"errors"`
			Data      map[string]string `json:"data"`
		} `json:"rows"`
		ValidRowCount   int `json:"valid_row_count"`
		InvalidRowCount int `json:"invalid_row_count"`
	} `json:"tables"`
	Valid bool `json:"valid"`
}

// A package-local counter.
// The first thing we do is post an event. This will return a "magic index"
// or a foreign key, that we will use in our post of the data. This associates
// every piece of data entered with a session, and indexes the post in that session.
// That way, we can say "this set of data was entry 293 of session ABC."
// If it isn't an event object, we won't get a magic_index back, and it will
// be returned as -1. Hopefully, we'll be ignoring it in those cases...
var magic_index int = 0

func postJSON(cfg *config.Config, tok *config.AuthConfig, uri string, data []map[string]string) (int, error) {
	log.Println("postjson: storing JSON to", uri)
	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	var reqBody []byte
	var err error

	// FIXME
	// Directus takes posted data directly.
	// ReVal is currently looking for it to be "wrapped" in an array.
	// We should modify ReVal so that it takes the exact same POSTed data
	// as Directus, so that we cannot tell the difference from the client-side.
	// switch svr.Name {
	// case "directus":
	// 	reqBody, err = json.Marshal(data)
	// case "reval":
	// 	source := map[string][]map[string]string{"source": {data}}
	// 	reqBody, err = json.Marshal(source)
	// }

	// UPDATE 20210401 MCJ
	// We are now sending an object that has a single key, "source"
	// That is keyed to an array of objects. We're still sending a singleton.
	// But, it's wrapped in an object and an array.
	// We can no longer post directly to Directus.
	// source := map[string][]map[string]string{"source": {data}}
	// UPDATE 20210401 MCJ a little while later...
	// We don't need the source key, but we do need an array of objects.
	// arr := []map[string]string{data}
	reqBody, err = json.Marshal(data)

	if err != nil {
		return -1, errors.New("postjson: unable to marshal post of data to JSON")
	}

	req, err := http.NewRequest("POST", uri, bytes.NewBuffer(reqBody))
	if err != nil {
		return -1, errors.New("postjson: unable to construct request for data POST")
	}

	req.Header.Set("Content-type", "application/json")
	if tok != nil {
		log.Printf("Access token length: %v\n", len(tok.Token))
		req.Header.Set("X-Api-Key", tok.Token)
	} else {
		log.Printf("postjson: failed to set headers for authorization.")
	}

	// Clean up the log string... no tokens in the log
	reqLogString := strings.Replace(fmt.Sprint(req), tok.Token, "APITOKEN", -1)

	log.Printf("postjson:req:\n%v\n", reqLogString)
	resp, err := client.Do(req)
	log.Printf("postjson:resp: %v\n", resp)
	if err != nil {
		log.Printf("postjson:err resp: %v\n", resp)
		return -1, fmt.Errorf("api: failure in client attempt to POST to %v", uri)
	} else {
		// If we get things back, the errors will be encoded within the JSON.
		if resp.StatusCode < 200 || resp.StatusCode > 299 {
			log.Printf("postjson: bad status on POST to: %v\n", uri)
			log.Printf("postjson: bad status on POST response: [ %v ]\n", resp.Status)
		} else {
			// Bump the index on success
			magic_index += 1
			var dat RevalResponse
			body, _ := ioutil.ReadAll(resp.Body)
			err := json.Unmarshal(body, &dat)
			if err != nil {
				return magic_index, fmt.Errorf("postjson: could not unmarshal response body")
			}
			// 2021/03/26 14:00:18 resp.Body {"data":{"magic_index":12,"device_uuid":"1000000089bbf88b","lib_user":"10x@gsa.gov","session_id":"effc67d0068b4e7f","localtime":"2021-03-26T18:00:17Z","servertime":"2021-03-26T18:00:17Z","tag":"nil","info":"{}"}}
			// 20210420 MCJ UPDATE
			// We're now going through ReVal. This returns a different format...
			log.Println("postjson: resp.Body", string(body))
		}
	}
	// Close the body at function exit.
	defer resp.Body.Close()

	return magic_index, nil
}
