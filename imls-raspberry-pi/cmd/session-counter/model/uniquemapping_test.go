package model

import (
	"log"
	"path/filepath"
	"runtime"
	"testing"

	"gsa.gov/18f/internal/state"
)

func mapEqual(m1 map[string]int, m2 map[string]int) bool {
	equal := true
	for k, v1 := range m1 {
		v2, found := m2[k]
		equal = equal && found
		equal = equal && (v1 == v2)
	}
	return equal
}

func copyMap(in map[string]int) map[string]int {
	new := make(map[string]int)
	for k, v := range in {
		new[k] = v
	}

	return new
}

type maps struct {
	want map[string]int
	got  map[string]int
}

/*
var tests = []struct {
		mac  string
		mfgs string
	}{
		{"f4:39:09", "HewlettP"},
		{"48:00:33", "Technico"},
		{"3c:37:86", "unknown"},
		{"dc:a6:32", "Raspberr"},
		{"b0:34:95", "Apple"},
		{"60:38:e0:bd:15", "BelkinIn"},
	}
*/

func TestAsUserMappings(t *testing.T) {
	//cfg := config.ReadConfig()
	state.NewConfig()
	cfg := state.GetConfig()
	_, filename, _, _ := runtime.Caller(0)
	path := filepath.Dir(filename)
	cfg.Databases.ManufacturersPath = filepath.Join(path, "..", "test", "manufacturers.sqlite")
	state.InitConfig()

	umdb := NewUMDB(cfg)
	m1 := umdb.AsUserMappings()

	umdb.UpdateMapping("f4:39:09")
	m2got := umdb.AsUserMappings()
	m2want := make(map[string]int)
	m2want["0:0"] = 0

	// Advance the time. The device we just saw
	// should now be associated with a 1.
	umdb.AdvanceTime()
	m3got := umdb.AsUserMappings()
	m3want := make(map[string]int)
	m3want["0:0"] = 1

	//Add a new device.
	umdb.UpdateMapping("dc:a6:32:aa")
	m4got := umdb.AsUserMappings()
	m4want := make(map[string]int)
	m4want["0:0"] = 1
	m4want["1:1"] = 0

	// Tick
	umdb.AdvanceTime()
	m5got := umdb.AsUserMappings()
	m5want := make(map[string]int)
	m5want["0:0"] = 2
	m5want["1:1"] = 1

	// Poke the RPi
	umdb.UpdateMapping("dc:a6:32:aa")
	m6got := umdb.AsUserMappings()
	m6want := make(map[string]int)
	m6want["0:0"] = 2
	m6want["1:1"] = 0

	// Tick
	umdb.AdvanceTime()
	// Add a new, unique RPi
	umdb.UpdateMapping("dc:a6:32:bb")
	m7got := umdb.AsUserMappings()
	m7want := make(map[string]int)
	m7want["0:0"] = 3
	m7want["1:1"] = 1
	m7want["1:2"] = 0

	// Tick
	// Poke the first device
	umdb.AdvanceTime()
	umdb.UpdateMapping("f4:39:09")
	m8got := umdb.AsUserMappings()
	m8want := make(map[string]int)
	m8want["0:0"] = 0
	m8want["1:1"] = 2
	m8want["1:2"] = 1

	tests := [...]*maps{
		{want: m1, got: make(map[string]int)},
		{want: m2want, got: m2got},
		{want: m3want, got: m3got},
		{want: m4want, got: m4got},
		{want: m5want, got: m5got},
		{want: m6want, got: m6got},
		{want: m7want, got: m7got},
		{want: m8want, got: m8got},
	}

	for ndx, test := range tests {
		eq := mapEqual(test.want, test.got)
		if eq {
			log.Println(test.want, "==", test.got)
		} else {
			log.Println("want", test.want)
			log.Println("got", test.got)
			t.Fatalf("test %v: maps not equal", ndx+1)
		}
	}

	// Wipe the DB and re-run the tests. They should "just pass."
	umdb.WipeDB()

	for ndx, test := range tests {
		eq := mapEqual(test.want, test.got)
		if eq {
			log.Println(test.want, "==", test.got)
		} else {
			log.Println("want", test.want)
			log.Println("got", test.got)
			t.Fatalf("test %v: maps not equal", ndx+1)
		}
	}

}
