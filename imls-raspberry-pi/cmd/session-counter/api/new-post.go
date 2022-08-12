package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	resty "github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
	"gsa.gov/18f/cmd/session-counter/state"
	"gsa.gov/18f/internal/config"
)

type AuthSuccess struct {
	/* variables */
}
type AuthError struct {
	/* variables */
}

var timeOut int = 15

func asCSV(durations []*state.Duration) []byte {
	// TODO: import encoding/text
	return []byte{}
}

func PostDurations(durations []*state.Duration) error {

	fscs := config.GetFSCSID()
	uri := config.GetDurationsURI()
	key := config.GetAPIKey()

	log.Debug().Str("fscs", fscs).Str("uri", uri).Msg("sending")

	data := asCSV(durations)

	//***for testing***
	fmt.Println("Running Post")
	myJson, _ := json.Marshal(data)
	fmt.Println(string(myJson))

	//***create client and conditions needed for the client***
	client := resty.New()
	client.AddRetryCondition(
		func(r *resty.Response, err error) bool {
			return r.StatusCode() == http.StatusTooManyRequests
		},
	)
	client.SetTimeout(time.Duration(timeOut) * time.Second)

	// TODO: chunking in case we send more than 2MB data

	//***Client post data***
	resp, err := client.R().
		SetBody(data).
		SetAuthToken(key).
		//SetResult(&AuthSuccess{}). Could be incorperated once we have defined response
		//SetError(&AuthError{}). Could be incorperated once we have defined response
		Post(uri)

	if err != nil {
		log.Fatal().Err(err).Msg("could not send")
		return err
	} else {
		fmt.Println("  Body       :\n", resp)
	}

	//***for testing response code***
	//fmt.Println("Response Info:")
	//fmt.Println("Status Code:", resp.StatusCode())
	//fmt.Println("Status:", resp.Status())
	return nil
}
