package config

import (
	"flag"
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v2"
	"gsa.gov/18f/session-counter/constants"
	"gsa.gov/18f/session-counter/model"
)

/* FUNC checkEnvVars
 * Checks to see if the username and password for
 * working with Directus is in memory.
 * If not, it quits.
 */
func checkEnvVars() {
	if os.Getenv(constants.EnvUsername) == "" {
		fmt.Printf("%s must be set in the env!\n", constants.EnvUsername)
		os.Exit(constants.ExitNoUsername)
	}
	if os.Getenv(constants.EnvPassword) == "" {
		fmt.Printf("%s must be set in the env!\n", constants.EnvPassword)
		os.Exit(constants.ExitNoPassword)
	}
}

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
	checkEnvVars()
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

func ReadAuth() (*model.AuthConfig, error) {
	_, err := os.Stat(constants.AuthPath)
	if err != nil {
		return &model.AuthConfig{}, fmt.Errorf("readToken: cannot find default token file at [%v]", constants.AuthPath)
	}

	f, err := os.Open(constants.AuthPath)
	if err != nil {
		log.Fatal("readToken: could not open token file. Exiting.")
	}
	defer f.Close()
	var auth *model.AuthConfig
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&auth)
	if err != nil {
		log.Fatalf("readToken: could not decode YAML:\n%v\n", err)
	}

	return auth, nil
}

func ReadConfig() *Config {
	// We expect config to be here:
	//   * /etc/session-counter/config.yaml
	// We expect there to be a token file at
	//   * /etc/session-counter/access-token
	//
	// If neither of those is true, we can check for a
	// the username and password to be in the ENV, and
	// for the config to be passed via command line.

	cfg, err := parseConfigFile(constants.ConfigPath)
	if err != nil {
		fmt.Printf("config: could not find config at default path [%v]\n", constants.ConfigPath)
		fmt.Println("config: loading dev config")
		return devConfig()
	}

	// FIXME 20210323 MCJ
	// ARE THESE USED ANYWHERE?

	// auth, err := ReadAuth()
	// if err != nil {
	// 	log.Fatal("readConfig: cannot find auth token")
	// }

	// Stick the username/token into the environment.
	// This will be used by the Get_token() auth dance.
	// os.Setenv(constants.AuthTokenKey, auth.Token)
	// os.Setenv(constants.AuthEmailKey, auth.User)

	return cfg
}
