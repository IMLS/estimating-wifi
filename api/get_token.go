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

func unmarshalResponse(cfg *config.Server, authcfg *config.AuthConfig, body []byte) (tok *model.Auth, err error) {
	a := new(model.Auth)
	switch cfg.Name {
	case "directus":
		res := new(model.DirectusToken)
		err := json.Unmarshal(body, &res)

		if err != nil {
			log.Println(err)
			return nil, errors.New("api: error unmarshalling directus token")
		}

		log.Println("directus token: ", res)
		a.User = authcfg.Directus.User
		a.Token = model.GetToken(res)

		return a, nil
	case "reval":
		res := new(model.RevalToken)
		err := json.Unmarshal(body, &res)

		if err != nil {
			log.Println(err)
			return nil, errors.New("api: error unmarshalling reval token")
		}
		a.User = authcfg.Reval.User
		a.Token = model.GetToken(res)
		return a, nil
	default:
		return nil, errors.New("api: no parser found for token response")
	}

	// Never get here
	// log.Fatal("Should never get here. API :micdrop:")
	// return nil, errors.New("api: catastrophic fail.")
}

// FUNC GetToken
// Fetches a token from a service for authenticating
// subsequent interactions with the service.
// Requires environment variables to be set
func GetTokenStub(cfg *config.Server) (tok *model.Auth, err error) {
	return nil, nil
}

func GetToken(cfg *config.Server) (tok *model.Auth, err error) {
	var uri string = (cfg.Host + cfg.Authpath)

	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	// The auth structure contains a username and password
	// for each of the services.
	authcfg, err := config.ReadAuth()
	auth := new(model.Auth)

	if err != nil {
		log.Println("api: could not read auth from filesystem")
	}

	if cfg.Name == "directus" {
		auth.User = authcfg.Directus.User
		auth.Token = authcfg.Directus.Token
	} else {
		auth.User = authcfg.Reval.User
		auth.Token = authcfg.Reval.Token
	}

	reqBody, err := json.Marshal(map[string]string{
		cfg.User: auth.User,
		cfg.Pass: auth.Token,
	})

	if err != nil {
		return nil, fmt.Errorf("api: could not marshal POST body for %v", cfg.Name)
	}

	req, err := http.NewRequest("POST", uri, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-type", "application/json")
	log.Println("req: ", req)

	if err != nil {
		return nil, errors.New("api: unable to construct URI for authentication")
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("api: error in client request to %v %v", cfg.Name, cfg.Authpath)
		log.Println(err)
		return nil, fmt.Errorf("api: error in client request to %v %v", cfg.Name, cfg.Authpath)
	}
	// Closes the connection at function exit.
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	log.Println("resp: ", string(body))
	// FIXME
	// Handle error conditions better.
	// Directus
	// 2021/03/19 12:10:52 resp:  {"errors":[{"message":"Invalid user credentials.","extensions":{"code":"INVALID_CREDENTIALS"}}]}
	// Reval
	// {"non_field_errors":["Unable to log in with provided credentials."]}

	if err != nil {
		return nil, fmt.Errorf("api: unable to read body of response from %v %v", cfg.Name, cfg.Authpath)
	}

	errRes := new(model.AuthError)
	errR := json.Unmarshal(body, &errRes)
	if errR != nil {
		log.Println("api: could not unmarshal error response")
		log.Println("user ", auth.User, " token ", auth.Token)
		log.Println(uri)
		log.Println(errR)
		log.Println("api: err message ", errRes.Errors)
		for _, e := range errRes.Errors {
			log.Println("api: err code ", e.Message)
		}
	}

	tok, err = unmarshalResponse(cfg, authcfg, body)

	return tok, err
}
