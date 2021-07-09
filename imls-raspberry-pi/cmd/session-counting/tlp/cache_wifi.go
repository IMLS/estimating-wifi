package tlp

import (
	"fmt"
	"path/filepath"
	"time"

	"gsa.gov/18f/analysis"
	"gsa.gov/18f/config"
	"gsa.gov/18f/logwrapper"
	"gsa.gov/18f/session-counter/constants"
	"gsa.gov/18f/session-counter/model"
)

func newTempDbInFS(cfg *config.Config) *model.TempDB {
	lw := logwrapper.NewLogger(nil)

	t := time.Now()
	todaysDB := fmt.Sprintf("%04d%02d%02d-%v.sqlite", t.Year(), int(t.Month()), int(t.Day()), constants.WIFIDB)
	path := filepath.Join(cfg.Local.WebDirectory, todaysDB)
	tdb := model.NewSqliteDB(todaysDB, path)
	lw.Info("Created temporary db: [ ", todaysDB, " ]")
	// First, remove the table if it exists
	// If we reboot midday, this means we will start a fresh table.
	tdb.DropTable(constants.WIFIDB)
	// Add in the table.
	tdb.AddStructAsTable(constants.WIFIDB, analysis.WifiEvent{})
	lw.Info("Created table ", constants.WIFIDB)
	return tdb
}

func newTempDbInMemory(cfg *config.Config) *model.TempDB {
	lw := logwrapper.NewLogger(nil)
	todaysDB := constants.WIFIDB
	path := fmt.Sprintf(`file:%v?mode=memory&cache=shared`, todaysDB)
	tdb := model.NewSqliteDB(todaysDB, path)
	lw.Info("Created memory db: [ ", todaysDB, " ]")
	tdb.DropTable(constants.WIFIDB)
	tdb.AddStructAsTable(constants.WIFIDB, analysis.WifiEvent{})
	lw.Info("Created table ", constants.WIFIDB)
	return tdb
}

func newTempDb(cfg *config.Config) *model.TempDB {
	lw := logwrapper.NewLogger(nil)
	if cfg.IsProductionMode() {
		lw.Debug("using in-mem DB for wifi (prod)")
		return newTempDbInMemory(cfg)
	} else {
		lw.Debug("using in-filesystem DB for wifi (dev)")
		return newTempDbInFS(cfg)
	}
}

func CacheWifi(ka *Keepalive, cfg *config.Config, rb *ResetBroker, kb *KillBroker,
	ch_data <-chan []analysis.WifiEvent,
	ch_db chan<- *model.TempDB,
	ch_ack <-chan Ping) {
	lw := logwrapper.NewLogger(nil)
	lw.Debug("starting CacheWifi")
	var ping, pong chan interface{} = nil, nil
	var ch_kill chan interface{} = nil
	if kb != nil {
		ch_kill = kb.Subscribe()
	} else {
		ping, pong = ka.Subscribe("CacheWifi", 30)
	}
	ch_reset := rb.Subscribe()

	tdb := newTempDb(cfg)

	for {
		select {
		case <-ping:
			pong <- "CacheWifi"
		case <-ch_kill:
			// TDB is now opened/closed automatically on all transactions.
			// tdb.Close()
			lw.Debug("exiting CacheWifi")
			return

		case <-ch_reset:
			// At reset, we pass the DB pointer down the stream
			// and let interesting things happen.
			lw.Info("received reset message")
			ch_db <- tdb
			// BAD! NOW FIXED! RACE HAZARD!
			// We continue immediately, meaning the DB gets flushed. We need to
			// wait until wifi processing is complete. That means GenerateDurations must
			// complete before we continue.
			<-ch_ack
			tdb = newTempDb(cfg)

		case wifiarr := <-ch_data:
			lw.Info("storing temporary data")
			data := make([]interface{}, 0)
			for _, elem := range wifiarr {
				data = append(data, elem)
			}
			tdb.InsertManyStructs(constants.WIFIDB, data)
		}
	}
}
