package tlp

import (
	"database/sql"
	"log"
	"os"

	"gsa.gov/18f/analysis"
	"gsa.gov/18f/config"
)

//func store(service string, cfg *config.Config, session_id int, h map[string]int) error {

func getSummaryDB(cfg *config.Config) *sql.DB {
	if _, err := os.Stat(cfg.Local.SummaryDB); os.IsNotExist(err) {
		file, err := os.Create(cfg.Local.SummaryDB)
		if err != nil {
			log.Println("sqlite: could not create sqlite summary db file.")
			log.Fatal(err.Error())
		}
		file.Close()
	}

	db, err := sql.Open("sqlite3", cfg.Local.SummaryDB)
	if err != nil {
		log.Fatal("sqlite: could not open summary db.")
	}

	// Create tables if it doesn't exist
	createTableStatement := `
	CREATE TABLE IF NOT EXISTS summary (
		id text PRIMARY KEY,
		pi_serial character text,
		fcfs_seq_id character text,
		device_tag character text,
		session_id text,
		minimum_minutes integer,
		maximum_minutes integer,
		patron_count integer,
		patron_minutes integer,
		device_count integer,
		device_minutes integer,
		transient_count integer,
		transient_minutes integer
	);`

	_, err = db.Exec(createTableStatement)
	if err != nil {
		log.Fatal("sqlite: could not create table in db.")
	}

	return db
}

func newInMemoryDB() *sql.DB {
	const DB_STRING = ":memory:"

	db, err := sql.Open("sqlite3", DB_STRING)
	if err != nil {
		log.Fatal("sqlite: Could not create in-memory DB.")
	}
	// Create tables.
	createTableStatement := `
	DROP TABLE IF EXISTS wifi;
	CREATE TABLE wifi (
		id text PRIMARY KEY,
		event_id integer,
		fcfs_seq_id character text,
		device_tag character text,
		"localtime" date,
		session_id text,
		manufacturer_index integer,
		patron_index integer
	);`

	_, err = db.Exec(createTableStatement)
	if err != nil {
		log.Fatal("sqlite: could not create table in db.")
	}

	return db
}

func extractWifiEvents(memdb *sql.DB) []analysis.WifiEvent {
	events := make([]analysis.WifiEvent, 0)
	q := `SELECT * FROM wifi`
	rows, err := memdb.Query(q)
	if err != nil {
		log.Println("sqlite: error in sqlite query.")
		log.Fatal(err.Error())
	}

	for rows.Next() {
		e := analysis.WifiEvent{}
		err = rows.Scan(&e.ID, &e.EventId, &e.FCFSSeqId, &e.DeviceTag, &e.Localtime, &e.SessionId, &e.ManufacturerIndex, &e.PatronIndex)
		events = append(events, e)
	}

	return events
}

func storeSummary(db *sql.DB, c *analysis.Counter) {

}

func processDataFromDay(cfg *config.Config, memdb *sql.DB) {
	summarydb := getSummaryDB(cfg)
	events := extractWifiEvents(memdb)
	c := analysis.Summarize(events)
	storeSummary(summarydb, c)
}

// FIXME
// On reset, we need to process and clear the sqlite tables. This should ping once daily.
func StoreToSqlite(ka *Keepalive, cfg *config.Config, ch_data <-chan []map[string]string, ch_reset <-chan Ping) {
	log.Println("Starting StoreToSqlite")
	ping, pong := ka.Subscribe("StoreToSqlite", 30)

	// If we aren't logging events...
	event_ndx := 0
	// We'll use an in-memory DB for the recording of data throughout the day.
	db := newInMemoryDB()

	for {
		select {
		case <-ping:
			pong <- "StoreToSqlite"
		// This is the [ uid -> ticks ] map (uid looks like "Next:0")
		case <-ch_reset:
			// Process the data from the day.
			processDataFromDay(cfg, db)
			// Close the existing DB.
			db.Close()
			// Open a new one.
			db = newInMemoryDB()

		case h := <-ch_data:
			log.Println(h)
			event_ndx += 1
		}
	}
}
