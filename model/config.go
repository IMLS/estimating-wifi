package model

import "log"

func GetServer(cfg *Config, which string) *Server {
	log.Println("cfg:", cfg)
	log.Println("servers: ", cfg.Servers)
	for _, s := range cfg.Servers {
		log.Printf("model: considering: %v", s)
		if s.Name == which {
			return new(Server)
		}
	}
	log.Printf("model: could not retrieve server matching name '%v'", which)
	return nil
}

type Server struct {
	Name     string `yaml:"name"`
	Host     string `yaml:"host"`
	Authpath string `yaml:"authpath"`
	User     string `yaml:"user"`
	Pass     string `yaml:"pass"`
}

type Config struct {
	Monitoring struct {
		PingInterval          int `yaml:"pinginterval"`
		MaxHTTPErrorCount     int `yaml:"max_http_error_count"`
		HTTPErrorIntervalMins int `yaml:"http_error_interval_mins"`
		UniquenessWindow      int `yaml:"uniqueness_window"`
		DisconnectionWindow   int `yaml:"disconnection_window"`
	} `yaml:"monitoring"`
	Servers []struct {
		Name     string `yaml:"name"`
		Host     string `yaml:"host"`
		Authpath string `yaml:"authpath"`
		User     string `yaml:"user"`
		Pass     string `yaml:"pass"`
	} `yaml:"servers"`
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
