package model

type Entry struct {
	MAC   string
	Mfg   string
	Count int
}

type Token struct {
	Data struct {
		AccessToken  string `json:"access_token"`
		Expires      int    `json:"expires"`
		RefreshToken string `json:"refresh_token"`
	} `json:"data"`
}
