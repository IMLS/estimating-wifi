package tlp

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
	"github.com/stretchr/testify/suite"
	"gsa.gov/18f/analysis"
	"gsa.gov/18f/config"
	"gsa.gov/18f/logwrapper"
	"gsa.gov/18f/tempdb"
)

const PASS = true
const FAIL = false

type TLPSuite struct {
	suite.Suite
	cfg  *config.Config
	mock *clock.Mock
	lw   *logwrapper.StandardLogger
}

// Make sure that VariableThatShouldStartAtFive is set to five
// before each test
func (suite *TLPSuite) SetupTest() {

	// Get a local config, so we have paths...
	_, filename, _, _ := runtime.Caller(0)
	fmt.Println(filename)
	path := filepath.Dir(filename)

	configPath := filepath.Join(path, "..", "test", "config.yaml")
	//suite.lw.Debug("path ", configPath)
	cfg, _ := config.NewConfigFromPath(configPath)

	if cfg == nil {
		suite.Fail("config is nil")
	}
	suite.cfg = cfg
	suite.lw = logwrapper.NewLogger(suite.cfg)
	suite.lw.SetLogLevel("DEBUG")

	mock := clock.NewMock()
	suite.mock = mock
	if mock == nil {
		suite.Fail("mock is nil")
	}
	mt, _ := time.Parse("2006-01-02T15:04", "1975-10-11T18:00")
	mock.Set(mt)
	suite.cfg.Clock = mock
	suite.lw.Debug("mock is now ", suite.cfg.Clock.Now())
	suite.cfg.RunMode = "test"
	suite.cfg.StorageMode = "sqlite"
	suite.cfg.Manufacturers.Db = filepath.Join(path, "..", "test", "manufacturers.sqlite")
	suite.cfg.Local.WebDirectory = filepath.Join(path, "..", "test", "www")
	os.Mkdir(suite.cfg.Local.WebDirectory, 0755)
	suite.cfg.NewSessionId()

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

func RunFakeWireshark(ka *Keepalive, cfg *config.Config, kb *KillBroker, in <-chan bool, out chan []string) {
	NUMMACS := 200
	NUMRANDOM := 10
	lw := logwrapper.NewLogger(nil)
	lw.Debug("RunFakeWireshark in the house.")

	ch_kill := kb.Subscribe()
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

		case <-ch_kill:
			log.Println("Exiting RunFakeWireshark")
			return
		}
	}
}

func PingAtBogoMidnight(ka *Keepalive, cfg *config.Config,
	rb *ResetBroker,
	kb *KillBroker,
	m *clock.Mock) {
	// counter := 0
	// ch_kill := kb.Subscribe()
	lw := logwrapper.NewLogger(nil)
	pinged := false
	for {
		if m.Now().Hour() == 0 && !pinged {
			pinged = true
			lw.Debug("IT IS BOGOMIDNIGHT.")
			rb.Publish(Ping{})
		}
		if m.Now().Hour() != 0 {
			pinged = false
		}
	}
}
func (suite *TLPSuite) TestManyTLPCycles() {

	// Create channels for process network
	ch_sec := make(chan bool)

	ch_nsec := make(chan bool)
	ch_macs := make(chan []string)
	ch_macs_counted := make(chan map[string]int)
	ch_data_for_report := make(chan []analysis.WifiEvent)
	ch_wifidb := make(chan *tempdb.TempDB)
	ch_durations_db := make(chan *tempdb.TempDB)
	ch_ack := make(chan Ping)
	ch_ddb_par := make([]chan *tempdb.TempDB, 2)
	for i := 0; i < 2; i++ {
		ch_ddb_par[i] = make(chan *tempdb.TempDB)
	}

	// See if we can wait and shut down the test...
	var wg sync.WaitGroup
	wg.Add(1)

	resetbroker := NewResetBroker()
	go resetbroker.Start()
	killbroker := NewKillBroker()
	go killbroker.Start()

	// Tock every two seconds.
	go TockEveryN(nil, killbroker, 2, ch_sec, ch_nsec)

	// Need a fake RunWireshark
	// go tlp.RunWireshark(nil, cfg, ch_nsec1, ch_macs, KC[1])
	go RunFakeWireshark(nil, suite.cfg, killbroker, ch_nsec, ch_macs)

	go AlgorithmTwo(nil, suite.cfg, resetbroker, killbroker, ch_macs, ch_macs_counted)
	go PrepEphemeralWifi(nil, suite.cfg, killbroker, ch_macs_counted, ch_data_for_report)
	// At midnight, flush internal structures and restart.
	//go tlp.PingAtMidnight(nil, cfg, chs_reset[0], KC[4])
	go PingAtBogoMidnight(nil, suite.cfg, resetbroker, killbroker, suite.mock)
	go CacheWifi(nil, suite.cfg, resetbroker, killbroker, ch_data_for_report, ch_wifidb, ch_ack)
	// Make sure we don't hang...
	go GenerateDurations(nil, suite.cfg, killbroker, ch_wifidb, ch_durations_db, ch_ack)

	go ParDeltaTempDB(killbroker, ch_durations_db, ch_ddb_par...)
	go BatchSend(nil, suite.cfg, killbroker, ch_ddb_par[0])
	go WriteImages(nil, suite.cfg, killbroker, ch_ddb_par[1])

	NUMCYCLESTORUN := 400

	go func() {
		minutes := 0
		skip := 20
		m, _ := time.ParseDuration(fmt.Sprintf("%vm", skip))
		for secs := 0; secs < NUMCYCLESTORUN; secs++ {
			ch_sec <- true
			if secs%2 == 0 {
				suite.mock.Add(m)
				suite.lw.Debug("MOCK NOW ", suite.mock.Now())
				// var m runtime.MemStats
				// runtime.ReadMemStats(&m)
				minutes += skip
				// memstats := fmt.Sprintf("Alloc[%vMB] Sys[%vMB], NumGC[%v]", bToMb(m.Alloc), bToMb(m.Sys), m.NumGC)
				// log.Println(days, "d", hours, "h", minutes%60, "m", memstats)
			}
		}
		log.Println("Killing the test network.")
		killbroker.Publish(Ping{})
		wg.Done()
	}()
	wg.Wait()
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

func (suite *TLPSuite) TestRawToUid(t *testing.T) {
	log.Println("TestRawToUid")
	cfg := new(config.Config)

	_, filename, _, _ := runtime.Caller(0)
	fmt.Println(filename)
	path := filepath.Dir(filename)
	cfg.Manufacturers.Db = filepath.Join(path, "test", "manufacturers.sqlite")

	ka := NewKeepalive(cfg)

	// var buf bytes.Buffer
	// log.SetOutput(&buf)
	// defer func() {
	// 	log.SetOutput(os.Stderr)
	// }()

	for testNdx, e := range tests {
		t.Logf("Test #%v: %v\n", testNdx, e.description)
		cfg.Monitoring.UniquenessWindow = e.uniqueness_window

		var wg sync.WaitGroup
		resetbroker := NewResetBroker()
		go resetbroker.Start()
		// The kill broker lets us poison the network.
		// var killbroker *tlp.Broker = nil
		killbroker := NewKillBroker()
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
		go AlgorithmTwo(ka, cfg, resetbroker, killbroker, ch_macs, ch_uniq)

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
			killbroker.Publish(Ping{})
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

func TestTLPSuite(t *testing.T) {
	suite.Run(t, new(TLPSuite))
}
