package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"testing"

	"gsa.gov/18f/session-counter/csp"
	"gsa.gov/18f/session-counter/model"
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
			{Mfg: "Eriksson", Id: 2}: 3,
			{Mfg: "unknown", Id: 3}:  0,
		},
	},
}

func assertEqual(t *testing.T, a interface{}, b interface{}, message string) {
	if a == b {
		return
	}
	if len(message) == 0 {
		message = fmt.Sprintf("%v != %v", a, b)
	}
	t.Fatal(message)
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

	f, err := os.OpenFile("testlogfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

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
				t.Log("send ", h)
				ch_macs <- h
			}
			for _, h := range e.loopMaps {
				ch_macs <- h
			}
			defer wg.Done()
		}()

		go rawToUids(ka, cfg, ch_macs, ch_uniq, ch_poison)

		wg.Add(1)
		go func() {
			count := len(e.loopMaps) + len(e.initMaps) - 1
			t.Log("receiving ", count)
			for i := 0; i < count; i++ {
				h := <-ch_uniq
				t.Log("receive ", h)
			}
			u := <-ch_uniq
			ch_poison <- true
			defer wg.Done()

			if e.passfail {
				s1 := fmt.Sprint(u)
				s2 := fmt.Sprint(e.resultMap)
				t.Log("s1: ", s1)
				t.Log("s2: ", s2)
				t.Log("s1 == s2: ", s1 == s2)

				assertEqual(t, fmt.Sprint(u), fmt.Sprint(e.resultMap), "maps not equal")
			} else {
				assertNotEqual(t, fmt.Sprint(u), fmt.Sprint(e.resultMap), "maps incorrectly equal")
			}
		}()

		wg.Wait()
	}
}
