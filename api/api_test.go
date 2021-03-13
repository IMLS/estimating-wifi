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
	}

	for ndx, tc := range tests {
		t.Run(fmt.Sprintf("Get Manufactuerer = %d", ndx), func(t *testing.T) {
			got := Mac_to_mfg(cfg, tc.mac)
			if got != tc.mfgs {
				t.Fatalf("got %v; want %v", got, tc.mfgs)
			} else {
				t.Logf("Success !")
			}

		})
	}
}
