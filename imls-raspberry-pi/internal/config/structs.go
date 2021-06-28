package config

// At /opt/imls/config.yaml
type Config struct {
	Auth struct {
		Token     string `yaml:"api_token"`
		DeviceTag string `yaml:"device_tag"`
		FCFSId    string `yaml:"fcfs_seq_id"`
	} `yaml:"auth"`
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
		Logfile      string `yaml:"logfile"`
		Crontab      string `yaml:"crontab"`
		SummaryDB    string `yaml:"summary_db"`
		TemporaryDB  string `yaml:"temporary_db"`
		WebDirectory string `yaml:"web_directory"`
	} `yaml:"local"`
}

func (cfg *Config) SetDefaults() {
	cfg.Monitoring.PingInterval = 30
	cfg.Monitoring.MaxHTTPErrorCount = 8
	cfg.Monitoring.HTTPErrorIntervalMins = 10
	cfg.Monitoring.UniquenessWindow = 120
	cfg.Monitoring.MinimumMinutes = 30
	cfg.Monitoring.MaximumMinutes = 600

	cfg.Umbrella.Scheme = "https"
	cfg.Umbrella.Host = "api.data.gov"
	cfg.Umbrella.Data = "/TEST/10x-imls/v1/wifi/"
	cfg.Umbrella.Logging = "/TEST/10x-imls/v1/events/"

	cfg.Wireshark.Duration = 45
	cfg.Wireshark.Path = "/usr/bin/tshark"
	cfg.Wireshark.CheckWlan = "1"

	cfg.Manufacturers.Db = "/opt/imls/manufacturers.sqlite"

	cfg.StorageMode = "sqlite"

	cfg.Local.Logfile = "/opt/imls/log.json"
	cfg.Local.Crontab = "0 */6 * * *"
	cfg.Local.SummaryDB = "/opt/imls/summary.sqlite"
	cfg.Local.TemporaryDB = "/tmp/imls.sqlite"
	cfg.Local.WebDirectory = "/www/imls"
}
