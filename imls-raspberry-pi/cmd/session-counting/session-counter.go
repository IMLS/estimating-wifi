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

	// WARNING: If you get this length wrong, we have deadlock.
	// That is, every one of these needs to be used/written to/read from.
	const RESET_CHANS = 3
	// The kill channel lets us poison the network for shutdown. Really only for testing.
	var NIL_KILL_CHANNEL chan tlp.Ping = nil

	var chs_reset [RESET_CHANS]chan tlp.Ping
	for ndx := 0; ndx < RESET_CHANS; ndx++ {
		chs_reset[ndx] = make(chan tlp.Ping)
	}

	// Run the process network.
	// Driven by a 1s `tick` process.
	// Thread the keepalive through the network
	go tlp.TockEveryMinute(ka, ch_nsec, NIL_KILL_CHANNEL)
	go tlp.RunWireshark(ka, cfg, ch_nsec, ch_macs, NIL_KILL_CHANNEL)
	// The reset will never be triggered in AlgoTwo unless we're rnuning in "sqlite" storage mode.
	go tlp.AlgorithmTwo(ka, cfg, ch_macs, ch_macs_counted, chs_reset[1], NIL_KILL_CHANNEL)
	go tlp.PrepareDataForStorage(ka, cfg, ch_macs_counted, ch_data_for_report, NIL_KILL_CHANNEL)

	// We need a multiplexer of sorts that will route data to the appropriate
	// storage processes. The storage processes can then sit there waiting.
	// go tlp.Multiplexer(ch_data_for_report, ch_out_to_api, ch_out_to_sqlite)
	go tlp.StoreToCloud(ka, cfg, ch_data_for_report, chs_reset[2], NIL_KILL_CHANNEL)

	// We need to think about how our state is managed; currently, it is burried
	// in StoreToSqlite, but really, that temporary state is not unique to that
	// backend. Perhaps the temporary state is a process that lives between
	// PrepareDataForStorage and the storage procs?  (Sorry... proc == gofunc)
	// Or, perhaps the state management moves back into PrepareDataForStorage?
	// Either way, the kind of state we need to reset is... well, it's threaded
	// through the network. Anything that listens to, and takes action on, a Ping
	// on the chs_reset[] lines, is something that we need to think about.

	// As I write this... it might be that StoreToSqlite becomes the state manager.
	// Pull out the (tiny) bit that actually writes data into its own process?

	// Another thing... and I hate to mention this...
	// should we consider moving from live storage of data to the cloud, and instead to a batch
	// store every night. That way, we *always* store to Sqlite, and every night,
	// we try and figure out what we have and have not submitted. This way,
	// if there are network problems, we don't lose data. Instead, we keep track
	// (locally) what sessions (days) have been transmitted, and which havent, and
	// when we get through (say there's a network outage), we just send everything that hasn't
	// been sent prior.

	// Or, something to that effect. It might be a "next step" kinda thing. But,
	// given that we have the infra in place to write wifi data to temporary SQLite
	// dbs in the fs already... it's not actually that big a jump to write that data in
	// a batch mode (vs. daily/live mode)...

	// Writes a Ping{} to chs_reset[0]
	// it would be more readable to have the output channel be named,
	// and only have the Par outputs be an array. I cheated...
	go tlp.PingAtMidnight(ka, cfg, chs_reset[0], NIL_KILL_CHANNEL)
	// Listens for a ping to know when to reset internal state.
	// That, too, should be abstracted out of the storage layer.
	go tlp.StoreToSqlite(ka, cfg, ch_data_for_report, chs_reset[2], NIL_KILL_CHANNEL)
	// Fan out the ping to multiple PROCs
	go tlp.ParDelta(NIL_KILL_CHANNEL, chs_reset[0], chs_reset[1:]...)

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
