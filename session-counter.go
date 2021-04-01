package main

import (
	"log"
	"sync"
	"time"

	"gsa.gov/18f/session-counter/api"
	"gsa.gov/18f/session-counter/config"
	"gsa.gov/18f/session-counter/tlp"
)

func run(ka *tlp.Keepalive, cfg *config.Config) {
	log.Println("Starting run")
	// Create channels for process network
	ch_sec := make(chan bool)
	ch_nsec := make(chan bool)
	ch_macs := make(chan []string)
	ch_macs_counted := make(chan map[string]int)

	// Run the process network.
	// Driven by a 1s `tick` process.
	// Thread the keepalive through the network
	go tlp.Tick(ka, ch_sec)
	go tlp.TockEveryN(ka, 60, ch_sec, ch_nsec)
	go tlp.RunWireshark(ka, cfg, ch_nsec, ch_macs)
	go tlp.AlgorithmTwo(ka, cfg, ch_macs, ch_macs_counted, nil)
	go tlp.ReportOut(ka, cfg, ch_macs_counted)
}

func keepalive(ka *tlp.Keepalive, cfg *config.Config) {
	log.Println("Starting keepalive")
	var counter int64 = 0
	for {
		time.Sleep(time.Duration(cfg.Monitoring.PingInterval) * time.Second)
		ka.Publish(counter)
		counter = counter + 1
	}
}

func main() {
	// Read in a config
	cfg := config.ReadConfig()
	// Set the session ID for this entire run
	cfg.SessionId = config.CreateSessionId()
	// Store this so we don't keep hitting /proc/cpuinfo
	cfg.Serial = config.GetSerial()
	// Make sure the mfg database is in place and can be loaded.
	api.CheckMfgDatabaseExists(cfg)

	el := api.NewEventLogger(cfg)
	el.Log("startup", nil)

	ka := tlp.NewKeepalive(cfg)
	go ka.Start()
	go keepalive(ka, cfg)
	go run(ka, cfg)

	// Wait forever.
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
