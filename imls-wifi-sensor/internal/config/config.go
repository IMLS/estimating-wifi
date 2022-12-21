package config

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
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

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	err := viper.ReadInConfig()
	if err != nil {
		log.Info().Msg("no configuration found: writing")
		viper.SafeWriteConfig()
	}
	log.Info().Msg(fmt.Sprintf("using configuration: %s", viper.ConfigFileUsed()))
	// configure logging.
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	logLevel := GetLogLevel()
	switch lvl := logLevel; lvl {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}

func GetFSCSID() string {
	return viper.GetString("device.fscs_id")
}

func GetDeviceTag() string {
	return viper.GetString("device.tag")
}

func GetSentryDSN() string {
	return viper.GetString("sentry.dsn")
}

// GetAPIKey decodes the api key stored in the ini file.
func GetAPIKey() string {
	return viper.GetString("device.api_key")
}

func SetFSCSID(id string) {
	viper.Set("device.fscs_id", id)
}

func SetDeviceTag(tag string) {
	viper.Set("device.tag", tag)
}

func SetAPIKey(key string) {
	viper.Set("device.api_key", key)
}

func SetRunMode(mode string) {
	viper.Set("mode.run", mode)
}

func GetLogLevel() string {
	return strings.ToLower(viper.GetString("log.level"))
}

func GetLoggers() []string {
	// viper.GetStringSlice does not work with ini file defaults
	loggers := viper.GetString("log.loggers")
	return strings.Split(loggers, ",")
}

func createURI(what string) string {
	scheme := viper.GetString("api.scheme")
	host := viper.GetString("api.host")
	//port := viper.GetInt("api.port")
	return (scheme + "://" +
		// strings.TrimSuffix(strings.TrimPrefix(host, "/"), "/") +
		// ":" + fmt.Sprint(port) + "/" +
		// strings.TrimPrefix(what, "/"))
		strings.TrimSuffix(strings.TrimPrefix(host, "/"), "/") +
		"/" +
		strings.TrimPrefix(what, "/"))

}

func GetDurationsURI() string {
	path := viper.GetString("api.presences_uri")
	return createURI(path)
}

func GetHeartbeatURI() string {
	path := viper.GetString("api.heartbeat_uri")
	return createURI(path)
}

func GetLoginURI() string {
	path := viper.GetString("api.login_uri")
	return createURI(path)
}

func GetLoginPort() int {
	return viper.GetInt("api.port")
}

func IsProductionMode() bool {
	mode := viper.GetString("mode.run")
	return strings.Contains(strings.ToLower(mode), "prod")
}

func IsDeveloperMode() bool {
	return !IsProductionMode()
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

func GetDeviceMemory() int {
	return viper.GetInt("config.device_memory")
}

func GetResetCron() string {
	return viper.GetString("cron.reset")
}

func GetHeartbeatCron() string {
	return viper.GetString("cron.heartbeat")
}

func GetFakesharkNumPerMinute() int {
	return viper.GetInt("wireshark.fakeshark_num_found_per_minute")
}

func GetFakesharkNumMacs() int {
	return viper.GetInt("wireshark.fakeshark_num_macs")
}

func GetFakesharkMinFound() int {
	return viper.GetInt("wireshark.fakeshark_min_found")
}

func GetFakesharkMaxFound() int {
	return viper.GetInt("wireshark.fakeshark_max_found")
}

func GetDataCollectionCron() string {
	return viper.GetString("cron.data_collection_cron")
}

func SetConfigDefaults() {
	// these must be filled in by the user. NB: these settings will _not_ be
	// present in the config and are set here for explicitness.
	viper.SetDefault("device.api_key", "")
	viper.SetDefault("device.fscs_id", "")
	viper.SetDefault("device.tag", "")
	// defaults for running in production
	viper.SetDefault("config.minimum_minutes", 5)
	viper.SetDefault("config.maximum_minutes", 600)
	viper.SetDefault("config.device_memory", 7200)
	viper.SetDefault("log.level", "DEBUG")
	viper.SetDefault("log.loggers", "local:stderr,local:tmp")
	viper.SetDefault("mode.run", "prod")
	viper.SetDefault("api.scheme", "https")
	viper.SetDefault("api.host", "rabbit-phase-4.app.cloud.gov")
	viper.SetDefault("api.port", 3000)
	viper.SetDefault("api.login_uri", "/rpc/login")
	viper.SetDefault("api.heartbeat_uri", "/rpc/beat_the_heart")
	viper.SetDefault("api.presences_uri", "/rpc/update_presences")
	// At midnight every night
	viper.SetDefault("cron.reset", "0 0 * * *")
	// Every hour
	viper.SetDefault("cron.heartbeat", "0 0/1 * * *")
	viper.SetDefault("cron.data_collection_cron", "*/1 * * * *")
	viper.SetDefault("wireshark.duration", 45)
	viper.SetDefault("wireshark.fakeshark_num_macs", 500)
	viper.SetDefault("wireshark.fakeshark_min_found", 20)
	viper.SetDefault("wireshark.fakeshark_max_found", 100)
	if runtime.GOOS == "windows" {
		viper.SetDefault("wireshark.path", "c:/Program Files/Wireshark/tshark.exe")
		viper.SetDefault("wlanhelper.path", "c:/Windows/System32/Npcap/WlanHelper.exe")
	} else {
		viper.SetDefault("iw.path", "/usr/sbin/iw")
		viper.SetDefault("ip.path", "/usr/sbin/ip")
		viper.SetDefault("wireshark.path", "/usr/bin/tshark")
		viper.SetDefault("lshw.path", "/usr/bin/lshw")
	}
}
