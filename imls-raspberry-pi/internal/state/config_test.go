package state

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/benbjohnson/clock"
)

func Test_Config(t *testing.T) {
	cfg := NewConfig()
	//cfg.Validate()
	if cfg.Databases.DurationsPath != "/www/imls/durations.sqlite" {
		t.Fatal()
	}
}

func TestReadConfig(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	fmt.Println(filename)
	path := filepath.Dir(filename)
	configPath := filepath.Join(path, "..", "..", "cmd", "session-counter", "test", "config.yaml")
	UnsafeNewConfigFromPath(configPath)
	cfg := GetConfig()
	expected := map[string]string{
		cfg.Databases.DurationsPath: "/opt/imls/www/durations.sqlite",
		cfg.Device.DeviceTag:        "MYDEVICETAG",
	}
	unexpected := false
	for k, v := range expected {
		log.Println("comparing", k, v)
		if k != v {
			unexpected = true
			log.Println(k, "not equal to", v)
		}
	}
	if unexpected {
		t.Fail()
	}
}

func TestMock(t *testing.T) {
	cfg := NewConfig()
	cfg.Clock = clock.NewMock()
	year := cfg.Clock.Now().UTC().Year()
	if year != 1970 {
		t.Log("year is", year)
		t.Log(cfg.Clock.Now())
		t.Fail()
	}
}

func TestSetMock(t *testing.T) {
	cfg := NewConfig()
	mock := clock.NewMock()
	cfg.Clock = mock
	if cfg.Clock.Now().UTC().Year() != 1970 {
		t.Fail()
	}

	d, e := time.ParseDuration("24h")
	if e != nil {
		t.Log("could not parse duration")
		t.Log(e.Error())
		t.Fail()
	}
	mock.Set(cfg.Clock.Now().UTC().Add(3 * 365 * d))
	year := cfg.Clock.Now().UTC().Year()
	log.Println(year)
	if year != 1972 {
		t.Log("year is", year)
		t.Fail()
	}
}
