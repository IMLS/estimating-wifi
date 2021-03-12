package model

type Config struct {
	Server struct {
		Scheme     string `yaml:"scheme"`
		Host       string `yaml:"host"`
		Collection string `yaml:"collection"`
	} `yaml:"server"`
	Wireshark struct {
		Duration  int    `yaml:"duration"`
		Rounds    int    `yaml:"rounds"`
		Threshold int    `yaml:"threshold"`
		Adapter   string `yaml:"adapter"`
		Path      string `yaml:"path"`
	} `yaml:"wireshark"`
	Manufacturers struct {
		Db string `yaml:"db"`
	} `yaml:"manufacturers`
}
