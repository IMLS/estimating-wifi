package tlp

import (
	"fmt"
	"path/filepath"
	"time"

	"gsa.gov/18f/config"
	"gsa.gov/18f/logwrapper"
	"gsa.gov/18f/session-counter/model"
)

func newTempDbInFS(cfg *config.Config) *model.TempDB {
	lw := logwrapper.NewLogger(nil)

	t := time.Now()
	todaysDB := fmt.Sprintf("%v-%02d-%02d-wifi.sqlite", t.Year(), int(t.Month()), int(t.Day()))
	path := filepath.Join(cfg.Local.WebDirectory, todaysDB)
	tdb := model.NewSqliteDB(todaysDB, path)
	lw.Info("Created temporary db: %v", todaysDB)
	// First, remove the table if it exists
	// If we reboot midday, this means we will start a fresh table.
	tdb.DropTable("wifi")
	// Add in the table.
	tdb.AddTable("wifi", map[string]string{
		"id":                 "INTEGER PRIMARY KEY AUTOINCREMENT",
		"fcfs_seq_id":        "TEXT",
		"device_tag":         "TEXT",
		"localtime":          "DATE",
		"session_id":         "TEXT",
		"manufacturer_index": "INTEGER",
		"patron_index":       "INTEGER"})
	lw.Info("Created table wifi")
	return tdb
}

func CacheWifi(ka *Keepalive, cfg *config.Config, rb *Broker, kb *Broker,
	ch_data <-chan []map[string]string, ch_db chan<- *model.TempDB) {
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

	event_ndx := 0
	tdb := newTempDbInFS(cfg)

	for {
		select {
		case <-ping:
			pong <- "StoreToSqlite"
		case <-ch_kill:
			tdb.Close()
			lw.Debug("exiting StoreToSqlite")
			return

		case <-ch_reset:
			// At reset, we pass the DB pointer down the stream
			// and let interesting things happen.
			lw.Info("received reset message")
			ch_db <- tdb
			// Once we come back, we should init a new DB.
			tdb = newTempDbInFS(cfg)

		case arr := <-ch_data:
			lw.Info("storing temporary data")
			for _, h := range arr {
				tdb.Insert("wifi", h)
			}
			event_ndx += 1
		}
	}
}
