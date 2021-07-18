package state

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/benbjohnson/clock"
	"github.com/jmoiron/sqlx"
	yaml "gopkg.in/yaml.v2"
	"gsa.gov/18f/internal/cryptopasta"
	"gsa.gov/18f/internal/wifi-hardware-search/config"
)

var the_config *CFG
var once sync.Once

// The config is a singleton. We can get a new,
// empty config once.
func NewConfig() *CFG {
	once.Do(func() {
		the_config = &CFG{}
		setDefaults()
	})
	return the_config
}

// Or, we can get a new config from a path.
// We cannot get both.
func NewConfigFromPath(path string) {
	once.Do(func() {
		the_config = &CFG{}
		setDefaults()
		readConfig(path)
		if the_config.Clock == nil {
			log.Println("clock should not be nil")
			log.Fatal()
		}
		the_config.InitializeSessionId()
	})
}

func readConfig(path string) {
	_, err := os.Stat(path)
	// Stat will set an error if the file cannot be found.
	if err == nil {
		f, err := os.Open(path)
		if err != nil {
			log.Fatal("config: could not open configuration file. Exiting.")
		}
		defer f.Close()
		var newcfg *CFG
		decoder := yaml.NewDecoder(f)
		err = decoder.Decode(&newcfg)
		if err != nil {
			log.Fatalf("config: could not decode YAML:\n%v\n", err)
		}
		the_config = newcfg

		// The API key will need to be decoded into memory.
		if len(the_config.Device.Token) > 0 {
			decodeAuthToken()
		}

		// FIXME: Because this is an external lib, we need to set the
		// path there. Or, pass the config?
		if len(the_config.Executables.LshwPath) > 0 {
			config.SetLSHWLocation(the_config.Executables.LshwPath)
		}

		// Need to reset the clock pointer...
		// Gets wiped out by the read.
		the_config.Clock = clock.New()

	} else {
		log.Fatalf("could not find config: %v\n", path)
	}
}

// We can return the existing singleton as many
// times as we like. However, if it is not initialized,
// the program should exit.
func GetConfig() *CFG {
	if the_config == nil {
		log.Fatal("cannot retrieve nil config")
	}
	return the_config
}

// interface Config

// See serial.go for GetSerial()

func (cfg *CFG) GetFCFSSeqId() string {
	return the_config.Device.FCFSId
}

func (cfg *CFG) GetDeviceTag() string {
	return the_config.Device.DeviceTag
}

func (cfg *CFG) GetAPIKey() string {
	return the_config.Device.Token
}

func (cfg *CFG) GetLogLevel() string {
	valid := []string{"DEBUG", "INFO", "WARN", "ERROR"}
	ok := false
	for _, v := range valid {
		if the_config.Logging.LogLevel == v {
			ok = true
		}
		if the_config.Logging.LogLevel == strings.ToLower(v) {
			ok = true
		}
	}
	if !ok {
		log.Printf("invalid log level in config: %v\n",
			the_config.Logging.LogLevel)
		return "ERROR"
	} else {
		return the_config.Logging.LogLevel
	}
}

func (cfg *CFG) GetLoggers() []string {
	return the_config.Logging.Loggers
}

func (cfg *CFG) GetEventsUri() string {
	return the_config.Umbrella.Paths.Events
}

func (cfg *CFG) GetDurationsUri() string {
	return the_config.Umbrella.Paths.Durations
}

// See sessionid.go for related interface implementation

func IsStoringToApi() bool {
	return strings.Contains(strings.ToLower(the_config.StorageMode), "api")
}

func IsStoringLocally() bool {
	either := false
	for _, s := range []string{"local", "sqlite"} {
		either = either || strings.Contains(strings.ToLower(the_config.StorageMode), s)
	}
	return either
}

func IsProductionMode() bool {
	return strings.Contains(strings.ToLower(the_config.RunMode), "prod")
}

func IsDeveloperMode() bool {
	either := false
	for _, s := range []string{"dev", "test"} {
		either = either || strings.Contains(strings.ToLower(the_config.RunMode), s)
	}
	if either {
		log.Println("running in developer mode")
	}
	return either
}

func IsTestMode() bool {
	return strings.Contains(strings.ToLower(the_config.RunMode), "test")
}

func decodeAuthToken() string {
	decodeSerial()
	the_config.Device.Token = decodeAuthToken()
	// It is a B64 encoded string
	// of the API key encrypted with the device's serial.
	// This is obscurity, but it is all we can do on a RPi
	serial := []byte(GetSerial())
	// ("serial", fmt.Sprintf("%v", serial))
	var key [32]byte
	copy(key[:], serial)
	// log.Println("token", cfg.Auth.Token)
	b64, err := base64.StdEncoding.DecodeString(the_config.Device.Token)
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

func setDefaults() {
	the_config.Logging.LogLevel = "DEBUG"
	the_config.Logging.Loggers = []string{"local:stderr", "local:tmp", "api:directus"}

	the_config.Monitoring.UniquenessWindow = 24 * 60
	the_config.Monitoring.MinimumMinutes = 5
	the_config.Monitoring.MaximumMinutes = 600

	the_config.Umbrella.Scheme = "https"
	the_config.Umbrella.Host = "api.data.gov"
	the_config.Umbrella.Paths.Durations = "/TEST/10x-imls/v2/durations/"
	the_config.Umbrella.Paths.Events = "/TEST/10x-imls/v2/events/"

	the_config.Executables.Wireshark.Duration = 45
	the_config.Executables.Wireshark.Path = "/usr/bin/tshark"

	the_config.Databases.Manufacturers = "/opt/imls/manufacturers.sqlite"

	the_config.StorageMode = "api"
	the_config.RunMode = "prod"
	the_config.Clock = clock.New()
	the_config.Monitoring.ResetCron = "0 0 * * *"

	the_config.Paths.WWW.Root = "/www/imls"
	the_config.Paths.WWW.Images = "/www/imls/images"
}

type CFG struct {
	Device struct {
		Token     string `yaml:"api_token"`
		DeviceTag string `yaml:"device_tag"`
		FCFSId    string `yaml:"fcfs_seq_id"`
	} `yaml:"device"`
	Logging struct {
		LogLevel string   `yaml:"log_level"`
		Loggers  []string `yaml:"loggers"`
	} `yaml:"logging"`
	Monitoring struct {
		UniquenessWindow int    `yaml:"uniqueness_window"`
		MinimumMinutes   int    `yaml:"minimum_minutes"`
		MaximumMinutes   int    `yaml:"maximum_minutes"`
		ResetCron        string `yaml:"reset_cron"`
	} `yaml:"monitoring"`
	Umbrella struct {
		Scheme string `yaml:"scheme"`
		Host   string `yaml:"host"`
		Paths  struct {
			Durations string `yaml:"durations"`
			Events    string `yaml:"events"`
		} `yaml:"paths"`
	} `yaml:"umbrella"`
	Executables struct {
		Wireshark struct {
			Duration int    `yaml:"duration"`
			Path     string `yaml:"wireshark_path"`
		} `yaml:"wireshark"`
		LshwPath string `yaml:"lshw_path"`
		IpPath   string `yaml:"ip_path"`
	} `yaml:"executables"`
	Databases struct {
		Manufacturers    string   `yaml:"manufacturers"`
		ManufacturersPtr *sqlx.DB `yaml:"-"`
		Durations        string   `yaml:"durations"`
		DurationsPtr     *sqlx.DB `yaml:"-"`
		Queues           string   `yaml:"queues"`
		QueuesPtr        *sqlx.DB `yaml:"-"`
	} `yaml:"databases"`
	Paths struct {
		WWW struct {
			Root   string `yaml:"root"`
			Images string `yaml:"images"`
		} `yaml:"www"`
	} `yaml:"paths"`
	Serial      string      `yaml:"serial"`
	StorageMode string      `yaml:"storagemode"`
	RunMode     string      `yaml:"runmode"`
	SessionId   int         `yaml:"-"` // ignore
	Clock       clock.Clock `yaml:"-"` // ignore

}
