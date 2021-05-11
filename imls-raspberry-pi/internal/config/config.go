package config

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

var Verbose bool = false
var configPath = "/opt/imls/config.yaml"

func SetConfigPath(path string) {
	configPath = path
}

func GetConfigPath() string {
	return configPath
}

func parseConfigFile(path string) (*Config, error) {
	_, err := os.Stat(path)

	// Stat will set an error if the file cannot be found.
	if err == nil {
		f, err := os.Open(path)
		if err != nil {
			log.Fatal("parseConfigFile: could not open configuration file. Exiting.")
		}
		defer f.Close()
		var cfg *Config
		decoder := yaml.NewDecoder(f)
		err = decoder.Decode(&cfg)
		if err != nil {
			log.Fatalf("parseConfigFile: could not decode YAML:\n%v\n", err)
		}

		return cfg, nil
	} else {
		log.Printf("parseConfigFile: could not find config: %v\n", path)
	}
	return nil, fmt.Errorf("config: could not find config file [%v]", path)
}

func ReadConfig() *Config {
	// We expect config to be here:
	//   * /opt/imls/config.yaml

	cfg, err := parseConfigFile(GetConfigPath())
	if err != nil {
		log.Fatal("config: please provide alternative configuration path.")
	}

	return cfg
}
