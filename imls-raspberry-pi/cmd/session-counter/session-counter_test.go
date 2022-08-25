package main

import (
	"fmt"
	"net/http"
	"runtime"
	"testing"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/rs/zerolog/log"
	"gsa.gov/18f/cmd/session-counter/state"
	"gsa.gov/18f/cmd/session-counter/tlp"
)

var NUMMACS int
var NUMFOUNDPERMINUTE int

func isItMidnight(now time.Time) bool {
	return (now.Hour() == 0 &&
		now.Minute() == 0 &&
		now.Second() == 0)

}

func MockRun(rundays int, nummacs int, numfoundperminute int) *state.DurationsDB {
	// The MAC database MUST be ephemeral. Put it in RAM.

	sq := state.NewQueue[int64]("sent")
	durationsdb := state.NewDurationsDB()

	for days := 0; days < rundays; days++ {
		// Pretend we run once per minute for 24 hours
		for minutes := 0; minutes < 60*24; minutes++ {
			fakeWiresharkHelper(NUMFOUNDPERMINUTE, nummacs)
			// Add one minute to the fake clock
			state.GetClock().(*clock.Mock).Add(1 * time.Minute)

			if isItMidnight(state.GetClock().Now().In(time.Local)) {
				// Then we run the processing at midnight (once per 24 hours)
				log.Info().
					Str("time", fmt.Sprint(state.GetClock().Now().In(time.Local))).
					Msg("RUNNING PROCESSDATA")
				// Copy ephemeral durations over to the durations table
				tlp.ProcessData(durationsdb, sq)
				// Draw images of the data
				// tlp.WriteImages(durationsdb)
				// Try sending the data
				tlp.SimpleSend(durationsdb)
				// Increment the session counter
				state.IncrementSessionID()
				// Clear out the ephemeral data for the next day of monitoring
				state.ClearEphemeralDB()
			}
		}

	}
	return durationsdb
}

func TestAllUp(t *testing.T) {

	// Provides profiling.
	// curl -sK -v http://localhost:8080/debug/pprof/heap > heap.out
	// go tool pprof heap.out
	go http.ListenAndServe("localhost:8080", nil)

	// DO NOT USE LOGGING YET
	_, filename, _, _ := runtime.Caller(0)
	fmt.Println(filename)

	log.Info().
		Int64("session id", state.GetCurrentSessionID()).
		Msg("initial session")

	// Fake the clock
	mock := clock.NewMock()
	mt, _ := time.Parse("2006-01-02T15:04", "1975-10-11T00:01")
	mock.Set(mt)
	state.SetClock(mock)

	db := MockRun(1, 200000, 10)
	db.ClearDurationsDB()

	log.Info().Msg("WAITING")
	time.Sleep(5 * time.Second)

	state.IncrementSessionID()

	log.Info().
		Int64("session id", state.GetCurrentSessionID()).
		Msg("next session")

	// Fake the clock
	mt, _ = time.Parse("2006-01-02T15:04", "1976-11-12T00:01")
	mock.Set(mt)
	state.SetClock(mock)

	db = MockRun(5, 200000, 10)
	db.ClearDurationsDB()

	mt, _ = time.Parse("2006-01-02T15:04", "1978-01-01T00:01")
	mock.Set(mt)
	state.SetClock(mock)

	db = MockRun(90, 2000000, 10)
	db.ClearDurationsDB()

}
