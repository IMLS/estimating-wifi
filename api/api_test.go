package api

import (
	"fmt"
	"testing"

	"gsa.gov/18f/session-counter/config"
)

// This should be much higher, like 2000
// But, it slows down practical testing... :/
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
			got := MacToMfg(&cfg, tc.mac)
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
	// cfg := config.Config{}
	// cfg.Manufacturers.Db = "/home/pi/git/imls/session-counter/manufacturer-db/manufacturers.sqlite"
	cfg := config.ReadConfig()

	for ndx := 0; ndx < dbIterations; ndx++ {
		t.Run(fmt.Sprintf("Thrash DB = %d", ndx), func(t *testing.T) {
			got := MacToMfg(cfg, "aa:bb:cc")
			if got != "unknown" {
				t.Fatalf("got %v; want %v", got, "unknown")
			} else {
				t.Logf("Success !")
			}

		})
	}
}

func Test_ReadAuth(t *testing.T) {
	a, e := config.ReadAuth()
	if e != nil {
		t.Fatal("failure in reading auth")
	}
	if a == nil {
		t.Fatal("auth is nil")
	}

}

func Test_GetToken(t *testing.T) {
	auth, err := config.ReadAuth()

	if err != nil {
		t.Log(err)
		t.Fatal("Failed to read token.")
	}

	if len(auth.Token) < 2 {
		t.Log(auth)
		t.Fatal("Failed to find token in auth struct.")
	}

}

func Test_StoreContent(t *testing.T) {
	cfg := config.ReadConfig()
	// Fill in the rest of the config.
	cfg.SessionId = config.CreateSessionId()
	cfg.Serial = config.GetSerial()

	auth, _ := config.ReadAuth()
	t.Log("auth", auth)
	StoreDeviceCount(cfg, auth, 42, "Next:0")
}

func Test_LogEvent(t *testing.T) {
	cfg := config.ReadConfig()
	// Fill in the rest of the config.
	cfg.SessionId = config.CreateSessionId()
	cfg.Serial = config.GetSerial()
	auth, _ := config.ReadAuth()
	t.Log("auth", auth)

	// Create a new logger

	el := NewEventLogger(cfg)
	el.Log("startup", map[string]string{"msg": "starting session-counter"})
	el.Log("empty", map[string]string{})
	el.Log("nil", nil)

}
