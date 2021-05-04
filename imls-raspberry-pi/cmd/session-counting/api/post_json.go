package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
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

var slash_warned bool = false

func postJSON(cfg *config.Config, tok *config.AuthConfig, uri string, data []map[string]string) (int, error) {
	log.Println("postjson: storing JSON to", uri)

	matched, _ := regexp.MatchString(".*/$", uri)
	if !slash_warned && !matched {
		slash_warned = true
		log.Println("WARNING: api.data.gov wants a trailing slash on URIs")
	}

	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	// Lets not send too much data at once. So, we'll walk through the data array in steps of 20 elements.
	chunkSize := 20
	var divided [][]map[string]string
	for i := 0; i < len(data); i += chunkSize {
		end := i + chunkSize
		if end > len(data) {
			end = len(data)
		}

		divided = append(divided, data[i:end])
	}

	// FIXME: If there is no data, throw an event that logs there was no data.

	// Now that the incoming array has been chopped up into subarrays of length chunkSize,
	// lets send those out into the world.
	for _, arr := range divided {
		var reqBody []byte
		var err error

		// First, try marshalling the data.
		// We have to give up if this doesn't work.
		reqBody, err = json.Marshal(arr)
		if err != nil {
			return -1, errors.New("postjson: unable to marshal post of data to JSON")
		}

		// Next, it's time to create a request object. Again, fail if it doesn't work.
		req, err := http.NewRequest("POST", uri, bytes.NewBuffer(reqBody))
		if err != nil {
			return -1, errors.New("postjson: unable to construct request for data POST")
		}

		// Lets set some headers. Perhaps we should quit here... but, we'll try and keep going
		// if anything fails.
		req.Header.Set("Content-type", "application/json")
		if tok != nil {
			if config.Verbose {
				fmt.Printf("Access token length: %v\n", len(tok.Token))
			}
			req.Header.Set("X-Api-Key", tok.Token)
		} else {
			log.Printf("postjson: failed to set headers for authorization.")
		}

		// Clean up the log string... no tokens in the log
		// FIXME This makes a mess if there is no token in the `tok` structure...
		// APITOKEN&APITOKEN{APITOKENPAPITOKENOAPITOKENSAPITOKENTAPITOKEN A
		reqLogString := strings.Replace(fmt.Sprint(req), tok.Token, "APITOKEN", -1)
		if config.Verbose {
			log.Printf("postjson:req:\n%v\n", reqLogString)
		}

		// MAKE THE REQUEST
		resp, err := client.Do(req)
		// Show the response from the server. Helpful in debugging.
		if config.Verbose {
			log.Printf("postjson:resp: %v\n", resp)
		}
		// If there was an error in the post, log it, and exit the function.
		if err != nil {
			if config.Verbose {
				log.Printf("postjson:err resp: %v\n", resp)
			}
			return -1, fmt.Errorf("postjson: failure in client attempt to POST to %v", uri)
		} else {
			// If we get things back, the errors will be encoded within the JSON.
			if resp.StatusCode < 200 || resp.StatusCode > 299 {
				log.Printf("postjson: bad status on POST uri: %v\n", uri)
				log.Printf("postjson: bad status on POST response: [ %v ]\n", resp.Status)
				return magic_index, fmt.Errorf("postjson: bad status on POST response: [ %v ]", resp.Status)
			} else {
				// Parse the response. Everything comes from ReVal in our current formulation.
				var dat RevalResponse
				body, _ := ioutil.ReadAll(resp.Body)
				err := json.Unmarshal(body, &dat)
				if err != nil {
					// If we can't parse the response, return a valid index but also include an error.
					return magic_index, fmt.Errorf("postjson: could not unmarshal response body")
				}
				if config.Verbose {
					log.Println("postjson: resp.Body", string(body))
				}
			}
		}
		// Close the body at function exit.
		defer resp.Body.Close()
	}

	magic_index += 1
	return magic_index, nil
}
