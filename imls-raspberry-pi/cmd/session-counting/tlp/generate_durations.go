package tlp

import (
	"fmt"
	"path/filepath"

	"gsa.gov/18f/analysis"
	"gsa.gov/18f/config"
	"gsa.gov/18f/logwrapper"
	"gsa.gov/18f/session-counter/constants"
	"gsa.gov/18f/session-counter/model"
)

func GetDurationsDB(cfg *config.Config) *model.TempDB {
	lw := logwrapper.NewLogger(nil)
	fullpath := filepath.Join(cfg.Local.WebDirectory, constants.DURATIONSDB)
	tdb := model.NewSqliteDB(constants.DURATIONSDB, fullpath)
	tdb.AddStructAsTable("durations", analysis.Duration{})
	lw.Info("Created durations table in db [", fullpath, "]")
	return tdb
}

func storeSummary(cfg *config.Config, tdb *model.TempDB, c *analysis.Counter, durations map[int]analysis.Duration) {
	for _, d := range durations {
		tdb.InsertStruct("durations", d)
	}
}

func processDataFromDay(cfg *config.Config, table string, wifidb *model.TempDB) *model.TempDB {
	lw := logwrapper.NewLogger(nil)
	ddb := GetDurationsDB(cfg)
	events := []analysis.WifiEvent{}

	lw.Debug("selecting all events from wifi table")
	wifidb.Open()
	err := wifidb.Ptr.Select(&events, fmt.Sprintf("SELECT * FROM %v", table))
	wifidb.Close()

	if err != nil {
		lw.Info("error in extracting all events: %v", table)
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
	ch_db <-chan *model.TempDB,
	ch_durations_db chan<- *model.TempDB,
	ch_ack chan<- Ping) {
	lw := logwrapper.NewLogger(nil)
	lw.Debug("starting GenerateDurations")
	var ping, pong chan interface{} = nil, nil
	var ch_kill chan interface{} = nil
	if kb != nil {
		ch_kill = kb.Subscribe()
	} else {
		ping, pong = ka.Subscribe("GenerateDurations", 30)
	}

	// Queues for processing duration data
	sq := model.NewQueue(cfg, "sent")
	iq := model.NewQueue(cfg, "images")

	for {
		select {
		case <-ping:
			pong <- "GenerateDurations"
		case <-ch_kill:
			lw.Debug("exiting GenerateDurations")
			return

		case wifidb := <-ch_db:
			// When we're passed the DB pointer, that means a reset has been triggered
			// up the chain. So, we need to process the events from the day.
			durationsdb := processDataFromDay(cfg, constants.WIFIDB, wifidb)
			// Creates the table if it does not exist.
			//durationsdb.AddStructAsTable("batches", model.Batch{})
			yestersession := model.GetYesterdaySessionId()
			sq.Enqueue(yestersession)
			iq.Enqueue(yestersession)
			// When we're done processing everything, let CacheWifi know
			// that it is safe to continue.
			ch_ack <- Ping{}
			// Everything else is related to the duraitons DB, so that
			// can happen in parallel with other work.
			ch_durations_db <- durationsdb

		}
	}
}
