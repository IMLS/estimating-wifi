package api

import (
	"bytes"
	"encoding/csv"
	"net/http"
	"strconv"
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
	b := new(bytes.Buffer)
	w := csv.NewWriter(b)
	w.Write([]string{"start_time", "end_time"})
	for _, d := range durations {
		s := strconv.FormatInt(d.Start, 10)
		e := strconv.FormatInt(d.End, 10)
		w.Write([]string{s, e})
	}
	return b.Bytes()
}

func PostDurations(durations []*state.Duration) error {

	fscs := config.GetFSCSID()
	uri := config.GetDurationsURI()
	key := config.GetAPIKey()

	log.Debug().Str("fscs", fscs).Str("uri", uri).Msg("sending")

	data := asCSV(durations)

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
		log.Debug().Str("body", resp.Status()).Msg("received response")
	}

	//***for testing response code***
	//fmt.Println("Response Info:")
	//fmt.Println("Status Code:", resp.StatusCode())
	//fmt.Println("Status:", resp.Status())
	return nil
}
