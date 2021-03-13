package api

import (
	"fmt"
	"testing"

	"gsa.gov/18f/session-counter/model"
)

func Test_get_manufactuerer(t *testing.T) {
	cfg := model.Config{}
	cfg.Manufacturers.Db = "/home/pi/git/imls/session-counter/manufacturer-db/manufacturers.sqlite"

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

	for _, tc := range tests {
		t.Run(fmt.Sprintf("Get Manufactuerer = %v", tc.mfgs), func(t *testing.T) {
			got := Mac_to_mfg(cfg, tc.mac)
			if got != tc.mfgs {
				t.Fatalf("got [ %v ] want [ %v ]", got, tc.mfgs)
			} else {
				t.Logf("Success !")
			}

		})
	}
}

// I'm hoping that if we're leaking DB connections that
// this loop will find it. When the DB isn't closed properly,
// this will fail around 1078ish connections.
func Test_thrash_db(t *testing.T) {
	cfg := model.Config{}
	cfg.Manufacturers.Db = "/home/pi/git/imls/session-counter/manufacturer-db/manufacturers.sqlite"

	for ndx := 0; ndx < 2000; ndx++ {
		t.Run(fmt.Sprintf("Thrash DB = %d", ndx), func(t *testing.T) {
			got := Mac_to_mfg(cfg, "aa:bb:cc")
			if got != "unknown" {
				t.Fatalf("got %v; want %v", got, "unknown")
			} else {
				t.Logf("Success !")
			}

		})
	}
}
