package config

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"time"

	yaml "gopkg.in/yaml.v2"
	"gsa.gov/18f/cryptopasta"
)

func NewConfig() *Config {
	cfg := Config{}
	cfg.setDefaults()
	return &cfg
}

func NewConfigFromPath(path string) (*Config, error) {
	cfg := Config{}
	cfg.setDefaults()
	err := cfg.ReadConfig(path)
	return &cfg, err
}

func (cfg *Config) ReadConfig(path string) error {
	_, err := os.Stat(path)
	// Stat will set an error if the file cannot be found.
	if err == nil {
		f, err := os.Open(path)
		if err != nil {
			log.Fatal("config: could not open configuration file. Exiting.")
		}
		defer f.Close()
		var newcfg *Config
		decoder := yaml.NewDecoder(f)
		err = decoder.Decode(&newcfg)
		if err != nil {
			log.Fatalf("config: could not decode YAML:\n%v\n", err)
		}
		*cfg = *newcfg
		// The API key will need to be decoded into memory.
		if len(cfg.Auth.Token) > 0 {
			cfg.Auth.Token = cfg.decodeAuthToken()
		}
		return nil
	} else {
		log.Printf("config: could not find config: %v\n", path)
	}
	return fmt.Errorf("config: could not find config file [%v]", path)
}

func (cfg *Config) NewSessionId() {
	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("%v", time.Now())))
	sid := fmt.Sprintf("%x", h.Sum(nil))[0:16]
	cfg.SessionId = sid
}

func (cfg *Config) GetLoggers() []string {
	return cfg.Loggers
}

func (cfg *Config) GetLogLevel() string {
	if cfg.LogLevel == "" {
		return "ERROR"
	} else {
		return cfg.LogLevel
	}
}

func (cfg *Config) setDefaults() {
	cfg.LogLevel = "INFO"
	cfg.Loggers = []string{"local:stderr", "local:tmp"}

	cfg.Monitoring.PingInterval = 30
	cfg.Monitoring.MaxHTTPErrorCount = 8
	cfg.Monitoring.HTTPErrorIntervalMins = 10
	cfg.Monitoring.UniquenessWindow = 24 * 60
	cfg.Monitoring.MinimumMinutes = 30
	cfg.Monitoring.MaximumMinutes = 600

	cfg.Umbrella.Scheme = "https"
	cfg.Umbrella.Host = "api.data.gov"
	cfg.Umbrella.Data = "/TEST/10x-imls/v2/durations/"
	cfg.Umbrella.Logging = "/TEST/10x-imls/v2/events/"

	cfg.Wireshark.Duration = 45
	cfg.Wireshark.Path = "/usr/bin/tshark"
	cfg.Wireshark.CheckWlan = "1"

	cfg.Manufacturers.Db = "/opt/imls/manufacturers.sqlite"

	cfg.StorageMode = "sqlite"

	cfg.Local.Crontab = "0 0 * * *"
	cfg.Local.SummaryDB = "/opt/imls/summary.sqlite"
	cfg.Local.TemporaryDB = "/tmp/imls.sqlite"
	cfg.Local.WebDirectory = "/www/imls"
}

func (cfg *Config) decodeAuthToken() string {
	// It is a B64 encoded string
	// of the API key encrypted with the device's serial.
	// This is obscurity, but it is all we can do on a RPi
	serial := []byte(cfg.Serial)
	var key [32]byte
	copy(key[:], serial)
	b64, err := base64.StdEncoding.DecodeString(cfg.Auth.Token)
	if err != nil {
		log.Println("config: cannot b64 decode auth token.")
	}
	dec, err := cryptopasta.Decrypt(b64, &key)
	if err != nil {
		log.Println("config: failed to decrypt auth token after decoding")
	}

	return string(dec)
}

type Config struct {
	Auth struct {
		Token     string `yaml:"api_token"`
		DeviceTag string `yaml:"device_tag"`
		FCFSId    string `yaml:"fcfs_seq_id"`
	} `yaml:"auth"`
	LogLevel     string `yaml:"log_level"`
	Loggers    []string `yaml:"loggers"`
	Monitoring struct {
		PingInterval          int `yaml:"pinginterval"`
		MaxHTTPErrorCount     int `yaml:"max_http_error_count"`
		HTTPErrorIntervalMins int `yaml:"http_error_interval_mins"`
		UniquenessWindow      int `yaml:"uniqueness_window"`
		MinimumMinutes        int `yaml:"minimum_minutes"`
		MaximumMinutes        int `yaml:"maximum_minutes"`
	} `yaml:"monitoring"`
	Umbrella struct {
		Scheme  string `yaml:"scheme"`
		Host    string `yaml:"host"`
		Data    string `yaml:"data"`
		Logging string `yaml:"logging"`
	} `yaml:"umbrella"`
	Wireshark struct {
		Duration  int    `yaml:"duration"`
		Adapter   string `yaml:"adapter"`
		Path      string `yaml:"path"`
		CheckWlan string `yaml:"check_wlan"`
	} `yaml:"wireshark"`
	Manufacturers struct {
		Db string `yaml:"db"`
	} `yaml:"manufacturers"`
	SessionId   string
	Serial      string `yaml:"serial"`
	StorageMode string `yaml:"storagemode"`
	Local       struct {
		Crontab      string `yaml:"crontab"`
		SummaryDB    string `yaml:"summary_db"`
		TemporaryDB  string `yaml:"temporary_db"`
		WebDirectory string `yaml:"web_directory"`
	} `yaml:"local"`
}
