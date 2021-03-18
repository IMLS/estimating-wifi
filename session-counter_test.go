package main

import (
	"fmt"
	"sync"
	"testing"

	"gsa.gov/18f/session-counter/csp"
	"gsa.gov/18f/session-counter/model"
	"gsa.gov/18f/session-counter/tlp"
)

const PASS = true
const FAIL = false

var tests = []struct {
	passfail             bool
	uniqueness_window    int
	disconnection_window int
	initMaps             []map[string]int
	loopMaps             []map[string]int
	resultMap            map[model.UserMapping]int
}{
	// One input hash.
	{PASS, 10, 5,
		[]map[string]int{
			{"de:ad:be:ef": 42},
		},
		[]map[string]int{
			{"de:ad:be:ef": 87},
		},
		map[model.UserMapping]int{
			{Mfg: "unknown", Id: 0}: 0,
		},
	},
	// Two input hashes
	{PASS, 10, 5,
		[]map[string]int{
			{"de:ad:be:ef": 42},
			{"be:ef:ca:fe": 137},
		},
		[]map[string]int{
			{"de:ad:be:ef": 87},
		},
		// Why zero and one?
		// Zero for deadbeef, because we send it in the loop.
		// One for beefcafe, because it was only sent once, and
		// one tick goes by.
		map[model.UserMapping]int{
			{Mfg: "unknown", Id: 0}: 0,
			{Mfg: "unknown", Id: 1}: 1,
		},
	},
	// Three hashes, three minutes
	{PASS, 10, 5,
		[]map[string]int{
			{"00:00:0F": 42},  // Next
			{"00:03:93": 137}, // Apple
			{"00:01:ec": 4},   // Eriksson
		},
		[]map[string]int{
			{"de:ad:be:ef": 87},
			{"de:ad:be:ef": 87},
			{"de:ad:be:ef": 87},
		},
		// Why zero and one?
		// Zero for deadbeef, because we send it in the loop.
		// One for beefcafe, because it was only sent once, and
		// one tick goes by.
		map[model.UserMapping]int{
			{Mfg: "Next", Id: 0}:     5,
			{Mfg: "Apple", Id: 1}:    4,
			{Mfg: "Ericsson", Id: 2}: 3,
			{Mfg: "unknown", Id: 3}:  0,
		},
	},

	// Next times out, because it is considered to
	// have "disconnected" after 5 minutes.
	{PASS, 10, 5,
		[]map[string]int{
			{"00:00:0F": 42},  // Next
			{"00:03:93": 137}, // Apple
			{"00:01:ec": 4},   // Eriksson
		},
		[]map[string]int{
			{"de:ad:be:ef": 87},
			{"de:ad:be:ef": 87},
			{"de:ad:be:ef": 87},
			{"de:ad:be:ef": 87},
		},
		// Why zero and one?
		// Zero for deadbeef, because we send it in the loop.
		// One for beefcafe, because it was only sent once, and
		// one tick goes by.
		map[model.UserMapping]int{
			{Mfg: "Apple", Id: 1}:    5,
			{Mfg: "Ericsson", Id: 2}: 4,
			{Mfg: "unknown", Id: 3}:  0,
		},
	},

	// Next times out, comes back. Still ID 0.
	// Apple is considered to have disconnected.
	{PASS, 10, 5,
		[]map[string]int{
			{"00:00:0f": 42},  // Next
			{"00:03:93": 137}, // Apple
			{"00:01:ec": 4},   // Eriksson
		},
		[]map[string]int{
			{"de:ad:be:ef": 87},
			{"de:ad:be:ef": 87},
			{"de:ad:be:ef": 87},
			{"de:ad:be:ef": 87},
			{"00:00:0f": 4},
		},
		// Why zero and one?
		// Zero for deadbeef, because we send it in the loop.
		// One for beefcafe, because it was only sent once, and
		// one tick goes by.
		map[model.UserMapping]int{
			{Mfg: "Next", Id: 0}:     0,
			{Mfg: "Ericsson", Id: 2}: 5,
			{Mfg: "unknown", Id: 3}:  1,
		},
	},
}

func assertEqual(t *testing.T, a interface{}, b interface{}, message string) {
	if a == b {
		return
	}
	t.Fatal(message, "\n\ta: ", a, "\n\tb: ", b)
}

func assertNotEqual(t *testing.T, a interface{}, b interface{}, message string) {
	if a != b {
		return
	}
	if len(message) == 0 {
		message = fmt.Sprintf("%v != %v", a, b)
	}
	t.Fatal(message)
}

func TestRawToUid(t *testing.T) {
	cfg := new(model.Config)
	cfg.Manufacturers.Db = "/etc/session-counter/manufacturers.sqlite"
	ka := csp.NewKeepalive()

	for testNdx, e := range tests {
		t.Logf("Test #%v\n", testNdx)
		cfg.Monitoring.UniquenessWindow = e.uniqueness_window
		cfg.Monitoring.DisconnectionWindow = e.disconnection_window
		var wg sync.WaitGroup

		ch_macs := make(chan map[string]int)
		ch_uniq := make(chan map[model.UserMapping]int)
		ch_poison := make(chan bool)

		wg.Add(1)
		go func() {
			for _, h := range e.initMaps {
				//t.Log("send ", h)
				ch_macs <- h
			}
			for _, h := range e.loopMaps {
				ch_macs <- h
			}
			defer wg.Done()
		}()

		go tlp.RawToUids(ka, cfg, ch_macs, ch_uniq, ch_poison)

		wg.Add(1)
		go func() {
			count := len(e.loopMaps) + len(e.initMaps) - 1
			//t.Log("receiving ", count)
			for i := 0; i < count; i++ {
				<-ch_uniq
				//t.Log("receive ", h)
			}
			u := <-ch_uniq
			ch_poison <- true
			defer wg.Done()

			s1 := fmt.Sprint(u)
			s2 := fmt.Sprint(e.resultMap)
			if e.passfail {
				assertEqual(t, s1, s2, "maps not equal")
			} else {
				assertNotEqual(t, s1, s2, "maps incorrectly equal")
			}
		}()

		wg.Wait()
	}
}
