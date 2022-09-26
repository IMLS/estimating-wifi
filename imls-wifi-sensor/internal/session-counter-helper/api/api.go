package api

import (
	"bytes"
	"encoding/csv"
	"net/http"
	"strconv"
	"time"

	resty "github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
	"gsa.gov/18f/internal/config"
	"gsa.gov/18f/internal/session-counter-helper/state"
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
	w.Flush()
	return b.Bytes()
}

func PostDurations(durations []*state.Duration) error {

	// fscs := config.GetFSCSID()
	uri := config.GetDurationsURI()
	key := config.GetAPIKey()

	data := asCSV(durations)

	client := resty.New()
	client.AddRetryCondition(
		func(r *resty.Response, err error) bool {
			return r.StatusCode() == http.StatusTooManyRequests
		},
	)
	client.SetTimeout(time.Duration(timeOut) * time.Second)

	// TODO: chunking in case we send more than 2MB data

	resp, err := client.R().
		SetBody(data).
		SetAuthToken(key).
		SetHeader("Content-Type", "text/csv").
		//SetResult(&AuthSuccess{}). Could be incorperated once we have defined response
		//SetError(&AuthError{}). Could be incorperated once we have defined response
		Post(uri)

	if err != nil {
		log.Fatal().
			Err(err).
			Str("response", resp.String()).
			Msg("could not send")
		return err
	}

	return nil
}
