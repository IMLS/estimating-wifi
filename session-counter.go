package main

import (
	"crypto/sha256"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"gopkg.in/yaml.v2"
	"gsa.gov/18f/session-counter/constants"
	"gsa.gov/18f/session-counter/csp"
	"gsa.gov/18f/session-counter/model"
	"gsa.gov/18f/session-counter/tlp"
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

func parseConfigFile(filepath string) (*model.Config, error) {
	_, err := os.Stat(filepath)

	// Stat will set an error if the file cannot be found.
	if err == nil {
		f, err := os.Open(filepath)
		if err != nil {
			log.Fatal("parseConfigFile: could not open configuration file. Exiting.")
		}
		defer f.Close()
		var cfg *model.Config
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
func devConfig() *model.Config {
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

func readAuth() (*model.Auth, error) {
	_, err := os.Stat(constants.AuthPath)
	if err != nil {
		return &model.Auth{}, fmt.Errorf("readToken: cannot find default token file at [%v]", constants.AuthPath)
	}

	f, err := os.Open(constants.AuthPath)
	if err != nil {
		log.Fatal("readToken: could not open token file. Exiting.")
	}
	defer f.Close()
	var auth *model.Auth
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&auth)
	if err != nil {
		log.Fatalf("readToken: could not decode YAML:\n%v\n", err)
	}

	return auth, nil
}

func readConfig() *model.Config {
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

	auth, err := readAuth()
	if err != nil {
		log.Fatal("readConfig: cannot find auth token")
	}

	os.Setenv(constants.AuthTokenKey, auth.Token)
	os.Setenv(constants.AuthEmailKey, auth.Email)
	return cfg
}

func run(ka *csp.Keepalive, cfg *model.Config) {
	log.Println("Starting run")
	// Create channels for process network
	ch_sec := make(chan bool)
	ch_nsec := make(chan bool)
	ch_macs := make(chan map[string]int)
	ch_macs_counted := make(chan map[string]int)
	ch_mfg := make(chan map[string]model.Entry)

	// Run the process network.
	// Driven by a 1s `tick` process.
	// Thread the keepalive through the network
	go tlp.Tick(ka, ch_sec)
	go tlp.TockEveryN(ka, 60, ch_sec, ch_nsec)
	go tlp.RunWireshark(ka, cfg, ch_nsec, ch_macs)
	go tlp.MacToEntry(ka, cfg, ch_macs_counted, ch_mfg)
	go tlp.RingBuffer(ka, cfg, ch_macs, ch_macs_counted)
	go tlp.ReportMap(ka, cfg, ch_mfg)
}

func keepalive(ka *csp.Keepalive, cfg *model.Config) {
	log.Println("Starting keepalive")
	var counter int64 = 0
	for {
		time.Sleep(time.Duration(cfg.Monitoring.PingInterval) * time.Second)
		ka.Publish(counter)
		counter = counter + 1
	}
}

func calcSessionId() string {
	h := sha256.New()
	email := os.Getenv(constants.AuthEmailKey)
	// FIXME: Use the email instead of the token.
	// Guaranteed to be unique. Current time along with our auth token, hashed.
	h.Write([]byte(fmt.Sprintf("%v%x", time.Now(), email)))
	sid := fmt.Sprintf("%x", h.Sum(nil))
	log.Println("Session id: ", sid)
	return sid
}

func main() {
	// Read in a config
	cfg := readConfig()
	// Add a "sessionId" to the mix.
	cfg.SessionId = calcSessionId()

	ka := csp.NewKeepalive()
	go ka.Start()
	go keepalive(ka, cfg)
	go run(ka, cfg)

	// Wait forever.
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
