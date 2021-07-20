package tlp

import (
	"fmt"
	"path/filepath"

	"gsa.gov/18f/internal/analysis"
	"gsa.gov/18f/internal/interfaces"
	"gsa.gov/18f/internal/logwrapper"
	"gsa.gov/18f/internal/state"
	"gsa.gov/18f/internal/structs"
)

func GetDurationsDB() interfaces.Database {
	cfg := state.GetConfig()
	fullpath := filepath.Join(cfg.Paths.WWW.Root, state.DURATIONSDB)
	tdb := state.NewSqliteDB(fullpath)
	tdb.CreateTableFromStruct(structs.Duration{})
	return tdb
}

func storeSummary(tdb interfaces.Database,
	c *analysis.Counter,
	durations map[int]structs.Duration) {
	for _, d := range durations {
		tdb.GetTableFromStruct(structs.Duration{}).InsertStruct(d)
	}
}

func processDataFromDay(wifidb interfaces.Database) interfaces.Database {
	cfg := state.GetConfig()
	ddb := GetDurationsDB()
	cfg.Log().Debug("selecting all events from wifi table")
	events := structs.WifiEvent{}.SelectAll(wifidb)
	if len(events) > 0 {
		c, d := analysis.Summarize(events)
		storeSummary(ddb, c, d)
	} else {
		cfg.Log().Info("no `events` to summarize")
	}
	return ddb
}

func GenerateDurations(ka *Keepalive, kb *KillBroker,
	chDB <-chan interfaces.Database,
	chDurationsDB chan<- interfaces.Database,
	chAck chan<- Ping) {
	cfg := state.GetConfig()
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
	sq := state.NewQueue("sent")
	iq := state.NewQueue("images")

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
			durationsdb := processDataFromDay(wifidb)
			thissession := cfg.GetCurrentSessionId()
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
