package tlp

import (
	"fmt"
	"path/filepath"

	"gsa.gov/18f/internal/analysis"
	"gsa.gov/18f/internal/config"
	"gsa.gov/18f/internal/logwrapper"
	"gsa.gov/18f/internal/state"
	"gsa.gov/18f/internal/structs"
)

func GetDurationsDB(cfg *config.Config) *state.TempDB {
	lw := logwrapper.NewLogger(nil)
	fullpath := filepath.Join(cfg.Local.WebDirectory, state.DURATIONSDB)
	tdb := state.NewSqliteDB(state.DURATIONSDB, fullpath)
	tdb.AddStructAsTable("durations", structs.Duration{})
	lw.Info("Created durations table in db [", fullpath, "]")
	return tdb
}

func storeSummary(cfg *config.Config, tdb *state.TempDB,
	c *analysis.Counter,
	durations map[int]structs.Duration) {
	for _, d := range durations {
		tdb.InsertStruct("durations", d)
	}
}

func processDataFromDay(cfg *config.Config, table string, wifidb *state.TempDB) *state.TempDB {
	lw := logwrapper.NewLogger(nil)
	ddb := GetDurationsDB(cfg)
	events := []structs.WifiEvent{}

	lw.Debug("selecting all events from wifi table")
	wifidb.Open()
	err := wifidb.Ptr.Select(&events, fmt.Sprintf("SELECT * FROM %v", table))
	wifidb.Close()

	if err != nil {
		lw.Info("error in extracting all events: ", table)
		lw.Fatal(err.Error())
	}
	if len(events) > 0 {
		//lw.Length("events", events)
		c, d := analysis.Summarize(cfg, events)
		storeSummary(cfg, ddb, c, d)
		//writeImages(cfg, events)
		//writeSummaryCSV(cfg, events)
	} else {
		lw.Info("no `events` to summarize")
	}
	return ddb
}

func GenerateDurations(ka *Keepalive, cfg *config.Config, kb *KillBroker,
	chDB <-chan *state.TempDB,
	chDurationsDB chan<- *state.TempDB,
	chAck chan<- Ping) {
	lw := logwrapper.NewLogger(nil)
	lw.Debug("starting GenerateDurations")
	var ping, pong chan interface{} = nil, nil
	var chKill chan interface{} = nil
	if kb != nil {
		chKill = kb.Subscribe()
	} else {
		ping, pong = ka.Subscribe("GenerateDurations", 30)
	}

	// Queues for processing duration data
	sq := state.NewQueue(cfg, "sent")
	iq := state.NewQueue(cfg, "images")

	for {
		select {
		case <-ping:
			pong <- "GenerateDurations"
		case <-chKill:
			lw.Debug("exiting GenerateDurations")
			return

		case wifidb := <-chDB:
			// When we're passed the DB pointer, that means a reset has been triggered
			// up the chain. So, we need to process the events from the day.
			durationsdb := processDataFromDay(cfg, state.WIFIDB, wifidb)
			thissession := state.GetCurrentSessionId(cfg) //state.GetPreviousSessionId(cfg)
			lw.Debug("queueing current session [ ", thissession, " ] to images and send queue... ")
			if thissession >= 0 {
				sq.Enqueue(fmt.Sprint(thissession))
				iq.Enqueue(fmt.Sprint(thissession))
			}
			// When we're done processing everything, let CacheWifi know
			// that it is safe to continue.
			chAck <- Ping{}
			// Everything else is related to the duraitons DB, so that
			// can happen in parallel with other work.
			chDurationsDB <- durationsdb

		}
	}
}
