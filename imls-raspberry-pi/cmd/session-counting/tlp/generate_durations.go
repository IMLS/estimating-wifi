package tlp

import (
	"path/filepath"

	"gsa.gov/18f/analysis"
	"gsa.gov/18f/config"
	"gsa.gov/18f/logwrapper"
	"gsa.gov/18f/session-counter/model"
)

func createDurationsTable(cfg *config.Config) *model.TempDB {
	lw := logwrapper.NewLogger(nil)
	durationsDB := "durations.sqlite"
	fullpath := filepath.Join(cfg.Local.WebDirectory, durationsDB)
	tdb := model.NewSqliteDB(durationsDB, fullpath)
	lw.Info("Created durations db: %v", fullpath)
	// Add in the table.
	tdb.AddStructAsTable("durations", analysis.Duration{})
	lw.Info("Created durations table")
	return tdb
}

func storeSummary(cfg *config.Config, tdb *model.TempDB, c *analysis.Counter, durations map[int]*analysis.Duration) {
	for _, d := range durations {
		tdb.InsertStruct("durations", d)
	}
}

func processDataFromDay(cfg *config.Config, tdb *model.TempDB) {
	lw := logwrapper.NewLogger(nil)
	events := []analysis.WifiEvent{}
	tdb.SelectAll("events", events)
	if len(events) > 0 {
		lw.Length("events", events)
		c, d := analysis.Summarize(cfg, events)
		storeSummary(cfg, tdb, c, d)
		//writeImages(cfg, events)
		//writeSummaryCSV(cfg, events)
	} else {
		lw.Info("no `events` to summarize")
	}
}

func GenerateDurations(ka *Keepalive, cfg *config.Config, kb *Broker,
	ch_db <-chan *model.TempDB) {
	lw := logwrapper.NewLogger(nil)
	lw.Debug("starting GenerateDurations")
	var ping, pong chan interface{} = nil, nil
	var ch_kill chan interface{} = nil
	if kb != nil {
		ch_kill = kb.Subscribe()
	} else {
		ping, pong = ka.Subscribe("GenerateDurations", 30)
	}

	for {
		select {
		case <-ping:
			pong <- "GenerateDurations"
		case <-ch_kill:
			lw.Debug("exiting GenerateDurations")
			return

		case tdb := <-ch_db:
			// When we're passed the DB pointer, that means a reset has been triggered
			// up the chain. So, we need to process the events from the day.
			processDataFromDay(cfg, tdb)
		}
	}
}
