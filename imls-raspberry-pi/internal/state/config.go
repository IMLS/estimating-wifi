package state

import (
	"encoding/base64"
	"runtime"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"gsa.gov/18f/internal/cryptopasta"
	"gsa.gov/18f/internal/interfaces"
	"gsa.gov/18f/internal/structs"
)

func SetConfigAtPath(configPath string) {
	SetConfigDefaults()
	viper.AddConfigPath(".")
	if runtime.GOOS == "linux" {
		viper.AddConfigPath("/etc/imls/")
	}
	if runtime.GOOS == "windows" {
		viper.AddConfigPath("%PROGRAMDATA%\\imls")
	}
	viper.SetConfigType("ini")
	viper.SetConfigFile(configPath)
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err == nil {
		log.Info().Msg(viper.ConfigFileUsed())
	} else {
		log.Fatal().
			Err(err).
			Msg("could not find configuration file")
	}
}

func GetSerial() string {
	// allow the serial to be stored so we can test out different serial and api
	// key settings. if not, default to reading from /proc (as a cached read)
	serial := viper.GetString("config.device_serial")
	if serial == "" {
		return getCachedSerial()
	}
	return serial
}

func GetFCFSSeqID() string {
	return viper.GetString("config.fcfs_id")
}

func GetDeviceTag() string {
	return viper.GetString("config.device_tag")
}

// GetAPIKey decodes the api key stored in the ini file.
func GetAPIKey() string {
	apiKey := viper.GetString("config.api_key")
	serial := GetSerial()
	var key [32]byte
	copy(key[:], serial)
	b64, err := base64.StdEncoding.DecodeString(apiKey)
	if err != nil {
		log.Error().
			Err(err).
			Msg("cannot b64 decode")
	}
	dec, err := cryptopasta.Decrypt(b64, &key)
	if err != nil {
		log.Error().
			Err(err).
			Msg("failed to decrypt after decoding")
	}
	return string(dec)
}

func SetFCFSSeqID(id string) {
	viper.Set("config.fcfs_id", id)
}

func SetDeviceTag(tag string) {
	viper.Set("config.device_tag", tag)
}

func SetAPIKey(key string) {
	viper.Set("config.api_key", key)
}

func SetStorageMode(mode string) {
	viper.Set("mode.storage", mode)
}

func SetRunMode(mode string) {
	viper.Set("mode.run", mode)
}

func SetQueuesPath(where string) {
	viper.Set("db.queues", where)
}

func SetDurationsPath(where string) {
	viper.Set("db.durations", where)
}

func SetRootPath(mode string) {
	viper.Set("www.root", mode)
}

func SetImagesPath(mode string) {
	viper.Set("www.images", mode)
}

func SetUniquenessWindow(window int) {
	viper.Set("config.uniqueness_window", window)
}

func GetLogLevel() string {
	return viper.GetString("log.level")
}

func GetLoggers() []string {
	return viper.GetStringSlice("log.loggers")
}

func GetDurationsURI() string {
	scheme := viper.GetString("api.scheme")
	host := viper.GetString("api.host")
	path := viper.GetString("api.uri")
	return (scheme + "://" +
		removeLeadingAndTrailingSlashes(host) +
		startsWithSlash(removeLeadingSlashes(path)))
}

func IsStoringToAPI() bool {
	mode := viper.GetString("mode.storage")
	return strings.Contains(strings.ToLower(mode), "api")
}

func IsStoringLocally() bool {
	mode := viper.GetString("mode.storage")
	either := false
	for _, s := range []string{"local", "sqlite"} {
		either = either || strings.Contains(strings.ToLower(mode), s)
	}
	return either
}

func IsProductionMode() bool {
	mode := viper.GetString("mode.run")
	return strings.Contains(strings.ToLower(mode), "prod")
}

func IsDeveloperMode() bool {
	mode := viper.GetString("mode.run")
	either := false
	for _, s := range []string{"dev", "test"} {
		either = either || strings.Contains(strings.ToLower(mode), s)
	}
	if either {
		log.Info().Msg("running in developer mode")
	}
	return either
}

func IsTestMode() bool {
	mode := viper.GetString("mode.run")
	return strings.Contains(strings.ToLower(mode), "test")
}

func GetDurationsDatabase() interfaces.Database {
	path := viper.GetString("db.durations")
	// always make sure we have a durations db created
	db := NewSqliteDB(path)
	if !db.CheckTableExists("durations") {
		db.CreateTableFromStruct(structs.Duration{})
	}
	return db
}

func GetQueuesDatabase() interfaces.Database {
	path := viper.GetString("db.queues")
	return NewSqliteDB(path)
}

func GetWiresharkPath() string {
	return viper.GetString("wireshark.path")
}

func GetWiresharkDuration() int {
	return viper.GetInt("wireshark.duration")
}

func GetMinimumMinutes() int {
	return viper.GetInt("config.minimum_minutes")
}

func GetMaximumMinutes() int {
	return viper.GetInt("config.maximum_minutes")
}

func GetUniquenessWindow() int {
	return viper.GetInt("config.uniqueness_window")
}

func GetResetCron() string {
	return viper.GetString("cron.reset")
}

func GetWWWRoot() string {
	return viper.GetString("www.root")
}

func GetWWWImages() string {
	return viper.GetString("www.images")
}

func SetConfigDefaults() {
	// config.api_key filled in by user
	// config.fcfs_id filled in by user
	// config.device_tag filled in by user
	// config.device_serial filled in by device or user
	viper.SetDefault("config.minimum_minutes", 5)
	viper.SetDefault("config.maximum_minutes", 600)
	viper.SetDefault("log.level", "DEBUG")
	viper.SetDefault("log.loggers", []string{"local:stderr", "local:tmp", "api:directus"})
	viper.SetDefault("mode.storage", "api")
	viper.SetDefault("mode.run", "prod")
	viper.SetDefault("api.scheme", "https")
	viper.SetDefault("api.host", "rabbit-phase-4.app.cloud.gov")
	viper.SetDefault("api.uri", "/items/durations_v2/")
	viper.SetDefault("cron.reset", "0 0 * * *")
	viper.SetDefault("wireshark.duration", 45)
	if runtime.GOOS == "windows" {
		viper.SetDefault("wireshark.path", "c:/Program Files/Wireshark/tshark.exe")
		viper.SetDefault("www.root", "c:/imls")
		viper.SetDefault("www.images", "c:/imls/images")
		viper.SetDefault("db.durations", "c:/imls/durations.sqlite")
		viper.SetDefault("db.queues", "c:/imls/queues.sqlite")
	} else {
		viper.SetDefault("wireshark.path", "/usr/bin/tshark")
		viper.SetDefault("www.root", "/www/imls")
		viper.SetDefault("www.images", "/www/imls/images")
		viper.SetDefault("db.durations", "/www/imls/durations.sqlite")
		viper.SetDefault("db.queues", "/www/imls/queues.sqlite")
	}
}
