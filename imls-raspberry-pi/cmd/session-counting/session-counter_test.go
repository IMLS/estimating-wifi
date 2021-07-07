package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"testing"
	"time"

	"gsa.gov/18f/analysis"
	"gsa.gov/18f/config"
	"gsa.gov/18f/logwrapper"
	"gsa.gov/18f/session-counter/model"
	"gsa.gov/18f/session-counter/tlp"
)

const PASS = true
const FAIL = false

func macs(arr ...string) []string {
	// h := make(map[string]int)
	// for _, s := range arr {
	// 	h[s] = rand.Intn(1024)
	// }
	// return h
	return arr
}

func hashes(arr ...string) [][]string {
	// Return a list of hashes, one hash for each string
	harr := make([][]string, 0)
	for _, s := range arr {
		harr = append(harr, []string{s})
	}
	return harr
}

// IDs will be assigned in the raw_to_uids proc
// on a sorted list of MAC addrs. Therefore, we *know*
// in this test set that "next" will always be UID 0.
// (If "next" and "apple" are together)
var m = map[string]string{
	"next":     "00:00:0f:aa:bb:cc", // ID 0
	"ericsson": "00:01:ec:aa:bb:cc",
	"apple":    "00:03:93:aa:bb:cc",
	"next2":    "00:00:0f:ee:ff:00",
}

var tests = []struct {
	description       string
	passfail          bool
	uniqueness_window int
	initMap           []string
	loopMaps          [][]string
	resultMap         map[string]int
}{
	// One input hash.
	{"one input mac, one loop mac",
		PASS, 5,
		macs(m["next"]),
		hashes(m["next"]),
		map[string]int{
			"0:0": 0,
		},
	},
	// // Two input hashes
	{"two input macs, one loop mac",
		PASS, 5,
		macs(m["next"], m["apple"]),
		hashes(m["next"]),
		// Why zero and one?
		// Zero for deadbeef, because we send it in the loop.
		// One for beefcafe, because it was only sent once, and
		// one tick goes by.
		map[string]int{
			"0:0": 0,
			"1:1": 1,
		},
	},
	// // Two input hashes
	{"two input macs, one loop mac, both next",
		PASS, 5,
		macs(m["next"], m["next2"]),
		hashes(m["next"]),
		// Why zero and one?
		// Zero for deadbeef, because we send it in the loop.
		// One for beefcafe, because it was only sent once, and
		// one tick goes by.
		map[string]int{
			"0:0": 0,
			"0:1": 1,
		},
	},
	// Three hashes, three minutes
	{"three input macs, three comms in the middle",
		PASS, 5,
		// Next, Apple, Ericsson
		macs(m["next"], m["apple"], m["ericsson"]),
		hashes("de:ad:be:ef", "de:ad:be:ef", "de:ad:be:ef"),
		// IDs will be assigned by MAC address sort!
		map[string]int{
			"0:0": 3,
			"1:1": 3,
			"2:2": 3,
			"3:3": 0,
		},
	},

	// Next times out, because it is considered to
	// have "disconnected" after 5 minutes.
	{"Next should disappear",
		PASS, 5,
		macs(m["next"], m["apple"], m["ericsson"]),
		hashes(
			"de:ad:be:ef",
			"de:ad:be:ef",
			"de:ad:be:ef",
			m["apple"],
			m["ericsson"]),
		// Why zero and one?
		// Zero for deadbeef, because we send it in the loop.
		// One for beefcafe, because it was only sent once, and
		// one tick goes by.
		map[string]int{
			"1:1": 1,
			"2:2": 0,
			"3:3": 2,
		},
	},

	// Next times out, comes back. Still ID 0.
	// Apple is considered to have disconnected.
	{"Drop two",
		PASS, 5,
		macs(m["next"], m["apple"], m["ericsson"]),
		hashes(
			"de:ad:be:ef",
			"de:ad:be:ef",
			"de:ad:be:ef",
			"de:ad:be:ef",
			m["next"]),
		// Why zero and one?
		// Zero for deadbeef, because we send it in the loop.
		// One for beefcafe, because it was only sent once, and
		// one tick goes by.
		map[string]int{
			"0:0": 0,
			"3:3": 1,
		},
	},
}

func assertEqual(t *testing.T, a interface{}, b interface{}, message string) {
	if a == b {
		return
	}
	t.Fatal(message, "\n\texpected: ", a, "\n\treceived: ", b)
}

func TestRawToUid(t *testing.T) {
	log.Println("TestRawToUid")
	cfg := new(config.Config)
	logwrapper.NewLogger(cfg)

	_, filename, _, _ := runtime.Caller(0)
	fmt.Println(filename)
	path := filepath.Dir(filename)
	cfg.Manufacturers.Db = filepath.Join(path, "test", "manufacturers.sqlite")

	ka := tlp.NewKeepalive(cfg)

	// var buf bytes.Buffer
	// log.SetOutput(&buf)
	// defer func() {
	// 	log.SetOutput(os.Stderr)
	// }()

	for testNdx, e := range tests {
		t.Logf("Test #%v: %v\n", testNdx, e.description)
		cfg.Monitoring.UniquenessWindow = e.uniqueness_window

		var wg sync.WaitGroup
		resetbroker := tlp.NewBroker()
		go resetbroker.Start()
		// The kill broker lets us poison the network.
		// var killbroker *tlp.Broker = nil
		killbroker := tlp.NewBroker()
		go killbroker.Start()

		ch_macs := make(chan []string)
		ch_uniq := make(chan map[string]int)
		var u map[string]int = nil

		wg.Add(1)
		go func() {
			ch_macs <- e.initMap
			for _, sarr := range e.loopMaps {

				ch_macs <- sarr
			}
			defer wg.Done()
		}()

		// Not using the reset here.
		go tlp.AlgorithmTwo(ka, cfg, resetbroker, killbroker, ch_macs, ch_uniq)

		wg.Add(1)
		go func() {
			// The init map
			<-ch_uniq
			count := len(e.loopMaps) - 1
			for i := 0; i < count; i++ {
				// This reads in the intervening maps.
				<-ch_uniq
			}
			u = <-ch_uniq
			// ch_poison <- tlp.Ping{}
			killbroker.Publish(tlp.Ping{})
			defer wg.Done()
		}()

		wg.Wait()

		// The last value we receive needs to have its time updated.
		expected := fmt.Sprint(e.resultMap)
		received := fmt.Sprint(u)
		//log.Println("expected", expected, "received", received)

		if e.passfail {
			assertEqual(t, expected, received, "not equal")
		}
	} // end for over tests
}

func PingAfterNHours(ka *tlp.Keepalive, cfg *config.Config, rb *tlp.Broker, kb *tlp.Broker, n_hours int, ch_tick chan bool) {
	counter := 0
	ch_kill := kb.Subscribe()
	for {
		select {
		case <-ch_tick:
			counter += 1
			// Ping every
			if (counter != 0) && ((counter % (60 * n_hours)) == 0) {
				// ch_reset <- tlp.Ping{}
				rb.Publish(tlp.Ping{})
			}
		case <-ch_kill:
			log.Println("Exiting PingAfterNHours")
			return
		}
	}
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

func generateFakeMac() string {
	var letterRunes = []rune("ABCDEF0123456789")
	b := make([]rune, 17)
	colons := [...]int{2, 5, 8, 11, 14}
	for i := 0; i < 17; i++ {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]

		for v := range colons {
			if i == colons[v] {
				b[i] = rune(':')
			}
		}
	}
	return string(b)
}

func RunFakeWireshark(ka *tlp.Keepalive, cfg *config.Config, kb *tlp.Broker, in <-chan bool, out chan []string) {
	NUMMACS := 40
	NUMRANDOM := 10
	ch_kill := kb.Subscribe()
	// Lets have 30 consistent devices
	macs := make([]string, NUMMACS)
	for i := 0; i < NUMMACS-NUMRANDOM; i++ {
		macs[i] = generateFakeMac()
	}

	for {
		select {
		case <-in:
			// And 10 random devices
			for i := 0; i < NUMRANDOM; i++ {
				macs[30+i] = generateFakeMac()
			}
			out <- macs

		case <-ch_kill:
			log.Println("Exiting RunFakeWireshark")
			return
		}
	}

}

func TestManyTLPCycles(t *testing.T) {
	const NUMDAYSTORUN int = 1
	const NUMMINUTESTORUN int = NUMDAYSTORUN * 24 * 60
	const WRITESUMMARYNHOURS int = 1
	const SECONDSPERMINUTE int = 3
	resetbroker := tlp.NewBroker()
	go resetbroker.Start()
	killbroker := tlp.NewBroker()
	go killbroker.Start()

	// This runs the TLP through 10000 cycles. This is roughly the same as week.
	log.Println("Starting run for a week")

	// Get a local config, so we have paths...
	_, filename, _, _ := runtime.Caller(0)
	fmt.Println(filename)
	path := filepath.Dir(filename)

	configPath := filepath.Join(path, "test", "config.yaml")
	cfg, _ := config.NewConfigFromPath(configPath)
	cfg.Local.SummaryDB = filepath.Join(path, "summarydb.sqlite")
	cfg.Manufacturers.Db = filepath.Join(path, "test", "manufacturers.sqlite")
	cfg.Local.WebDirectory = filepath.Join(path, "test", "www")
	tt := time.Now()
	cfg.SessionId = fmt.Sprintf("%v%02d%02d", tt.Year(), tt.Month(), tt.Day())

	os.Mkdir(cfg.Local.WebDirectory, 0755)

	//cfg.NewSessionId()
	log.Println(cfg)

	// Create channels for process network
	ch_sec := make(chan bool)

	ch_nsec := make(chan bool)
	ch_nsec1 := make(chan bool)
	ch_nsec2 := make(chan bool)

	ch_macs := make(chan []string)
	ch_macs_counted := make(chan map[string]int)
	ch_data_for_report := make(chan []analysis.WifiEvent)
	ch_db := make(chan *model.TempDB)
	ch_durations_db := make(chan *model.TempDB)

	// See if we can wait and shut down the test...
	var wg sync.WaitGroup
	wg.Add(1)

	// Delta this out to RunWireshark and PingAfterNHours
	go tlp.TockEveryN(nil, killbroker, SECONDSPERMINUTE, ch_sec, ch_nsec)
	go func(in chan bool, o1 chan bool, o2 chan bool) {
		for {
			<-in
			o1 <- true
			o2 <- true
		}
	}(ch_nsec, ch_nsec1, ch_nsec2)

	// Need a fake RunWireshark
	// go tlp.RunWireshark(nil, cfg, ch_nsec1, ch_macs, KC[1])
	go RunFakeWireshark(nil, cfg, killbroker, ch_nsec1, ch_macs)

	go tlp.AlgorithmTwo(nil, cfg, resetbroker, killbroker, ch_macs, ch_macs_counted)
	go tlp.PrepEphemeralWifi(nil, cfg, killbroker, ch_macs_counted, ch_data_for_report)
	// At midnight, flush internal structures and restart.
	//go tlp.PingAtMidnight(nil, cfg, chs_reset[0], KC[4])
	go PingAfterNHours(nil, cfg, resetbroker, killbroker, WRITESUMMARYNHOURS, ch_nsec2)
	go tlp.CacheWifi(nil, cfg, resetbroker, killbroker, ch_data_for_report, ch_db)
	// Make sure we don't hang...
	go tlp.GenerateDurations(nil, cfg, killbroker, ch_db, ch_durations_db)
	go tlp.BatchSend(nil, cfg, killbroker, ch_durations_db)

	// We want 10000 minutes, but the tocker is every second.
	go func() {
		// Give the rest of the network time to come alive.
		time.Sleep(5 * time.Second)
		minutes := 0
		for secs := 0; secs < NUMMINUTESTORUN*60; secs++ {
			ch_sec <- true
			if secs%SECONDSPERMINUTE == 0 {
				var m runtime.MemStats
				runtime.ReadMemStats(&m)
				minutes += 1
				hours := (minutes / 60) % 24
				days := (minutes / (60 * 24))
				memstats := fmt.Sprintf("Alloc[%vMB] Sys[%vMB], NumGC[%v]", bToMb(m.Alloc), bToMb(m.Sys), m.NumGC)
				log.Println(days, "d", hours, "h", minutes%60, "m", memstats)
			}
		}
		log.Println("Killing the test network.")
		killbroker.Publish(tlp.Ping{})
		wg.Done()
	}()
	wg.Wait()
	log.Println("Done waiting... exiting in 10 seconds")
	time.Sleep(10 * time.Second)
}
