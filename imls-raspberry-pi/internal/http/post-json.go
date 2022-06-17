// Package http provides primitives around http communication.
package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"

	"github.com/rs/zerolog/log"
	"gsa.gov/18f/internal/state"
)

var slashWarned bool = false

func PostJSON(uri string, data []map[string]interface{}) error {

	if state.IsStoringLocally() {
		// do nothing.
		return nil
	}

	tok := state.GetAPIKey()
	matched, _ := regexp.MatchString(".*/$", uri)
	if !slashWarned && !matched {
		slashWarned = true
		log.Warn().Msg("missing a trailing slash on URIs")
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

	if len(data) == 0 {
		return errors.New("PostJSON: no events found")
	}

	// Now that the incoming array has been chopped up into subarrays of length chunkSize,
	// lets send those out into the world.
	for _, arr := range divided {
		var reqBody []byte
		var err error

		// First, try marshalling the data.
		// We have to give up if this doesn't work.
		reqBody, err = json.Marshal(arr)
		if err != nil {
			return errors.New("PostJSON: unable to marshal post of data to JSON")
		}

		// Next, it's time to create a request object. Again, fail if it doesn't work.
		req, err := http.NewRequest("POST", uri, bytes.NewBuffer(reqBody))
		if err != nil {
			return errors.New("PostJSON: unable to construct request for data POST")
		}

		req.Header.Set("Content-type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tok))
		resp, err := client.Do(req)

		if err != nil {
			message := fmt.Sprintf("PostJSON: failure in client attempt to POST to %v", uri)
			log.Warn().Str("uri", uri).Msg(message)
			return fmt.Errorf(message)
		}

		// If we get things back, the errors will be encoded within the JSON.
		if resp.StatusCode < 200 || resp.StatusCode > 299 {
			message := fmt.Sprintf("PostJSON: bad status from POST to %v [%v]\n", uri, resp.Status)
			body, _ := ioutil.ReadAll(resp.Body)
			log.Warn().Str("message", string(body)).Str("response", resp.Status).Msg(message)
			return fmt.Errorf(message)
		} else {
			// Parse the response. Everything comes from ReVal in our current formulation.
			var dat RevalResponse
			body, _ := ioutil.ReadAll(resp.Body)

			err := json.Unmarshal(body, &dat)
			if err != nil {
				message := fmt.Sprintf("PostJSON: could not unmarshal response body: %v", err)
				log.Warn().Err(err).Str("body", string(body)).Msg(message)
				return fmt.Errorf(message)
			}
		}
		// Close the body at function exit.
		defer resp.Body.Close()
	}

	return nil
}
