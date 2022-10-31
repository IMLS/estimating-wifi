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
	"gsa.gov/18f/internal/wifi-hardware-search/search"

	_ "net/http/pprof"
)

// var (
// 	cfgFile string
// 	mode    string
// )

func runEvery(crontab string, c *cron.Cron, fun func()) {
	id, err := c.AddFunc(crontab, fun)
	log.Debug().
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

	// gather wifi statistics every minute
	go runEvery("*/1 * * * *", c,
		func() {
			if config.IsDeveloperMode() {
				log.Debug().Msg("DEV MODE, RUNNING FAKE MOCK RUN")
				mock_hw.FakeWiresharkHelper(10, 200000)
			} else {
				log.Debug().Msg("RUNNING SIMPLESHARK")
				tlp.SimpleShark(
					search.SetMonitorMode,
					search.SearchForMatchingDevice,
					tlp.TSharkRunner)
			}
		})

	// send a heartbeat at the top of every hour
	go runEvery("0 * * * *", c, tlp.HeartBeat)

	// send the processed data at the end of each day
	go runEvery(config.GetResetCron(), c,
		func() {
			log.Info().
				Str("time", fmt.Sprintf("%v", state.GetClock().Now().In(time.Local))).
				Msg("RUNNING PROCESSDATA")
			// Copy ephemeral durations over to the durations table
			tlp.ProcessData(durationsdb, sq)
			// Try sending the data
			tlp.SimpleSend(durationsdb)
			// Increment the session counter
			state.IncrementSessionID()
			// Clear out the ephemeral data for the next day of monitoring
			state.ClearEphemeralDB()
			durationsdb.ClearDurationsDB()
		})

	// Start the cron jobs...
	c.Start()
}
