package model

type Config struct {
	Server struct {
		Scheme     string `yaml:"scheme"`
		Host       string `yaml:"host"`
		Collection string `yaml:"collection"`
	} `yaml:"server"`
	Wireshark struct {
		Duration int    `yaml:"duration"`
		Adapter  string `yaml:"adapter"`
	} `yaml:"wireshark"`
	Manufacturers struct {
		Db string `yaml:"db"`
	} `yaml:"manufacturers`
}
