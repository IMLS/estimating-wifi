package config

// At /etc/session-counter/config.yaml
type Config struct {
	Monitoring struct {
		PingInterval          int `yaml:"pinginterval"`
		MaxHTTPErrorCount     int `yaml:"max_http_error_count"`
		HTTPErrorIntervalMins int `yaml:"http_error_interval_mins"`
		UniquenessWindow      int `yaml:"uniqueness_window"`
		Rounds                int `yaml:"rounds"`
		Threshold             int `yaml:"threshold"`
	} `yaml:"monitoring"`
	Umbrella struct {
		Scheme  string `yaml:"scheme"`
		Host    string `yaml:"host"`
		Data    string `yaml:"data"`
		Logging string `yaml:"logging"`
	} `yaml:"umbrella"`
	Wireshark struct {
		Duration int    `yaml:"duration"`
		Adapter  string `yaml:"adapter"`
		Path     string `yaml:"path"`
	} `yaml:"wireshark"`
	Manufacturers struct {
		Db string `yaml:"db"`
	} `yaml:"manufacturers"`
	SessionId string
	Serial    string
}

// Located at /etc/session-counter/auth.yaml
type AuthConfig struct {
	Umbrella struct {
		Token string `yaml:"token"`
		Email string `yaml:"email"`
	}
}
