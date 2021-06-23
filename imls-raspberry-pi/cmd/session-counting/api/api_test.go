package api

import (
	"encoding/json"
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"testing"

	"gsa.gov/18f/config"
	"gsa.gov/18f/http"
)

// This should be much higher, like 2000
// But, it slows down practical testing... :/
const dbIterations = 2000

func Test_get_manufactuerer(t *testing.T) {
	cfg := config.Config{}
	_, filename, _, _ := runtime.Caller(0)
	path := filepath.Dir(filename)
	cfg.Manufacturers.Db = filepath.Join(path, "..", "test", "manufacturers.sqlite")

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
	_, filename, _, _ := runtime.Caller(0)
	path := filepath.Dir(filename)

	cfg, _ := config.ReadConfig(filepath.Join(path, "..", "test", "config.yaml"))

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

func Test_RevalResponseUnmarshall(t *testing.T) {
	testString := `{
		"tables": [
		  {
			"headers": [
			  "event_id",
			  "device_uuid",
			  "lib_user",
			  "localtime",
			  "servertime",
			  "session_id",
			  "device_id"
			],
			"whole_table_errors": [],
			"rows": [
			  {
				"row_number": 2,
				"errors": [],
				"data": {
				  "event_id": "-1",
				  "device_uuid": "1000000089bbf88b",
				  "lib_user": "matthew.jadud@gsa.gov",
				  "localtime": "2021-04-02T10:46:53-04:00",
				  "servertime": "2021-04-02T10:46:53-04:00",
				  "session_id": "9475068c05fea81f",
				  "device_id": "unknown:6"
				}
			  }
			],
			"valid_row_count": 1,
			"invalid_row_count": 0
		  }
		],
		"valid": true
	  }`

	var rev http.RevalResponse
	err := json.Unmarshal([]byte(testString), &rev)
	if err != nil {
		log.Println("unmarshalling error:", err)

	} else {
		log.Println(rev)

	}
}
