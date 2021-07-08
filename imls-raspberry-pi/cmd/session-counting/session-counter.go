package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sync"

	"gsa.gov/18f/analysis"
	"gsa.gov/18f/config"
	"gsa.gov/18f/logwrapper"
	"gsa.gov/18f/session-counter/api"
	"gsa.gov/18f/session-counter/model"
	"gsa.gov/18f/session-counter/tlp"
	"gsa.gov/18f/version"
)

func run(cfg *config.Config) {
	logwrapper.NewLogger(nil)
	// CHANNELS
	ch_nsec := make(chan bool)
	ch_macs := make(chan []string)
	ch_macs_counted := make(chan map[string]int)
	ch_data_for_report := make(chan []analysis.WifiEvent)
	ch_wifidb := make(chan *model.TempDB)
	ch_ddb := make(chan *model.TempDB)
	ch_ddb_par := make([]chan *model.TempDB, 2)
	for i := 0; i < 2; i++ {
		ch_ddb_par[i] = make(chan *model.TempDB)
	}
	// BROKERS
	resetbroker := &tlp.ResetBroker{tlp.NewBroker()}
	go resetbroker.Start()
	var killbroker *tlp.KillBroker = nil
	ka := tlp.NewKeepalive(cfg)

	// PROCESSES
	go tlp.StayinAlive(ka, cfg)
	go tlp.TockEveryMinute(ka, killbroker, ch_nsec)
	go tlp.RunWireshark(ka, cfg, killbroker, ch_nsec, ch_macs)
	go tlp.AlgorithmTwo(ka, cfg, resetbroker, killbroker, ch_macs, ch_macs_counted)
	go tlp.PrepEphemeralWifi(ka, cfg, killbroker, ch_macs_counted, ch_data_for_report)
	go tlp.CacheWifi(ka, cfg, resetbroker, killbroker, ch_data_for_report, ch_wifidb)

	go tlp.GenerateDurations(ka, cfg, killbroker, ch_wifidb, ch_ddb)
	go tlp.ParDeltaTempDB(killbroker, ch_ddb, ch_ddb_par...)
	go tlp.BatchSend(ka, cfg, killbroker, ch_ddb_par[0])
	go tlp.WriteImages(ka, cfg, killbroker, ch_ddb_par[1])

	go tlp.PingAtMidnight(ka, cfg, resetbroker, killbroker)
}

func handleFlags() *config.Config {
	versionPtr := flag.Bool("version", false, "Get the software version and exit.")
	showKeyPtr := flag.Bool("show-key", false, "Tests key decryption.")
	configPathPtr := flag.String("config", "", "Path to config.yaml. REQUIRED.")
	flag.Parse()

	// If they just want the version, print and exit.
	if *versionPtr {
		fmt.Println(version.GetVersion())
		os.Exit(0)
	}

	// Make sure a config is passed.
	if *configPathPtr == "" {
		log.Fatal("The flag --config MUST be provided.")
		os.Exit(1)
	}

	if _, err := os.Stat(*configPathPtr); os.IsNotExist(err) {
		log.Println("Looked for config at ", *configPathPtr)
		log.Fatal("Cannot find config file. Exiting.")
	}

	cfg, err := config.NewConfigFromPath(*configPathPtr)
	if err != nil {
		log.Fatal("session-counter: error loading config.")
	}

	if *showKeyPtr {
		fmt.Println(cfg.Auth.Token)
		os.Exit(0)
	}

	return cfg

}

func main() {
	// DO NOT USE LOGGING YET
	cfg := handleFlags()
	cfg.NewSessionId()
	// INIT THE LOGGER
	lw := logwrapper.NewLogger(cfg)
	// NOW YOU MAY USE LOGGING.

	cfg.DecodeSerial()
	// SINGLETON PATTERN
	// Once this is set up, all loggers (should)
	// log through the config passed here.
	lw.Info("startup")
	lw.Info("serial ", cfg.GetSerial())

	// Make sure the mfg database is in place and can be loaded.
	api.CheckMfgDatabaseExists(cfg)
	// also make sure the binary paths in the config are valid.
	_, err := os.Stat(cfg.Wireshark.Path)
	if os.IsNotExist(err) {
		//lw.ExeNotFound(cfg.Wireshark.Path)
		lw.Fatal("wireshark not found at ", cfg.Wireshark.Path)
	}

	// Run the network
	var wg sync.WaitGroup
	wg.Add(1)
	go run(cfg)
	// Wait forever.
	wg.Wait()
}
