package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sync"

	"gsa.gov/18f/cmd/session-counter/tlp"
	"gsa.gov/18f/internal/interfaces"
	"gsa.gov/18f/internal/logwrapper"
	"gsa.gov/18f/internal/state"
	"gsa.gov/18f/internal/structs"
	"gsa.gov/18f/internal/version"
)

func run(cfg interfaces.Config) {
	logwrapper.NewLogger(nil)

	// CHANNELS
	chNsec := make(chan bool)
	chMacs := make(chan []string)
	chMacsCounted := make(chan map[string]int)
	chDataForReport := make(chan []structs.WifiEvent)
	chWifidb := make(chan interfaces.Database)
	chDdb := make(chan interfaces.Database)
	chDdbPar := make([]chan interfaces.Database, 2)
	chAck := make(chan tlp.Ping)
	for i := 0; i < 2; i++ {
		chDdbPar[i] = make(chan interfaces.Database)
	}
	// BROKERS
	resetbroker := tlp.NewResetBroker()
	go resetbroker.Start()
	var killbroker *tlp.KillBroker = nil
	ka := tlp.NewKeepalive()

	// PROCESSES
	go tlp.StayinAlive(ka)
	go tlp.TockEveryMinute(ka, killbroker, chNsec)
	go tlp.RunWireshark(ka, killbroker, chNsec, chMacs)
	go tlp.AlgorithmTwo(ka, resetbroker, killbroker, chMacs, chMacsCounted)
	go tlp.PrepEphemeralWifi(ka, killbroker, chMacsCounted, chDataForReport)

	go tlp.CacheWifi(ka, resetbroker, killbroker, chDataForReport, chWifidb, chAck)
	go tlp.GenerateDurations(ka, killbroker, chWifidb, chDdb, chAck)

	go tlp.ParDeltaTempDB(killbroker, chDdb, chDdbPar...)
	go tlp.BatchSend(ka, killbroker, chDdbPar[0])
	go tlp.WriteImages(ka, killbroker, chDdbPar[1])

	go tlp.PingAtMidnight(ka, resetbroker, killbroker)
}

func initConfigFromFlags() {
	versionPtr := flag.Bool("version", false, "Get the software version and exit.")
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

	state.NewConfigFromPath(*configPathPtr)

}

func main() {
	// DO NOT USE LOGGING YET
	initConfigFromFlags()
	cfg := state.GetConfig()

	// NOW YOU MAY USE LOGGING.
	cfg.Log().Debug("session id at startup is is ", cfg.GetCurrentSessionId())

	// Make sure the mfg database is in place and can be loaded.
	// api.CheckMfgDatabaseExists(cfg)
	// also make sure the binary paths in the config are valid.
	// _, err := os.Stat(cfg.Wireshark.Path)
	// if os.IsNotExist(err) {
	// 	//lw.ExeNotFound(cfg.Wireshark.Path)
	// 	cfg.Log().Fatal("wireshark not found at ", cfg.Wireshark.Path)
	// }

	// Run the network
	var wg sync.WaitGroup
	wg.Add(1)
	go run(cfg)
	// Wait forever.
	wg.Wait()
}
