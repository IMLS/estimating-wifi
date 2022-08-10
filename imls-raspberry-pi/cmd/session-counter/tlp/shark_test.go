package tlp

import (
	"os"
	"testing"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/rs/zerolog/log"
	"gsa.gov/18f/cmd/session-counter/state"
	"gsa.gov/18f/internal/config"
	"gsa.gov/18f/internal/wifi-hardware-search/models"
)

// There's a lot of copypasta in these tests.

func setup() {
	temp, err := os.CreateTemp("", "shark-test.ini")
	if err != nil {
		log.Fatal().Err(err).Msg("could not create ini file")
	}

	config.SetConfigAtPath(temp.Name())
	config.SetRunMode("test")
	config.SetFCFSSeqID("ME0000-001")
	config.SetDeviceTag("testing")

	log.Debug().Int64("session id", state.GetCurrentSessionID()).Msg("setup")

	mock := clock.NewMock()
	mt, _ := time.Parse("2006-01-02T15:04", "1975-10-11T02:00")
	mock.Set(mt)
	state.SetClock(mock)
	state.ClearEphemeralDB()

	if state.GetClock() == nil {
		log.Fatal().Msg("clock should not be nil")
	}
	log.Debug().Time("now", state.GetClock().Now().In(time.Local)).Msg("mock")
}

func fakeMonitorFn(d *models.Device) {

}

func fakeSearchFn() (d *models.Device) {
	d = &models.Device{Exists: true, Logicalname: "fakewan0"}
	return d
}

func fakeShark2(dev string) []string {
	return []string{"DE:AD:BE:EF:00:00", "BE:EF:00:00:00:00"}
}

func fakeShark1(dev string) []string {
	return []string{"DE:AD:BE:EF:00:00"}
}

func TestOneHour(t *testing.T) {

	setup()

	startTime, _ := time.Parse(time.RFC3339, "1975-10-11T08:00:00-04:00")
	endTime, _ := time.Parse(time.RFC3339, "1975-10-11T09:00:00-04:00")

	mock := clock.NewMock()
	// Bump the clock forward
	mock.Set(startTime)
	state.SetClock(mock)
	// Run once at the initial time.
	SimpleShark(fakeMonitorFn, fakeSearchFn, fakeShark2)
	mock.Set(endTime)
	state.SetClock(mock)

	SimpleShark(fakeMonitorFn, fakeSearchFn, fakeShark2)

	macs := state.GetMACs()

	// We should now be able to check the DB.
	for _, testmac := range []string{"DE:AD:BE:EF:00:00", "BE:EF:00:00:00:00"} {
		p, ok := macs[testmac]
		if !ok {
			log.Fatal().Str("testmac", testmac).Msg("could not find test mac")
		}
		if p.Start != startTime.Unix() || p.End != endTime.Unix() {
			log.Fatal().
				Str("testmac", testmac).
				Int64("testmac start", p.Start).
				Int64("testmac end", p.End).
				Int64("stored start", startTime.Unix()).
				Int64("stored end", endTime.Unix()).
				Msg("values do not match")
		}
	}
}

func TestOneYear(t *testing.T) {

	setup()

	startTime, _ := time.Parse(time.RFC3339, "1975-10-11T08:00:00-04:00")
	endTime, _ := time.Parse(time.RFC3339, "1976-10-11T09:00:00-04:00")

	mock := clock.NewMock()
	// Bump the clock forward
	mock.Set(startTime)
	state.SetClock(mock)
	// Run once at the initial time.
	SimpleShark(fakeMonitorFn, fakeSearchFn, fakeShark2)
	mock.Set(endTime)
	state.SetClock(mock)

	SimpleShark(fakeMonitorFn, fakeSearchFn, fakeShark2)

	macs := state.GetMACs()

	// We should now be able to check the DB.
	for _, testmac := range []string{"DE:AD:BE:EF:00:00", "BE:EF:00:00:00:00"} {
		p, ok := macs[testmac]
		if !ok {
			log.Fatal().Str("testmac", testmac).Msg("could not find test mac")
		}
		if p.Start != startTime.Unix() || p.End != endTime.Unix() {
			log.Fatal().
				Str("testmac", testmac).
				Int64("testmac start", p.Start).
				Int64("testmac end", p.End).
				Int64("stored start", startTime.Unix()).
				Int64("stored end", endTime.Unix()).
				Msg("TestOneYear: values do not match")
		}
	}
}

func TestBumpOne(t *testing.T) {

	setup()

	startTime, _ := time.Parse(time.RFC3339, "1975-10-11T08:00:00-04:00")
	endTime, _ := time.Parse(time.RFC3339, "1975-10-11T09:00:00-04:00")

	mock := clock.NewMock()
	// Bump the clock forward
	mock.Set(startTime)
	state.SetClock(mock)
	// Run once at the initial time.
	SimpleShark(fakeMonitorFn, fakeSearchFn, fakeShark2)
	mock.Set(endTime)
	state.SetClock(mock)

	SimpleShark(fakeMonitorFn, fakeSearchFn, fakeShark1)

	macs := state.GetMACs()

	// We should now be able to check the DB.
	for _, testmac := range []string{"DE:AD:BE:EF:00:00"} {
		p, ok := macs[testmac]
		if !ok {
			log.Fatal().Str("testmac", testmac).Msg("could not find test mac")
		}
		if p.Start != startTime.Unix() || p.End != endTime.Unix() {
			log.Fatal().
				Str("testmac", testmac).
				Int64("testmac start", p.Start).
				Int64("testmac end", p.End).
				Int64("stored start", startTime.Unix()).
				Int64("stored end", endTime.Unix()).
				Msg("TestBumpOne: values do not match")
		}
	}

	for _, testmac := range []string{"BE:EF:00:00:00:00"} {
		p, ok := macs[testmac]
		if !ok {
			log.Fatal().Str("testmac", testmac).Msg("could not find test mac")
		}
		if p.Start == startTime.Unix() || p.End == endTime.Unix() {
			log.Fatal().
				Str("testmac", testmac).
				Int64("testmac start", p.Start).
				Int64("testmac end", p.End).
				Int64("stored start", startTime.Unix()).
				Int64("stored end", endTime.Unix()).
				Msg("TestBumpOne: values DO match")
		}
	}
}
