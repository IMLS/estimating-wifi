package tlp

import (
	"path/filepath"
	"reflect"
	"strings"

	"github.com/jmoiron/sqlx"
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

	tdb.AddTable("durations", map[string]string{
		"id":                 "INTEGER PRIMARY KEY AUTOINCREMENT",
		"pi_serial":          "TEXT",
		"fcfs_seq_id":        "TEXT",
		"device_tag":         "TEXT",
		"session_id":         "TEXT",
		"manufacturer_index": "INTEGER",
		"patron_index":       "INTEGER",
		"start":              "DATE",
		"end":                "DATE",
	})
	lw.Info("Created table wifi")
	return tdb
}

func extractWifiEvents(memdb *sqlx.DB) []analysis.WifiEvent {
	lw := logwrapper.NewLogger(nil)

	events := []analysis.WifiEvent{}
	err := memdb.Select(&events, `SELECT * FROM wifi`)
	if err != nil {
		lw.Info("error in extracting all wifi events.")
		lw.Fatal(err.Error())
	}

	lw.Length("events", events)
	return events
}

// key is "db"
func getFieldName(tag, key string, s interface{}) (fieldname string) {
	rt := reflect.TypeOf(s)
	if rt.Kind() != reflect.Struct {
		panic("bad type")
	}
	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		v := strings.Split(f.Tag.Get(key), ",")[0]
		if v == tag {
			return f.Name
		}
	}
	return ""
}

func storeSummary(cfg *config.Config, tdb *model.TempDB, c *analysis.Counter, durations map[int]*analysis.Duration) {
	lw := logwrapper.NewLogger(nil)

	fields := tdb.GetFields("durations")
	// For each duration we want to store
	for _, d := range durations {
		values := make(map[string]interface{})
		// For each field name in the DB
		for _, field := range fields {
			// Get the struct field name.
			structFieldName := getFieldName(field, "db", t)
			// Reflect on the duration
			r := reflect.ValueOf(d)
			v := reflect.Indirect(r).FieldByName(structFieldName)
			values[field] = v.String()
		}
		tdb.Insert("durations", values)
	}
}

func processDataFromDay(cfg *config.Config, tdb *model.TempDB) {
	lw := logwrapper.NewLogger(nil)

	events := []analysis.WifiEvent{}
	tdb.SelectAll("events", events)

	if len(events) > 0 {
		lw.Length("events", events)
		c, d := analysis.Summarize(cfg, events)
		//writeImages(cfg, events)
		//writeSummaryCSV(cfg, events)
		storeSummary(cfg, tdb, c, d)
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
