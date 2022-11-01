package session_counter

import (
	"fmt"
	"time"

	cron "github.com/robfig/cron/v3"
	"github.com/rs/zerolog/log"
	"gsa.gov/18f/internal/config"
	"gsa.gov/18f/internal/session-counter-helper/mock_hw"
	"gsa.gov/18f/internal/session-counter-helper/state"
	"gsa.gov/18f/internal/session-counter-helper/tlp"
	"gsa.gov/18f/internal/wifi-hardware-search/models"
	"gsa.gov/18f/internal/wifi-hardware-search/search"

	_ "net/http/pprof"
)

// var (
// 	cfgFile string
// 	mode    string
// )

func runEvery(process string, crontab string, c *cron.Cron, fun func()) {
	id, err := c.AddFunc(crontab, fun)
	log.Debug().
		Str("process", process).
		Str("crontab", crontab).
		Str("id", fmt.Sprintf("%v", id)).
		Msg("runEvery")
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("cron: could not set up crontab entry")
	}
}

func Run2() {
	sq := state.NewQueue[int64]("sent")
	durationsdb := state.NewDurationsDB()
	c := cron.New()
	var sem = make(chan int, 1)

	if config.IsDeveloperMode() {
		log.Debug().Msg("DEV MODE, RUNNING FAKESHARK (IT'S A DOLPHIN)")
		mock_hw.FakeWiresharkSetup()
	}

	// gather wifi statistics every minute
	go runEvery("Data collection loop", config.GetDataCollectionCron(), c,
		func() {
			// We have a race on the ephemeral DB.
			// If we're in the data collection loop, we shouldn't let the data send and clear the DB.
			sem <- 1
			if config.IsDeveloperMode() {
				tlp.SimpleShark(
					// search.SetMonitorMode,
					func(d *models.Device) {},
					// search.SearchForMatchingDevice,
					func() *models.Device { return &models.Device{Exists: true, Logicalname: "fakewan0"} },
					// tlp.TSharkRunner
					mock_hw.RunFakeWireshark)
			} else {
				log.Debug().Msg("RUNNING SIMPLESHARK")
				tlp.SimpleShark(
					search.SetMonitorMode,
					search.SearchForMatchingDevice,
					tlp.TSharkRunner)
			}
			// Release
			<-sem
		})

	// send a heartbeat at the top of every hour
	go runEvery("Heartbeat process", config.GetHeartbeatCron(), c, tlp.HeartBeat)

	// send the processed data at the end of each day
	go runEvery("End of day reset", config.GetResetCron(), c,
		func() {
			// We're in a race with the data collection loop.
			// Before sending, grab the semaphore.
			sem <- 1
			log.Info().
				Str("time", fmt.Sprintf("%v", state.GetClock().Now().In(time.Local))).
				Msg("RUNNING PROCESSDATA")
			// Copy ephemeral durations over to the durations table
			tlp.ProcessData(durationsdb, sq)
			// Try sending the data
			tlp.SimpleSend(durationsdb, sq)
			// Increment the session counter
			state.IncrementSessionID()
			// Clear out the ephemeral data for the next day of monitoring
			state.ClearEphemeralDB()
			durationsdb.ClearDurationsDB()
			// After clearing the DB, release the semaphore.
			<-sem
		})

	// Start the cron jobs...
	c.Start()
}
