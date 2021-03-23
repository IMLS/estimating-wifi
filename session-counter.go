package main

import (
	"crypto/sha256"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"gsa.gov/18f/session-counter/config"
	"gsa.gov/18f/session-counter/constants"
	"gsa.gov/18f/session-counter/csp"
	"gsa.gov/18f/session-counter/model"
	"gsa.gov/18f/session-counter/tlp"
)

func run(ka *csp.Keepalive, cfg *config.Config) {
	log.Println("Starting run")
	// Create channels for process network
	ch_sec := make(chan bool)
	ch_nsec := make(chan bool)
	ch_macs := make(chan []string)
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

func keepalive(ka *csp.Keepalive, cfg *config.Config) {
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
	sid := fmt.Sprintf("%x", h.Sum(nil))[0:8]
	// Keep it short.
	log.Println("Session id: ", sid)
	return sid
}

func main() {
	// Read in a config
	cfg := config.ReadConfig()
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
