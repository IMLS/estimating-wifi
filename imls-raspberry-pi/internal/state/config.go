package state

import (
	"encoding/base64"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/benbjohnson/clock"
	yaml "gopkg.in/yaml.v2"
	"gsa.gov/18f/internal/cryptopasta"
	"gsa.gov/18f/internal/interfaces"
	"gsa.gov/18f/internal/logwrapper"
	"gsa.gov/18f/internal/wifi-hardware-search/config"
)

var theConfig *CFG
var once sync.Once

// The config is a singleton. We can get a new,
// empty config once.
func NewConfig() *CFG {
	once.Do(func() {
		UnsafeNewConfig()
	})
	return theConfig
}

func UnsafeNewConfig() *CFG {
	theConfig = &CFG{}
	setDefaults()
	return theConfig
}

func (cfg *CFG) InitConfig() {
	theConfig.Logging.Log = logwrapper.NewLogger(theConfig)
	theConfig.Databases.DurationsDB = NewSqliteDB(theConfig.Databases.DurationsPath)
	theConfig.Databases.QueuesDB = NewSqliteDB(theConfig.Databases.QueuesPath)
	theConfig.Databases.ManufacturersDB = NewSqliteDB(theConfig.Databases.ManufacturersPath)
	theConfig.InitializeSessionId()
}

// Or, we can get a new config from a path.
// We cannot get both.
func NewConfigFromPath(path string) {
	once.Do(func() {
		UnsafeNewConfigFromPath(path)
	})
}

func UnsafeNewConfigFromPath(path string) {
	theConfig = &CFG{}
	setDefaults()
	readConfig(path)
	if theConfig.Clock == nil {
		log.Println("clock should not be nil")
		log.Fatal()
	}
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
		theConfig = newcfg

		// The API key will need to be decoded into memory.
		if len(theConfig.Device.Token) > 0 {
			decodeAuthToken()
		}

		// FIXME: Because this is an external lib, we need to set the
		// path there. Or, pass the config?
		if len(theConfig.Executables.LshwPath) > 0 {
			config.SetLSHWLocation(theConfig.Executables.LshwPath)
		}

		// Need to reset the clock pointer...
		// Gets wiped out by the read.
		theConfig.Clock = clock.New()
		theConfig.InitializeSessionId()
	} else {
		log.Fatalf("could not find config: %v\n", path)
	}
}

// We can return the existing singleton as many
// times as we like. However, if it is not initialized,
// the program should exit.
func GetConfig() *CFG {
	if theConfig == nil {
		log.Fatal("cannot retrieve nil config")
	}
	return theConfig
}

// interface Config

// See serial.go for GetSerial()

func (cfg *CFG) GetFCFSSeqId() string {
	return theConfig.Device.FCFSId
}

func (cfg *CFG) GetDeviceTag() string {
	return theConfig.Device.DeviceTag
}

func (cfg *CFG) GetAPIKey() string {
	return theConfig.Device.Token
}

func (cfg *CFG) GetLogLevel() string {
	valid := []string{"DEBUG", "INFO", "WARN", "ERROR"}
	ok := false
	for _, v := range valid {
		if theConfig.Logging.LogLevel == v {
			ok = true
		}
		if theConfig.Logging.LogLevel == strings.ToLower(v) {
			ok = true
		}
	}
	if !ok {
		log.Printf("invalid log level in config: %v\n",
			theConfig.Logging.LogLevel)
		return "ERROR"
	} else {
		return theConfig.Logging.LogLevel
	}
}

func (cfg *CFG) GetLoggers() []string {
	return theConfig.Logging.Loggers
}

func (cfg *CFG) Log() interfaces.Logger {
	return theConfig.Logging.Log
}

func (cfg *CFG) GetEventsUri() string {
	return theConfig.Umbrella.Paths.Events
}

func (cfg *CFG) GetDurationsUri() string {
	return theConfig.Umbrella.Paths.Durations
}

// See sessionid.go for related interface implementation

func (cfg *CFG) IsStoringToApi() bool {
	return strings.Contains(strings.ToLower(theConfig.StorageMode), "api")
}

func (cfg *CFG) IsStoringLocally() bool {
	either := false
	for _, s := range []string{"local", "sqlite"} {
		either = either || strings.Contains(strings.ToLower(theConfig.StorageMode), s)
	}
	return either
}

func (cfg *CFG) IsProductionMode() bool {
	return strings.Contains(strings.ToLower(theConfig.RunMode), "prod")
}

func (cfg *CFG) IsDeveloperMode() bool {
	either := false
	for _, s := range []string{"dev", "test"} {
		either = either || strings.Contains(strings.ToLower(theConfig.RunMode), s)
	}
	if either {
		log.Println("running in developer mode")
	}
	return either
}

func (cfg *CFG) IsTestMode() bool {
	return strings.Contains(strings.ToLower(theConfig.RunMode), "test")
}

func (cfg *CFG) GetDurationsDatabase() interfaces.Database {
	return cfg.Databases.DurationsDB
}

func (cfg *CFG) GetQueuesDatabase() interfaces.Database {
	return cfg.Databases.QueuesDB
}
func (cfg *CFG) GetManufacturerDatabase() interfaces.Database {
	return cfg.Databases.ManufacturersDB
}

func (cfg *CFG) GetClock() clock.Clock {
	return cfg.Clock
}

func (cfg *CFG) GetMinimumMinutes() int {
	return cfg.Monitoring.MinimumMinutes
}

func (cfg *CFG) GetMaximumMinutes() int {
	return cfg.Monitoring.MaximumMinutes
}

/////////

func decodeAuthToken() {
	decodeSerial()
	// theConfig.Device.Token = decodeAuthToken()
	// It is a B64 encoded string
	// of the API key encrypted with the device's serial.
	// This is obscurity, but it is all we can do on a RPi
	serial := []byte(theConfig.GetSerial())
	// ("serial", fmt.Sprintf("%v", serial))
	var key [32]byte
	copy(key[:], serial)
	// log.Println("token", cfg.Auth.Token)
	b64, err := base64.StdEncoding.DecodeString(theConfig.Device.Token)
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
	theConfig.Device.Token = string(dec)
}

// var states = []string{"AA,AE,AK,AL,AP,AR,AS,AZ,CA,CO,CT,CZ,DE,FL,FM,GA,GU,HI,ID,IL,IN,IA,KS,KY,LA,ME,MD,MA,MH,MI,MN,MS,MO,MT,NE,NV,NH,NJ,NM,NY,NC,ND,OH,OK,OR,PA,PR,PW,RI,SC,SD,TN,TX,UT,VI,VT,VA,WA,WV,WI,WY"}

func setDefaults() {
	theConfig.Logging.LogLevel = "DEBUG"
	theConfig.Logging.Loggers = []string{"local:stderr", "local:tmp", "api:directus"}

	theConfig.Monitoring.UniquenessWindow = 24 * 60
	theConfig.Monitoring.MinimumMinutes = 5
	theConfig.Monitoring.MaximumMinutes = 600

	theConfig.Umbrella.Scheme = "https"
	theConfig.Umbrella.Host = "api.data.gov"
	theConfig.Umbrella.Paths.Durations = "/TEST/10x-imls/v2/durations/"
	theConfig.Umbrella.Paths.Events = "/TEST/10x-imls/v2/events/"

	theConfig.Executables.Wireshark.Duration = 45
	theConfig.Executables.Wireshark.Path = "/usr/bin/tshark"

	theConfig.Databases.ManufacturersPath = "/opt/imls/manufacturers.sqlite"
	theConfig.Databases.DurationsPath = "/www/imls/durations.sqlite"
	theConfig.Databases.QueuesPath = "/www/imls/queues.sqlite"

	theConfig.StorageMode = "api"
	theConfig.RunMode = "prod"
	theConfig.Clock = clock.New()
	theConfig.Monitoring.ResetCron = "0 0 * * *"

	theConfig.Paths.WWW.Root = "/www/imls"
	theConfig.Paths.WWW.Images = "/www/imls/images"

}

type CFG struct {
	Device struct {
		Token     string `yaml:"api_token"`
		DeviceTag string `yaml:"device_tag"`
		FCFSId    string `yaml:"fcfs_seq_id"`
	} `yaml:"device"`
	Logging struct {
		LogLevel string            `yaml:"log_level"`
		Loggers  []string          `yaml:"loggers"`
		Log      interfaces.Logger `yaml:"-"`
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
		ManufacturersPath string              `yaml:"manufacturers_path"`
		ManufacturersDB   interfaces.Database `yaml:"-"`
		DurationsPath     string              `yaml:"durations_path"`
		DurationsDB       interfaces.Database `yaml:"-"`
		QueuesPath        string              `yaml:"queues_path"`
		QueuesDB          interfaces.Database `yaml:"-"`
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
