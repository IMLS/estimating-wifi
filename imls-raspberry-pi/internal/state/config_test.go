package config

import (
	"log"
	"testing"
	"time"

	"github.com/benbjohnson/clock"
)

func Test_Config(t *testing.T) {
	cfg := NewConfig()
	//cfg.Validate()
	log.Println(cfg)
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
