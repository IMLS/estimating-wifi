package config

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v2"
	"gsa.gov/18f/cryptopasta"
)

var Verbose bool = false

// This might not exist.
// A wrapper script should require a valid path.
var configPath = "config.yaml"

func SetConfigPath(path string) {
	configPath = path
}

func GetConfigPath() string {
	return configPath
}

func decodeAuthToken(token string) string {
	// It is a B64 encoded string
	// of the API key encrypted with the device's serial.
	// This is obscurity, but it is all we can do on a RPi
	serial := []byte(GetSerial())
	var key [32]byte
	copy(key[:], serial)
	b64, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		if Verbose {
			log.Println("config: cannot b64 decode auth token.")
		}
	}
	dec, err := cryptopasta.Decrypt(b64, &key)
	if err != nil {
		if Verbose {
			log.Println("config: failed to decrypt auth token after decoding")
		}
	}

	return string(dec)
}

func ReadConfig(path string) (*Config, error) {
	_, err := os.Stat(path)

	// Stat will set an error if the file cannot be found.
	if err == nil {
		f, err := os.Open(path)
		if err != nil {
			log.Fatal("config: could not open configuration file. Exiting.")
		}
		defer f.Close()
		var cfg *Config
		decoder := yaml.NewDecoder(f)
		err = decoder.Decode(&cfg)
		if err != nil {
			log.Fatalf("config: could not decode YAML:\n%v\n", err)
		}

		// The API key will need to be decoded into memory.
		cfg.Auth.Token = decodeAuthToken(cfg.Auth.Token)

		return cfg, nil
	} else {
		log.Printf("config: could not find config: %v\n", path)
	}
	return nil, fmt.Errorf("config: could not find config file [%v]", path)
}
