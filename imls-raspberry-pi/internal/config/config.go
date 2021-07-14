package config

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"reflect"
	"regexp"
	"strings"

	"github.com/benbjohnson/clock"
	yaml "gopkg.in/yaml.v2"
	"gsa.gov/18f/cryptopasta"
	"gsa.gov/18f/wifi-hardware-search/config"
)

const STATEDB = "state.sqlite"

func NewConfig() *Config {
	cfg := Config{}
	cfg.setDefaults()
	//cfg.Validate()
	return &cfg
}

func NewConfigFromPath(path string) (*Config, error) {
	cfg := Config{}
	cfg.setDefaults()
	err := cfg.ReadConfig(path)
	if err != nil {
		log.Println("cfg could not be read")
		log.Fatal()
	}
	log.Println(cfg)
	if cfg.Clock == nil {
		log.Println("clock should not be nil")
		log.Fatal()
	}
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
			cfg.DecodeSerial()
			cfg.Auth.Token = cfg.decodeAuthToken()
		}
		if len(cfg.LshwPath) > 0 {
			config.SetLSHWLocation(cfg.LshwPath)
		}

		// Need to reset the clock pointer...
		// Gets wiped out by the read.
		cfg.Clock = clock.New()

		// Validate the config before returning.
		//cfg.Validate()

		return nil
	} else {
		log.Printf("config: could not find config: %v\n", path)
	}
	return fmt.Errorf("config: could not find config file [%v]", path)
}

func (cfg *Config) NewSessionId() {
	// h := sha256.New()
	// h.Write([]byte(fmt.Sprintf("%v", time.Now())))
	// sid := fmt.Sprintf("%x", h.Sum(nil))[0:16]
	// cfg.SessionId = sid
	t := cfg.Clock.Now()
	cfg.SessionId = fmt.Sprintf("%v%02d%02d", t.Year(), t.Month(), t.Day())
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

func (cfg *Config) IsStoringToApi() bool {
	return strings.Contains(strings.ToLower(cfg.StorageMode), "api")
}

func (cfg *Config) IsStoringLocally() bool {
	either := false
	for _, s := range []string{"local", "sqlite"} {
		either = either || strings.Contains(strings.ToLower(cfg.StorageMode), s)
	}
	return either
}

func (cfg *Config) IsProductionMode() bool {
	return strings.Contains(strings.ToLower(cfg.RunMode), "prod")
}

func (cfg *Config) IsDeveloperMode() bool {
	either := false
	for _, s := range []string{"dev", "test"} {
		either = either || strings.Contains(strings.ToLower(cfg.RunMode), s)
	}
	if either {
		log.Println("running in developer mode")
	}
	return either
}

func (cfg *Config) IsTestMode() bool {
	return strings.Contains(strings.ToLower(cfg.RunMode), "test")
}

func (cfg *Config) decodeAuthToken() string {
	// It is a B64 encoded string
	// of the API key encrypted with the device's serial.
	// This is obscurity, but it is all we can do on a RPi
	serial := []byte(cfg.GetSerial())
	// ("serial", fmt.Sprintf("%v", serial))
	var key [32]byte
	copy(key[:], serial)
	// log.Println("token", cfg.Auth.Token)
	b64, err := base64.StdEncoding.DecodeString(cfg.Auth.Token)
	if err != nil {
		log.Println("config: cannot b64 decode auth token.")
		log.Println(err.Error())
	}
	dec, err := cryptopasta.Decrypt(b64, &key)
	if err != nil {
		log.Println("config: failed to decrypt auth token after decoding")
		// log.Println("key", fmt.Sprintf("%v", key))
		log.Println(err.Error())
	}

	return string(dec)
}

var states = []string{"AA,AE,AK,AL,AP,AR,AS,AZ,CA,CO,CT,CZ,DE,FL,FM,GA,GU,HI,ID,IL,IN,IA,KS,KY,LA,ME,MD,MA,MH,MI,MN,MS,MO,MT,NE,NV,NH,NJ,NM,NY,NC,ND,OH,OK,OR,PA,PR,PW,RI,SC,SD,TN,TX,UT,VI,VT,VA,WA,WV,WI,WY"}

var patterns = map[string]string{
	"Auth.FCFSId":    fmt.Sprintf(`[%v][0-9]{4}-[0-9]{3}`, states),
	"Auth.DeviceTag": `[a-zA-Z-].+`,
	"LogLevel":       "{DEBUG|debug|INFO|info|WARN|warn|ERROR|error|FATAL|fatal}",
	"RunMode":        "{DEV|dev|PROD|prod|DEVELOP|develop|PRODUCTION|production}",
}

//BROKEN
func getValue(chain string, s interface{}) interface{} {
	reflectType := reflect.TypeOf(s).Elem()
	reflectValue := reflect.ValueOf(s).Elem()
	next := strings.Split(chain, ".")[0]
	rest := strings.Split(chain, ".")[1:]
	for i := 0; i < reflectType.NumField(); i++ {
		typeName := reflectType.Field(i).Name
		// valueType := reflectValue.Field(i).Type()
		valueValue := reflectValue.Field(i).Interface()

		log.Println("typeName", typeName)
		if typeName == next && len(rest) == 0 {
			return valueValue
		} else {
			return getValue(strings.Join(rest, "."), valueValue)
		}
	}
	return nil
}

// THIS IS BROKEN.
func (cfg *Config) Validate() {

	for tag, pat := range patterns {
		value := fmt.Sprintf("%v", getValue(tag, cfg))
		log.Printf("tag [ %v ], pattern %v, value [ %v ]\n", tag, pat, value)
		b, err := regexp.Match(pat, []byte(value))
		if !b || err != nil {
			log.Fatalf("Tag [%v] is invalid; must match pattern '%v'", tag, pat)
		}

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

	cfg.StorageMode = "api"
	cfg.RunMode = "prod"
	cfg.Clock = clock.New()
	cfg.Local.Crontab = "0 0 * * *"
	cfg.Local.SummaryDB = "/opt/imls/summary.sqlite"
	cfg.Local.WebDirectory = "/www/imls"
}

type Config struct {
	Auth struct {
		Token     string `yaml:"api_token"`
		DeviceTag string `yaml:"device_tag"`
		FCFSId    string `yaml:"fcfs_seq_id"`
	} `yaml:"auth"`
	LogLevel   string   `yaml:"log_level"`
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
	LshwPath      string `yaml:"lshw_path"`
	Manufacturers struct {
		Db string `yaml:"db"`
	} `yaml:"manufacturers"`
	SessionId   string      // No YAML equiv.
	Serial      string      `yaml:"serial"`
	StorageMode string      `yaml:"storagemode"`
	RunMode     string      `yaml:"runmode"`
	Clock       clock.Clock // No YAML equiv.
	Local       struct {
		Crontab      string `yaml:"crontab"`
		SummaryDB    string `yaml:"summary_db"`
		WebDirectory string `yaml:"web_directory"`
	} `yaml:"local"`
}
