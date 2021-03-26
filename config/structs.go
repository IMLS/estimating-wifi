package config

import "log"

func GetServer(cfg *Config, which string) *Server {
	// log.Println("cfg:", cfg)
	// log.Println("servers: ", cfg.Servers)
	for _, s := range cfg.Servers {
		// log.Printf("config: considering: %v", s)
		if s.Name == which {
			return &s
		}
	}
	log.Printf("model: could not retrieve server matching name '%v'", which)
	return nil
}

type Server struct {
	Name     string `yaml:"name"`
	Host     string `yaml:"host"`
	Authpath string `yaml:"authpath"`
	Bearer   string `yaml:"bearer"`
	Postpath string `yaml:"postpath"`
	User     string `yaml:"user"`
	Pass     string `yaml:"pass"`
}

type Config struct {
	Monitoring struct {
		PingInterval          int `yaml:"pinginterval"`
		MaxHTTPErrorCount     int `yaml:"max_http_error_count"`
		HTTPErrorIntervalMins int `yaml:"http_error_interval_mins"`
		UniquenessWindow      int `yaml:"uniqueness_window"`
		Rounds                int `yaml:"rounds"`
		Threshold             int `yaml:"threshold"`
	} `yaml:"monitoring"`
	Servers   []Server `yaml:"servers"`
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
	Directus struct {
		Token string `yaml:"token"`
		User  string `yaml:"username"`
	} `yaml:"directus"`
	Reval struct {
		Token string `yaml:"token"`
		User  string `yaml:"username"`
	} `yaml:"reval"`
}
