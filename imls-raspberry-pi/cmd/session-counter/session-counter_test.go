package main

import (
	"fmt"
	"log"
	"math/rand"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/benbjohnson/clock"

	"gsa.gov/18f/cmd/session-counter/tlp"
	"gsa.gov/18f/internal/state"
	"gsa.gov/18f/internal/wifi-hardware-search/models"
)

var NUMMACS int
var NUMFOUNDPERMINUTE int
var consistentMACs = make([]string, 0)

func generateFakeMac() string {
	var letterRunes = []rune("ABCDEF0123456789")
	b := make([]rune, 17)
	colons := [...]int{2, 5, 8, 11, 14}
	for i := 0; i < 17; i++ {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]

		for v := range colons {
			if i == colons[v] {
				b[i] = rune(':')
			}
		}
	}
	return string(b)
}

func runFakeWireshark(device string) []string {

	thisTime := rand.Intn(NUMFOUNDPERMINUTE)
	send := make([]string, thisTime)
	// cfg := state.GetConfig()
	//cfg.Log().Info("sending ", len(send), " this minute")
	for i := 0; i < thisTime; i++ {
		send[i] = consistentMACs[rand.Intn(len(consistentMACs))]
	}
	return send
}

func isItMidnight(now time.Time) bool {
	return (now.Hour() == 0 &&
		now.Minute() == 0 &&
		now.Second() == 0)

}

func MockRun(rundays int, nummacs int, numfoundperminute int) {
	cfg := state.GetConfig()
	// The MAC database MUST be ephemeral. Put it in RAM.

	sq := state.NewQueue("sent")
	iq := state.NewQueue("images")
	durationsdb := cfg.GetDurationsDatabase()

	// Create a pool of NUMMACS devices to draw from.
	// We will send NUMFOUNDPERMINUTE each minute
	NUMMACS = nummacs
	NUMFOUNDPERMINUTE = numfoundperminute
	consistentMACs = make([]string, NUMMACS)
	for i := 0; i < NUMMACS; i++ {
		consistentMACs[i] = generateFakeMac()
	}

	for days := 0; days < rundays; days++ {
		// Pretend we run once per minute for 24 hours
		for minutes := 0; minutes < 60*24; minutes++ {
			tlp.SimpleShark(
				// search.SetMonitorMode,
				func(d *models.Device) {},
				// search.SearchForMatchingDevice,
				func() *models.Device { return &models.Device{Exists: true, Logicalname: "fakewan0"} },
				// tlp.TSharkRunner
				runFakeWireshark)
			// Add one minute to the fake clock
			state.GetClock().(*clock.Mock).Add(1 * time.Minute)

			if isItMidnight(state.GetClock().Now()) {
				// Then we run the processing at midnight (once per 24 hours)
				cfg.Log().Info("RUNNING PROCESSDATA at " + fmt.Sprint(state.GetClock().Now()))
				// Copy ephemeral durations over to the durations table
				tlp.ProcessData(durationsdb, sq, iq)
				// Draw images of the data
				tlp.WriteImages(durationsdb)
				// Try sending the data
				tlp.SimpleSend(durationsdb)
				// Increment the session counter
				cfg.IncrementSessionID()
				// Clear out the ephemeral data for the next day of monitoring
				state.ClearEphemeralDB()
			}
		}

	}
}

func TestAllUp(t *testing.T) {
	// DO NOT USE LOGGING YET
	_, filename, _, _ := runtime.Caller(0)
	fmt.Println(filename)
	path := filepath.Dir(filename)
	state.SetConfigAtPath(filepath.Join(path, "test", "config.sqlite"))
	cfg := state.GetConfig()
	cfg.SetStorageMode("local")
	cfg.SetRootPath(filepath.Join(path, "test", "www"))
	cfg.SetImagesPath(filepath.Join(path, "test", "www", "images"))

	cfg.Log().SetLogLevel("DEBUG")
	cfg.Log().Info("initial session id: ", cfg.GetCurrentSessionID())

	// Fake the clock
	mock := clock.NewMock()
	mt, _ := time.Parse("2006-01-02T15:04", "1975-10-11T00:00")
	mock.Set(mt)
	state.SetClock(mock)

	MockRun(1, 200000, 10)
	log.Println("WAITING")
	time.Sleep(5 * time.Second)

	cfg.IncrementSessionID()
	cfg.Log().Info("next session id: ", cfg.GetCurrentSessionID())

	// Fake the clock
	mt, _ = time.Parse("2006-01-02T15:04", "1976-11-12T00:00")
	mock.Set(mt)
	state.SetClock(mock)

	MockRun(5, 200000, 10)

	mt, _ = time.Parse("2006-01-02T15:04", "1978-01-01T00:00")
	mock.Set(mt)
	state.SetClock(mock)

	MockRun(90, 2000000, 10)
}
