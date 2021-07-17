package tlp

import (
	"fmt"
	"path/filepath"

	"gsa.gov/18f/internal/config"
	"gsa.gov/18f/internal/logwrapper"
	"gsa.gov/18f/internal/state"
	"gsa.gov/18f/internal/structs"
)

func newTempDBInFS(cfg *config.Config) *state.TempDB {
	lw := logwrapper.NewLogger(nil)

	t := cfg.Clock.Now()
	todaysDB := fmt.Sprintf("%04d%02d%02d-%v.sqlite", t.Year(), int(t.Month()), int(t.Day()), state.WIFIDB)
	path := filepath.Join(cfg.Local.WebDirectory, todaysDB)
	tdb := state.NewSqliteDB(todaysDB, path)
	lw.Info("Created temporary db: [ ", todaysDB, " ]")
	lw.Info("Path to DB: ", cfg.Local.WebDirectory)
	// First, remove the table if it exists
	// If we reboot midday, this means we will start a fresh table.
	tdb.DropTable(state.WIFIDB)
	// Add in the table.
	tdb.AddStructAsTable(state.WIFIDB, structs.WifiEvent{})
	lw.Info("Created table ", state.WIFIDB)
	return tdb
}

func newTempDBInMemory(cfg *config.Config) *state.TempDB {
	lw := logwrapper.NewLogger(nil)
	todaysDB := state.WIFIDB
	path := fmt.Sprintf(`file:%v?mode=memory&cache=shared`, todaysDB)
	tdb := state.NewSqliteDB(todaysDB, path)
	lw.Info("Created memory db: [ ", todaysDB, " ]")
	tdb.DropTable(state.WIFIDB)
	tdb.AddStructAsTable(state.WIFIDB, structs.WifiEvent{})
	lw.Info("Created table ", state.WIFIDB)
	return tdb
}

func newTempDB(cfg *config.Config) *state.TempDB {
	lw := logwrapper.NewLogger(nil)
	lw.Debug("IsProductionMode is ", cfg.IsProductionMode())
	lw.Debug("IsDeveloperMode is ", cfg.IsDeveloperMode())
	lw.Debug("cfg.RunMode is ", cfg.RunMode)

	if cfg.IsProductionMode() {
		lw.Debug("using in-mem DB for wifi (prod)")
		return newTempDBInMemory(cfg)
	} else {
		lw.Debug("using in-filesystem DB for wifi (dev)")
		return newTempDBInFS(cfg)
	}
}

func CacheWifi(ka *Keepalive, cfg *config.Config, rb *ResetBroker, kb *KillBroker,
	chData <-chan []structs.WifiEvent,
	chDB chan<- *state.TempDB,
	chAck <-chan Ping) {
	lw := logwrapper.NewLogger(nil)
	lw.Debug("starting CacheWifi")
	var ping, pong chan interface{} = nil, nil
	var chKill chan interface{} = nil
	if kb != nil {
		chKill = kb.Subscribe()
	} else {
		ping, pong = ka.Subscribe("CacheWifi", 30)
	}
	chReset := rb.Subscribe()

	tdb := newTempDB(cfg)

	for {
		select {
		case <-ping:
			pong <- "CacheWifi"
		case <-chKill:
			// TDB is now opened/closed automatically on all transactions.
			// tdb.Close()
			lw.Debug("exiting CacheWifi")
			return

		case <-chReset:
			// At reset, we pass the DB pointer down the stream
			// and let interesting things happen.
			lw.Info("received reset message")
			chDB <- tdb
			// BAD! NOW FIXED! RACE HAZARD!
			// We continue immediately, meaning the DB gets flushed. We need to
			// wait until wifi processing is complete. That means GenerateDurations must
			// complete before we continue.
			<-chAck

			cfg.SetSessionId(state.GetNextSessionID(cfg))
			lw.Info("UPDATING SESSION ID TO: ", cfg.SessionId)
			tdb = newTempDB(cfg)

		case wifiarr := <-chData:
			lw.Info("storing temporary data")
			data := make([]interface{}, 0)
			for _, elem := range wifiarr {
				data = append(data, elem)
			}
			tdb.InsertManyStructs(state.WIFIDB, data)
		}
	}
}
