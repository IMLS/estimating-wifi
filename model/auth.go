package model

import (
	"os"

	"gsa.gov/18f/session-counter/constants"
)

// Located at /etc/sc/auth.yaml
type Auth struct {
	Token string `yaml:"token"`
	User  string `yaml:"username"`
}

type Authorization interface {
	GetToken() string
	GetUser() string
}

type RevalToken struct {
	AccessToken string `json:"token"`
}

func (rt RevalToken) GetToken() string {
	return rt.AccessToken
}

func getUserFromEnv() string {
	user := os.Getenv(constants.EnvUsername)
	return user
}

func (rt RevalToken) GetUser() string {
	return getUserFromEnv()
}

type DirectusToken struct {
	Data struct {
		AccessToken  string `json:"access_token"`
		Expires      int    `json:"expires"`
		RefreshToken string `json:"refresh_token"`
	} `json:"data"`
}

func (dt DirectusToken) GetToken() string {
	return dt.Data.AccessToken
}

func (dt DirectusToken) GetUser() string {
	return getUserFromEnv()
}

type Entry struct {
	MAC   string
	Mfg   string
	Count int
}

func GetUser(a Authorization) string {
	return a.GetUser()
}

func GetToken(a Authorization) string {
	return a.GetToken()
}
