package api

import (
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

	if err != nil || resp.StatusCode() != http.StatusOK {
		log.Error().
			Err(err).
			Str("response", resp.String()).
			Msg("could not authenticate")
		return err
	}

	if json.Unmarshal(resp.Body(), &jwt) != nil {
		log.Error().
			Err(err).
			Str("response", resp.String()).
			Msg("could not unmarshal authentication response")
		return err
	}

	return nil
}

func PostDurations(durations []*state.Duration) error {
	token := JWTToken{}
	auth_err := PostAuthentication(&token)
	if auth_err != nil {
		return auth_err
	}

	fscs := config.GetFSCSID()
	uri := config.GetDurationsURI()

	// TODO: we need to chunk in case we send more than 2MB data
	client := resty.New()
	client.AddRetryCondition(
		func(r *resty.Response, err error) bool {
			return r.StatusCode() == http.StatusTooManyRequests
		},
	)
	client.SetTimeout(time.Duration(timeOut) * time.Second)

	for _, d := range durations {
		data := make(map[string]string)
		data["_start"] = strconv.FormatInt(d.Start, 10)
		data["_end"] = strconv.FormatInt(d.End, 10)
		data["_fscs"] = fscs

		resp, err := client.R().
			SetBody(data).
			SetAuthToken(token.Token).
			SetHeader("Content-Type", "application/json").
			SetError(&AuthError{}).
			Post(uri)

		if err != nil || resp.StatusCode() != http.StatusOK {
			log.Error().
				Err(err).
				Str("response", resp.String()).
				Msg("could not send durations")
			return err
		}
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

	if err != nil || resp.StatusCode() != http.StatusOK {
		log.Error().
			Err(err).
			Str("response", resp.String()).
			Msg("could not send heartbeat")
		return err
	}

	return nil
}
