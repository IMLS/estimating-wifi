package api

import (
	"fmt"
	"os"
	"testing"

	"gsa.gov/18f/session-counter/config"
	"gsa.gov/18f/session-counter/constants"
)

const dbIterations = 10

func Test_get_manufactuerer(t *testing.T) {
	cfg := config.Config{}
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
			got := Mac_to_mfg(&cfg, tc.mac)
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
	cfg := config.Config{}
	cfg.Manufacturers.Db = "/home/pi/git/imls/session-counter/manufacturer-db/manufacturers.sqlite"

	for ndx := 0; ndx < dbIterations; ndx++ {
		t.Run(fmt.Sprintf("Thrash DB = %d", ndx), func(t *testing.T) {
			got := Mac_to_mfg(&cfg, "aa:bb:cc")
			if got != "unknown" {
				t.Fatalf("got %v; want %v", got, "unknown")
			} else {
				t.Logf("Success !")
			}

		})
	}
}

func Test_GetToken(t *testing.T) {
	cfg := config.ReadConfig()
	a, _ := config.ReadAuth()
	os.Setenv(constants.EnvUsername, a.User)
	os.Setenv(constants.EnvPassword, a.Token)

	directusServer := config.GetServer(cfg, "directus")
	auth, err := GetToken(directusServer)
	if err != nil {
		t.Log(err)
		t.Fatal("Failed to get token.")
	}

	if len(auth.Token) < 8 {
		t.Fatal("Failed to get auth token.")
	}

}

func Test_GetAuth(t *testing.T) {
	a, e := config.ReadAuth()
	if e != nil {
		t.Fatal("failure in reading auth")
	}
	if a == nil {
		t.Fatal("auth is nil")
	}

	if len(a.User) < 3 {
		t.Fatal("auth username too short")
	}
	if len(a.Token) < 8 {
		t.Fatal("auth token too short")
	}
}
