package state

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
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
	viper.SetConfigName("session-counter")
	viper.SetConfigType("ini")
	viper.SetConfigFile(configPath)
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Info().Msg("no configuration found: writing")
		viper.SafeWriteConfig()
	}
	log.Info().Msg(fmt.Sprintf("using configuration: %s", viper.ConfigFileUsed()))
}

func GetSerial() string {
	return getCachedSerial()
}

func GetFCFSSeqID() string {
	return viper.GetString("device.fcfs_id")
}

func GetDeviceTag() string {
	return viper.GetString("device.tag")
}

// GetAPIKey decodes the api key stored in the ini file.
func GetAPIKey() string {
	return viper.GetString("device.api_key")
}

func SetFCFSSeqID(id string) {
	viper.Set("device.fcfs_id", id)
}

func SetDeviceTag(tag string) {
	viper.Set("device.tag", tag)
}

func SetAPIKey(key string) {
	viper.Set("device.api_key", key)
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
	// viper.GetStringSlice does not work with ini file defaults
	loggers := viper.GetString("log.loggers")
	return strings.Split(loggers, ",")
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

func GetDurationsPath() string {
	return viper.GetString("db.durations")
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

func GetIpPath() string {
	return viper.GetString("ip.path")
}

func GetIwPath() string {
	return viper.GetString("iw.path")
}

func GetLshwPath() string {
	return viper.GetString("lshw.path")
}

func GetWlanHelperPath() string {
	return viper.GetString("wlanhelper.path")
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
	// these must be filled in by the user. NB: these settings will _not_ be
	// present in the config and are set here for explicitness.
	viper.SetDefault("device.api_key", "")
	viper.SetDefault("device.fcfs_id", "")
	viper.SetDefault("device.tag", "")
	// defaults for running in production
	viper.SetDefault("config.minimum_minutes", 5)
	viper.SetDefault("config.maximum_minutes", 600)
	viper.SetDefault("log.level", "DEBUG")
	viper.SetDefault("log.loggers", "local:stderr,local:tmp,api:directus")
	viper.SetDefault("mode.storage", "api")
	viper.SetDefault("mode.run", "prod")
	viper.SetDefault("api.scheme", "https")
	viper.SetDefault("api.host", "rabbit-phase-4.app.cloud.gov")
	viper.SetDefault("api.uri", "/items/durations_v2/")
	viper.SetDefault("cron.reset", "0 0 * * *")
	viper.SetDefault("wireshark.duration", 45)
	if runtime.GOOS == "windows" {
		viper.SetDefault("wireshark.path", "c:/Program Files/Wireshark/tshark.exe")
		viper.SetDefault("wlanhelper.path", "c:/Windows/System32/Npcap/WlanHelper.exe")
		viper.SetDefault("www.root", "c:/imls")
		viper.SetDefault("www.images", "c:/imls/images")
		viper.SetDefault("db.durations", "c:/imls/durations.sqlite")
		viper.SetDefault("db.queues", "c:/imls/queues.sqlite")
	} else {
		viper.SetDefault("iw.path", "/usr/sbin/iw")
		viper.SetDefault("ip.path", "/usr/sbin/ip")
		viper.SetDefault("wireshark.path", "/usr/bin/tshark")
		viper.SetDefault("lshw.path", "/usr/bin/lshw")
		viper.SetDefault("www.root", "/www/imls")
		viper.SetDefault("www.images", "/www/imls/images")
		viper.SetDefault("db.durations", "/www/imls/durations.sqlite")
		viper.SetDefault("db.queues", "/www/imls/queues.sqlite")
	}
}
