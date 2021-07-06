package main

import (
	"flag"
	"fmt"
	"os"
	"sync"
	"time"

	"gsa.gov/18f/config"
	"gsa.gov/18f/logwrapper"
	"gsa.gov/18f/session-counter/api"
	"gsa.gov/18f/session-counter/model"
	"gsa.gov/18f/session-counter/tlp"
	"gsa.gov/18f/version"
)

func run(ka *tlp.Keepalive, cfg *config.Config) {
	logwrapper.NewLogger(nil)

	// Create channels for process network
	// ch_sec := make(chan bool)
	ch_nsec := make(chan bool)
	ch_macs := make(chan []string)
	ch_macs_counted := make(chan map[string]int)
	ch_data_for_report := make(chan []map[string]string)
	ch_db := make(chan *model.TempDB)

	// The reset broker signals midnight (for resetting the network/device)
	resetbroker := tlp.NewBroker()
	go resetbroker.Start()
	// The kill broker lets us poison the network.
	var killbroker *tlp.Broker = nil

	// Run the process network.
	// Driven by a 1s `tick` process.
	// Thread the keepalive through the network
	go tlp.TockEveryMinute(ka, killbroker, ch_nsec)
	go tlp.RunWireshark(ka, cfg, killbroker, ch_nsec, ch_macs)
	// The reset will never be triggered in AlgoTwo unless we're rnuning in "sqlite" storage mode.
	go tlp.AlgorithmTwo(ka, cfg, killbroker, resetbroker, ch_macs, ch_macs_counted)
	go tlp.PrepEphemeralWifi(ka, cfg, killbroker, ch_macs_counted, ch_data_for_report)

	go tlp.CacheWifi(ka, cfg, resetbroker, killbroker, ch_data_for_report, ch_db)

	go tlp.PingAtMidnight(ka, cfg, resetbroker, killbroker)
	// Listens for a ping to know when to reset internal state.
	// That, too, should be abstracted out of the storage layer.
	//go tlp.StoreToSqlite(ka, cfg, resetbroker, killbroker, ch_data_for_report)

}

func keepalive(ka *tlp.Keepalive, cfg *config.Config) {
	lw := logwrapper.NewLogger(nil)
	lw.Info("starting keepalive")
	var counter int64 = 0
	for {
		time.Sleep(time.Duration(cfg.Monitoring.PingInterval) * time.Second)
		ka.Publish(counter)
		counter = counter + 1
	}
}

func handleFlags() *config.Config {
	versionPtr := flag.Bool("version", false, "Get the software version and exit.")
	showKeyPtr := flag.Bool("show-key", false, "Tests key decryption.")
	configPathPtr := flag.String("config", "", "Path to config.yaml. REQUIRED.")
	flag.Parse()
	lw := logwrapper.NewLogger(nil)

	// If they just want the version, print and exit.
	if *versionPtr {
		fmt.Println(version.GetVersion())
		os.Exit(0)
	}

	// Make sure a config is passed.
	if *configPathPtr == "" {
		lw.Fatal("The flag --config MUST be provided.")
		os.Exit(1)
	}

	if _, err := os.Stat(*configPathPtr); os.IsNotExist(err) {
		lw.Info("Looked for config at: %v", *configPathPtr)
		lw.Fatal("Cannot find config file. Exiting.")
	}

	cfg, err := config.NewConfigFromPath(*configPathPtr)
	if err != nil {
		lw.Fatal("session-counter: error loading config.")
	}

	if *showKeyPtr {
		fmt.Println(cfg.Auth.Token)
		os.Exit(0)
	}

	return cfg

}

func main() {
	// Read in a config
	cfg := handleFlags()

	cfg.NewSessionId()

	lw := logwrapper.NewLogger(cfg)
	lw.Info("startup")

	// Store this so we don't keep hitting /proc/cpuinfo
	cfg.Serial = config.GetSerial()
	// Make sure the mfg database is in place and can be loaded.
	api.CheckMfgDatabaseExists(cfg)

	// also make sure the binary paths in the config are valid.
	_, err := os.Stat(cfg.Wireshark.Path)
	if os.IsNotExist(err) {
		lw.ExeNotFound(cfg.Wireshark.Path)
	}

	ka := tlp.NewKeepalive(cfg)
	go ka.Start()
	go keepalive(ka, cfg)
	go run(ka, cfg)

	// Wait forever.
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
