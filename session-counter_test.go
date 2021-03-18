package main

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"

	"gsa.gov/18f/session-counter/csp"
	"gsa.gov/18f/session-counter/model"
	"gsa.gov/18f/session-counter/tlp"
)

const PASS = true
const FAIL = false

func macs(arr ...string) map[string]int {
	h := make(map[string]int)
	for _, s := range arr {
		h[s] = rand.Intn(1024)
	}
	return h
}

func hashes(arr ...string) []map[string]int {
	// Return a list of hashes, one hash for each string
	harr := make([]map[string]int, 0)
	for _, s := range arr {
		harr = append(harr, map[string]int{s: rand.Intn(1024)})
	}
	return harr
}

// IDs will be assigned in the raw_to_uids proc
// on a sorted list of MAC addrs. Therefore, we *know*
// in this test set that "next" will always be UID 0.
// (If "next" and "apple" are together)
var m = map[string]string{
	"next":     "00:00:0f", // ID 0
	"ericsson": "00:01:ec",
	"apple":    "00:03:93",
}

var tests = []struct {
	description          string
	passfail             bool
	uniqueness_window    int
	disconnection_window int
	initMap              map[string]int
	loopMaps             []map[string]int
	resultMap            map[model.UserMapping]int
}{
	// One input hash.
	{"one input mac, one loop mac",
		PASS, 10, 5,
		macs(m["next"]),
		hashes(m["next"]),
		map[model.UserMapping]int{
			{Mfg: "Next", Id: 0}: 0,
		},
	},
	// Two input hashes
	{"two input macs, one loop mac",
		PASS, 10, 5,
		macs(m["next"], m["apple"]),
		hashes(m["next"]),
		// Why zero and one?
		// Zero for deadbeef, because we send it in the loop.
		// One for beefcafe, because it was only sent once, and
		// one tick goes by.
		map[model.UserMapping]int{
			{Mfg: "Next", Id: 0}:  0,
			{Mfg: "Apple", Id: 1}: 1,
		},
	},
	// Three hashes, three minutes
	{"three input macs, three comms in the middle",
		PASS, 10, 5,
		// Next, Apple, Ericsson
		macs(m["next"], m["apple"], m["ericsson"]),
		hashes("de:ad:be:ef", "de:ad:be:ef", "de:ad:be:ef"),
		// IDs will be assigned by MAC address sort!
		map[model.UserMapping]int{
			{Mfg: "Next", Id: 0}:     3,
			{Mfg: "Apple", Id: 2}:    3,
			{Mfg: "Ericsson", Id: 1}: 3,
			{Mfg: "unknown", Id: 3}:  0,
		},
	},

	// Next times out, because it is considered to
	// have "disconnected" after 5 minutes.
	{"Next should disappear",
		PASS, 10, 5,
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
		map[model.UserMapping]int{
			{Mfg: "Apple", Id: 2}:    1,
			{Mfg: "Ericsson", Id: 1}: 0,
			{Mfg: "unknown", Id: 3}:  2,
		},
	},

	// Next times out, comes back. Still ID 0.
	// Apple is considered to have disconnected.
	{"Drop two",
		PASS, 10, 5,
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
		map[model.UserMapping]int{
			{Mfg: "Next", Id: 0}:    0,
			{Mfg: "unknown", Id: 3}: 1,
		},
	},
}

func assertEqual(t *testing.T, a interface{}, b interface{}, message string) {
	if a == b {
		return
	}
	t.Fatal(message, "\n\texpected: ", a, "\n\treceived: ", b)
}

func assertNotEqual(t *testing.T, a interface{}, b interface{}, message string) {
	if a != b {
		return
	}
	t.Fatal(message, "\n\texpected: ", a, "\n\treceived: ", b)
}

func TestRawToUid(t *testing.T) {
	cfg := new(model.Config)
	cfg.Manufacturers.Db = "/etc/session-counter/manufacturers.sqlite"
	ka := csp.NewKeepalive()

	// var buf bytes.Buffer
	// log.SetOutput(&buf)
	// defer func() {
	// 	log.SetOutput(os.Stderr)
	// }()

	for testNdx, e := range tests {
		t.Logf("Test #%v: %v\n", testNdx, e.description)
		cfg.Monitoring.UniquenessWindow = e.uniqueness_window
		cfg.Monitoring.DisconnectionWindow = e.disconnection_window
		var wg sync.WaitGroup

		ch_macs := make(chan map[string]int)
		ch_uniq := make(chan map[model.UserMapping]int)
		ch_poison := make(chan bool)
		var u map[model.UserMapping]int = nil

		wg.Add(1)
		go func() {
			ch_macs <- e.initMap
			for _, h := range e.loopMaps {
				ch_macs <- h
			}
			defer wg.Done()
		}()

		go tlp.RawToUids(ka, cfg, ch_macs, ch_uniq, ch_poison)

		wg.Add(1)
		go func() {
			// The init map
			<-ch_uniq
			count := len(e.loopMaps) - 1
			//t.Log("receiving ", count)
			for i := 0; i < count; i++ {
				<-ch_uniq
				//t.Log("receive ", h)
			}
			u = <-ch_uniq
			ch_poison <- true
			defer wg.Done()
		}()

		wg.Wait()
		// t.Log(buf.String())

		expected := fmt.Sprint(e.resultMap)
		received := fmt.Sprint(u)

		if e.passfail {
			assertEqual(t, expected, received, "not equal")
		} else {
			assertNotEqual(t, expected, received, "incorrectly equal")
		}
	} // end for over tests
}
