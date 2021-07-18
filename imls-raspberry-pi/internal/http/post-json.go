// Package http provides primitives around http communication.
package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"time"

	"gsa.gov/18f/internal/interfaces"
)

var slashWarned bool = false

func PostJSON(cfg interfaces.Config, uri string, data []map[string]interface{}) error {

	tok := cfg.GetAPIKey()
	matched, _ := regexp.MatchString(".*/$", uri)
	if !slashWarned && !matched {
		slashWarned = true
		log.Println("WARNING: api.data.gov wants a trailing slash on URIs")
	}

	timeout := time.Duration(15 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	// Lets not send too much data at once. So, we'll walk through the data array in steps of 20 elements.
	chunkSize := 20
	var divided [][]map[string]interface{}
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
			return errors.New("postjson: unable to marshal post of data to JSON")
		}

		// Next, it's time to create a request object. Again, fail if it doesn't work.
		req, err := http.NewRequest("POST", uri, bytes.NewBuffer(reqBody))
		if err != nil {
			return errors.New("postjson: unable to construct request for data POST")
		}

		// Lets set some headers. Perhaps we should quit here... but, we'll try and keep going
		// if anything fails.
		req.Header.Set("Content-type", "application/json")
		req.Header.Set("X-Api-Key", tok)

		// MAKE THE REQUEST
		resp, err := client.Do(req)
		// If there was an error in the post, log it, and exit the function.
		if err != nil {
			message := fmt.Sprintf("postjson: failure in client attempt to POST to %v", uri)
			log.Print(message)
			return fmt.Errorf(message)
		} else {
			// If we get things back, the errors will be encoded within the JSON.
			if resp.StatusCode < 200 || resp.StatusCode > 299 {
				message := fmt.Sprintf("PostJSON: bad status from POST to %v [%v]\n", uri, resp.Status)
				log.Print(message)
				return fmt.Errorf(message)
			} else {
				// Parse the response. Everything comes from ReVal in our current formulation.
				var dat RevalResponse
				body, _ := ioutil.ReadAll(resp.Body)
				err := json.Unmarshal(body, &dat)
				if err != nil {
					message := fmt.Sprintf("PostJSON: could not unmarshal response body: %v", err)
					log.Print(message)
					return fmt.Errorf(message)
				}
			}
		}
		// Close the body at function exit.
		defer resp.Body.Close()
	}

	return nil
}
