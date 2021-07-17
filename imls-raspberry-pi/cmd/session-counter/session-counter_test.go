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

	"github.com/benbjohnson/clock"

	"gsa.gov/18f/cmd/session-counter/tlp"
	"gsa.gov/18f/internal/config"
	"gsa.gov/18f/internal/logwrapper"
	"gsa.gov/18f/internal/state"
	"gsa.gov/18f/internal/structs"
)

const PASS = true
const FAIL = false

// Lets mock the clock for testing.
type Application struct {
	Clock clock.Clock
}

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
	description      string
	passfail         bool
	uniquenessWindow int
	initMap          []string
	loopMaps         [][]string
	resultMap        map[string]int
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
		cfg.Monitoring.UniquenessWindow = e.uniquenessWindow

		var wg sync.WaitGroup
		resetbroker := tlp.NewResetBroker()
		go resetbroker.Start()
		killbroker := tlp.NewKillBroker()
		go killbroker.Start()

		chMacs := make(chan []string)
		chUniq := make(chan map[string]int)
		var u map[string]int = nil

		wg.Add(1)
		go func() {
			chMacs <- e.initMap
			for _, sarr := range e.loopMaps {

				chMacs <- sarr
			}
			defer wg.Done()
		}()

		// Not using the reset here.
		go tlp.AlgorithmTwo(ka, cfg, resetbroker, killbroker, chMacs, chUniq)

		wg.Add(1)
		go func() {
			// The init map
			<-chUniq
			count := len(e.loopMaps) - 1
			for i := 0; i < count; i++ {
				// This reads in the intervening maps.
				<-chUniq
			}
			u = <-chUniq
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

func PingAfterNHours(ka *tlp.Keepalive, cfg *config.Config, rb *tlp.ResetBroker, kb *tlp.KillBroker, nHours int, chTick chan bool) {
	counter := 0
	chKill := kb.Subscribe()
	lw := logwrapper.NewLogger(nil)
	for {
		select {
		case <-chTick:
			counter += 1
			if (counter != 0) && ((counter % (60 * nHours)) == 0) {
				// chReset <- tlp.Ping{}
				lw.Debug("PingAfterNHours is Pinging!")
				rb.Publish(tlp.Ping{})
			}
		case <-chKill:
			log.Println("Exiting PingAfterNHours")
			return
		}
	}
}

func PingAtBogoMidnight(ka *tlp.Keepalive, cfg *config.Config,
	rb *tlp.ResetBroker,
	kb *tlp.KillBroker,
	m *clock.Mock) {
	// counter := 0
	// chKill := kb.Subscribe()
	lw := logwrapper.NewLogger(nil)
	pinged := false
	for {
		if m.Now().Hour() == 0 && !pinged {
			pinged = true
			lw.Debug("IT IS BOGOMIDNIGHT.")
			rb.Publish(tlp.Ping{})
		}
		if m.Now().Hour() != 0 {
			pinged = false
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

var consistentMacs = []string{
	"00:03:93:60:BF:CB",
	"00:03:93:06:5C:0E",
	"00:03:93:3F:BB:F9",
	"00:03:93:51:9D:26",
	"00:00:F0:D0:25:52",
	"00:00:F0:59:41:80",
	"00:00:F0:F2:2C:13",
}

func RunFakeWireshark(ka *tlp.Keepalive, cfg *config.Config, kb *tlp.KillBroker, in <-chan bool, out chan []string) {
	NUMMACS := 200
	NUMRANDOM := 10
	lw := logwrapper.NewLogger(nil)
	lw.Debug("RunFakeWireshark in the house.")

	chKill := kb.Subscribe()
	// Lets have 30 consistent devices
	macs := make([]string, NUMMACS)
	for i := 0; i < NUMMACS; i++ {
		macs[i] = generateFakeMac()
	}

	for {
		select {
		case <-in:
			// Pick NUMRANDOM devices every minute
			send := make([]string, NUMRANDOM)
			for i := 0; i < NUMRANDOM; i++ {
				send[i] = macs[rand.Intn(len(macs))]
			}
			out <- send

		case <-chKill:
			log.Println("Exiting RunFakeWireshark")
			return
		}
	}

}

func TestManyTLPCycles(t *testing.T) {
	const NUMDAYSTORUN int = 6
	const NUMMINUTESTORUN int = NUMDAYSTORUN * 24 * 60

	const skip int = 20
	const NUMCYCLESTORUN = NUMMINUTESTORUN / skip

	const SECONDSPERMINUTE int = 2
	lw := logwrapper.NewLogger(nil)
	lw.SetLogLevel("DEBUG")

	resetbroker := tlp.NewResetBroker()
	go resetbroker.Start()
	killbroker := tlp.NewKillBroker()
	go killbroker.Start()

	// This runs the TLP through 10000 cycles. This is roughly the same as week.
	log.Println("Starting run for a week")

	// Get a local config, so we have paths...
	_, filename, _, _ := runtime.Caller(0)
	fmt.Println(filename)
	path := filepath.Dir(filename)

	configPath := filepath.Join(path, "test", "config.yaml")
	cfg, _ := config.NewConfigFromPath(configPath)
	mock := clock.NewMock()
	mt, _ := time.Parse("2006-01-02T15:04", "1975-10-11T18:00")
	mock.Set(mt)
	cfg.Clock = mock

	if cfg.Clock == nil {
		log.Println("clock should not be nil")
		t.Fail()
	}
	lw.Debug("mock is now ", cfg.Clock.Now())

	cfg.RunMode = "test"
	cfg.StorageMode = "sqlite"
	cfg.Local.SummaryDB = filepath.Join(path, "summarydb.sqlite")
	cfg.Manufacturers.Db = filepath.Join(path, "test", "manufacturers.sqlite")
	cfg.Local.WebDirectory = filepath.Join(path, "test", "www")
	os.Mkdir(cfg.Local.WebDirectory, 0755)
	cfg.SessionId = state.GetNextSessionID(cfg)

	// Create channels for process network
	chSec := make(chan bool)

	chNsec := make(chan bool)
	chNsec1 := make(chan bool)
	chNsec2 := make(chan bool)

	chMacs := make(chan []string)
	chMacsCounted := make(chan map[string]int)
	chDataForReport := make(chan []structs.WifiEvent)
	chWifiDB := make(chan *state.TempDB)
	chDurationsDB := make(chan *state.TempDB)
	chAck := make(chan tlp.Ping)
	chDdbPar := make([]chan *state.TempDB, 2)
	for i := 0; i < 2; i++ {
		chDdbPar[i] = make(chan *state.TempDB)
	}

	// See if we can wait and shut down the test...
	var wg sync.WaitGroup
	wg.Add(1)

	// Delta this out to RunWireshark and PingAfterNHours
	go tlp.TockEveryN(nil, killbroker, SECONDSPERMINUTE, chSec, chNsec)
	go func(in chan bool, o1 chan bool, o2 chan bool) {
		for {
			<-in
			o1 <- true
			o2 <- true
		}
	}(chNsec, chNsec1, chNsec2)

	go func(in chan bool) {
		for {
			<-in
		}
	}(chNsec2)

	// Need a fake RunWireshark
	// go tlp.RunWireshark(nil, cfg, ch_nsec1, ch_macs, KC[1])
	go RunFakeWireshark(nil, cfg, killbroker, chNsec1, chMacs)

	go tlp.AlgorithmTwo(nil, cfg, resetbroker, killbroker, chMacs, chMacsCounted)
	go tlp.PrepEphemeralWifi(nil, cfg, killbroker, chMacsCounted, chDataForReport)
	// At midnight, flush internal structures and restart.
	//go tlp.PingAtMidnight(nil, cfg, chs_reset[0], KC[4])
	go PingAtBogoMidnight(nil, cfg, resetbroker, killbroker, mock)
	go tlp.CacheWifi(nil, cfg, resetbroker, killbroker, chDataForReport, chWifiDB, chAck)
	// Make sure we don't hang...
	go tlp.GenerateDurations(nil, cfg, killbroker, chWifiDB, chDurationsDB, chAck)

	go tlp.ParDeltaTempDB(killbroker, chDurationsDB, chDdbPar...)
	go tlp.BatchSend(nil, cfg, killbroker, chDdbPar[0])
	go tlp.WriteImages(nil, cfg, killbroker, chDdbPar[1])

	go func() {
		minutes := 0

		m, _ := time.ParseDuration(fmt.Sprintf("%vm", skip))
		for secs := 0; secs < NUMCYCLESTORUN; secs++ {
			chSec <- true
			mock.Add(m)
			lw.Debug("MOCK NOW ", mock.Now())

			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			minutes += skip
		}
		log.Println("Killing the test network.")
		killbroker.Publish(tlp.Ping{})
		wg.Done()
	}()
	wg.Wait()

}
