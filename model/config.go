package model

type Config struct {
	Monitoring struct {
		PingInterval          int `yaml:"pinginterval"`
		MaxHTTPErrorCount     int `yaml:"max_http_error_count"`
		HTTPErrorIntervalMins int `yaml:"http_error_interval_mins"`
		UniquenessWindow      int `yaml:"uniqueness_window"`
		DisconnectionWindow   int `yaml:"disconnection_window"`
	} `yaml:"monitoring"`
	Server struct {
		Scheme     string `yaml:"scheme"`
		Host       string `yaml:"host"`
		Collection string `yaml:"collection"`
	} `yaml:"server"`
	Wireshark struct {
		Duration  int    `yaml:"duration"`
		Rounds    int    `yaml:"rounds"`
		Threshold int    `yaml:"threshold"`
		Adapter   string `yaml:"adapter"`
		Path      string `yaml:"path"`
	} `yaml:"wireshark"`
	Manufacturers struct {
		Db string `yaml:"db"`
	} `yaml:"manufacturers"`
	SessionId string
}
