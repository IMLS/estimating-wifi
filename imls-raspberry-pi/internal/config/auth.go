package config

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v2"
	"gsa.gov/18f/cryptopasta"
)

var authPath = "/opt/imls/auth.yaml"

func SetAuthPath(path string) {
	authPath = path
}
func GetAuthPath() string {
	return authPath
}

// Located at /opt/imls/auth.yaml
type AuthConfig struct {
	Token     string `yaml:"api_token"`
	DeviceTag string `yaml:"tag"`
	FCFSId    string `yaml:"fcfs_seq_id"`
}

func ReadAuth() (*AuthConfig, error) {
	_, err := os.Stat(GetAuthPath())
	if err != nil {
		return &AuthConfig{}, fmt.Errorf("readToken: cannot find default token file at [%v]", GetAuthPath())
	}

	f, err := os.Open(GetAuthPath())
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
		log.Println("readToken: failed to decrypt auth token after decoding")
	}
	auth.Token = string(dec)
	return auth, nil
}

func ReadAuthTest() (*AuthConfig, error) {
	_, err := os.Stat(GetAuthPath())
	if err != nil {
		return &AuthConfig{}, fmt.Errorf("readToken: cannot find default token file at [%v]", GetAuthPath())
	}

	f, err := os.Open(GetAuthPath())
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
	// dec, err := cryptopasta.Decrypt(b64, &key)
	// if err != nil {
	// 	log.Fatal("readToken: failed to decrypt auth token after decoding")
	// }
	auth.Token = string(b64)
	return auth, nil
}
