package model

type AuthError struct {
	// {"errors":[{"message":"Invalid user credentials.","extensions":{"code":"INVALID_CREDENTIALS"}}]}
	Errors []struct {
		Message    string `json:"message"`
		Extensions struct {
			Code string `json:"code"`
		} `json:"extensions"`
	} `json:"errors"`
}

type Auth struct {
	Token string `yaml:"token"`
	User  string `yaml:"username"`
}

type Authorization interface {
	GetToken() string
}

type RevalToken struct {
	AccessToken string `json:"token"`
}

func (rt RevalToken) GetToken() string {
	return rt.AccessToken
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

type Entry struct {
	MAC   string
	Mfg   string
	Count int
}

func GetToken(a Authorization) string {
	return a.GetToken()
}
