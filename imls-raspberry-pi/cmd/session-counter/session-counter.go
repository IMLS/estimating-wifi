package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sync"

	"gsa.gov/18f/cmd/session-counter/api"
	"gsa.gov/18f/cmd/session-counter/tlp"
	"gsa.gov/18f/internal/config"
	"gsa.gov/18f/internal/logwrapper"
	"gsa.gov/18f/internal/state"
	"gsa.gov/18f/internal/structs"
	"gsa.gov/18f/internal/version"
)

func run(cfg *config.Config) {
	logwrapper.NewLogger(nil)

	// CHANNELS
	chNsec := make(chan bool)
	chMacs := make(chan []string)
	chMacsCounted := make(chan map[string]int)
	chDataForReport := make(chan []structs.WifiEvent)
	chWifidb := make(chan *state.TempDB)
	chDdb := make(chan *state.TempDB)
	chDdbPar := make([]chan *state.TempDB, 2)
	chAck := make(chan tlp.Ping)
	for i := 0; i < 2; i++ {
		chDdbPar[i] = make(chan *state.TempDB)
	}
	// BROKERS
	resetbroker := tlp.NewResetBroker()
	go resetbroker.Start()
	var killbroker *tlp.KillBroker = nil
	ka := tlp.NewKeepalive(cfg)

	// PROCESSES
	go tlp.StayinAlive(ka, cfg)
	go tlp.TockEveryMinute(ka, killbroker, chNsec)
	go tlp.RunWireshark(ka, cfg, killbroker, chNsec, chMacs)
	go tlp.AlgorithmTwo(ka, cfg, resetbroker, killbroker, chMacs, chMacsCounted)
	go tlp.PrepEphemeralWifi(ka, cfg, killbroker, chMacsCounted, chDataForReport)

	go tlp.CacheWifi(ka, cfg, resetbroker, killbroker, chDataForReport, chWifidb, chAck)
	go tlp.GenerateDurations(ka, cfg, killbroker, chWifidb, chDdb, chAck)

	go tlp.ParDeltaTempDB(killbroker, chDdb, chDdbPar...)
	go tlp.BatchSend(ka, cfg, killbroker, chDdbPar[0])
	go tlp.WriteImages(ka, cfg, killbroker, chDdbPar[1])

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

	// INIT THE LOGGER
	// SINGLETON PATTERN
	// Once this is set up, all loggers (should)
	// log through the config passed here.
	lw := logwrapper.NewLogger(cfg)
	// NOW YOU MAY USE LOGGING.
	cfg.SessionId = state.GetInitialSessionID(cfg)
	cfg.SessionId = state.GetNextSessionID(cfg)
	lw.Debug("session id at startup is is ", cfg.SessionId)

	cfg.DecodeSerial()
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
