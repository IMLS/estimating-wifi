package config

import (
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v2"
	"gsa.gov/18f/session-counter/constants"
	"gsa.gov/18f/session-counter/cryptopasta"
)

var Verbose bool = false

func parseConfigFile(filepath string) (*Config, error) {
	_, err := os.Stat(filepath)

	// Stat will set an error if the file cannot be found.
	if err == nil {
		f, err := os.Open(filepath)
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
		log.Printf("parseConfigFile: could not find config: %v\n", filepath)
	}
	return nil, fmt.Errorf("config: could not find config file [%v]", filepath)
}

func devConfig() *Config {
	// FIXME consider turning this into an env var
	cfgPtr := flag.String("config", "config.yaml", "config file")
	flag.Parse()
	cfg, err := parseConfigFile(*cfgPtr)
	if err != nil {
		log.Println("config: could not load dev config. Exiting.")
		log.Fatalln(err)
	}
	return cfg
}


func ReadAuth() (*AuthConfig, error) {
	_, err := os.Stat(constants.AuthPath)
	if err != nil {
		return &AuthConfig{}, fmt.Errorf("readToken: cannot find default token file at [%v]", constants.AuthPath)
	}

	f, err := os.Open(constants.AuthPath)
	if err != nil {
		log.Fatal("readToken: could not open token file. Exiting.")
	}
	defer f.Close()
	var auth *AuthConfig
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&auth)
	if err != nil {
		log.Fatalf("readToken: could not decode YAML:\n%v\n", err)
	}

	// Unencrypt the API token.
	// It is a B64 encoded string
	// of the API key encrypted with the device's serial.
	// This is obscurity, but it is all we can do on a RPi
	serial := []byte(GetSerial())
	var key [32]byte
	copy(key[:], serial)
	b64, err := base64.StdEncoding.DecodeString(auth.Token)
	if err != nil {
		log.Fatal("readToken: cannot b64 decode auth token.")
	}
	dec, err := cryptopasta.Decrypt(b64, &key)
	if err != nil {
		log.Fatal("readToken: failed to decrypt auth token after decoding")
	}
	auth.Token = string(dec)
	return auth, nil
}

func ReadConfig() *Config {
	// We expect config to be here:
	//   * /opt/imls/config.yaml

	cfg, err := parseConfigFile(constants.ConfigPath)
	if err != nil {
		fmt.Printf("config: could not find config at default path [%v]\n", constants.ConfigPath)
		fmt.Println("config: loading dev config")
		return devConfig()
	}

	return cfg
}
