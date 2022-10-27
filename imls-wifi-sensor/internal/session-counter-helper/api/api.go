package api

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	resty "github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
	"gsa.gov/18f/internal/config"
	"gsa.gov/18f/internal/session-counter-helper/state"
)

type JWTToken struct {
	Token string `json:"token"`
}

// postgrest error response
type AuthError struct {
	code    string
	details string
	hint    string
	message string
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

func PostAuthentication(jwt *JWTToken) error {
	fscs := config.GetFSCSID()
	key := config.GetAPIKey()

	client := resty.New()
	client.AddRetryCondition(
		func(r *resty.Response, err error) bool {
			return r.StatusCode() == http.StatusTooManyRequests
		},
	)
	client.SetTimeout(time.Duration(timeOut) * time.Second)

	login_data := make(map[string]string)
	login_data["fscs_id"] = fscs
	login_data["api_key"] = key

	login := config.GetLoginURI()
	resp, err := client.R().
		SetBody(login_data).
		SetHeader("Content-Type", "application/json").
		SetError(&AuthError{}).
		Post(login)

	if err != nil {
		log.Fatal().
			Err(err).
			Str("response", resp.String()).
			Msg("could not authenticate")
		return err
	}

	if json.Unmarshal(resp.Body(), &jwt) != nil {
		log.Fatal().
			Err(err).
			Str("response", resp.String()).
			Msg("could not unmarshal authentication response")
		return err
	}

	return nil
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
		SetError(&AuthError{}).
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

func PostHeartBeat() error {
	token := JWTToken{}
	auth_err := PostAuthentication(&token)
	if auth_err != nil {
		return auth_err
	}

	fscs := config.GetFSCSID()
	serial := state.GetCachedSerial()
	uri := config.GetHeartbeatURI()
	data := make(map[string]string)
	data["_fscs_id"] = fscs
	data["_sensor_version"] = "1.0"
	data["_sensor_serial"] = serial

	client := resty.New()
	client.AddRetryCondition(
		func(r *resty.Response, err error) bool {
			return r.StatusCode() == http.StatusTooManyRequests
		},
	)
	client.SetTimeout(time.Duration(timeOut) * time.Second)
	resp, err := client.R().
		SetBody(data).
		SetAuthToken(token.Token).
		SetHeader("Content-Type", "application/json").
		SetError(&AuthError{}).
		Post(uri)

	if err != nil || resp.StatusCode() != 200 {
		log.Fatal().
			Err(err).
			Str("response", resp.String()).
			Msg("could not send")
		return err
	}

	return nil
}
