package tlp

import (
	"database/sql"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
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
		id integer PRIMARY KEY AUTOINCREMENT,
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

func newInMemoryDB() *sqlx.DB {
	const DB_STRING = ":memory:"

	db, err := sqlx.Open("sqlite3", DB_STRING)
	if err != nil {
		log.Fatal("sqlite: Could not create in-memory db.")
	}

	clearInMemoryDB(db)
	return db
}

func clearInMemoryDB(db *sqlx.DB) {
	// Create tables.
	createTableStatement := `
	DROP TABLE IF EXISTS wifi;
	CREATE TABLE wifi (
		id integer PRIMARY KEY AUTOINCREMENT,
		event_id integer,
		fcfs_seq_id character text,
		device_tag character text,
		"localtime" date,
		session_id text,
		manufacturer_index integer,
		patron_index integer
	);`

	_, err := db.Exec(createTableStatement)
	if err != nil {
		log.Fatal("sqlite: could not create table in db.")
	}
}

func extractWifiEvents(memdb *sqlx.DB) []analysis.WifiEvent {
	// events := make([]analysis.WifiEvent, 0)
	events := []analysis.WifiEvent{}
	err := memdb.Select(&events, `SELECT * FROM wifi`)
	if err != nil {
		log.Println("sqlite: error in sqlite query.")
		log.Fatal(err.Error())
	}
	log.Println("events", events)
	// for rows.Next() {
	// 	e := analysis.WifiEvent{}
	// 	err = rows.Scan(&e.ID, &e.EventId, &e.FCFSSeqId, &e.DeviceTag, &e.Localtime, &e.SessionId, &e.ManufacturerIndex, &e.PatronIndex)
	// 	events = append(events, e)
	// }
	// rows.Close()

	return events
}

func storeSummary(cfg *config.Config, c *analysis.Counter) {
	log.Println("sqlite: getting summary db")
	summarydb := getSummaryDB(cfg)
	insertS, err := summarydb.Prepare(`INSERT INTO summary (pi_serial, fcfs_seq_id, device_tag, session_id, minimum_minutes, maximum_minutes, patron_count, patron_minutes, device_count, device_minutes, transient_count, transient_minutes) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)`)
	if err != nil {
		log.Println("sqlite: could not prepare insert statement.")
		log.Fatal(err.Error())
	}
	tok, _ := config.ReadAuth()
	res, err := insertS.Exec(config.GetSerial(), tok.FCFSId, tok.DeviceTag, cfg.SessionId, cfg.Monitoring.MinimumMinutes, cfg.Monitoring.MaximumMinutes, c.Patrons, c.PatronMinutes, c.Devices, c.DeviceMinutes, c.Transients, c.TransientMinutes)
	if err != nil {
		log.Println("sqlite: could not insert into summary db")
		log.Println(res)
		log.Fatal(err.Error())
	}
	summarydb.Close()
}

func processDataFromDay(cfg *config.Config, memdb *sqlx.DB) {
	log.Println("sqlite: extracting wifi events")
	events := extractWifiEvents(memdb)
	log.Println(len(events), "events found")
	log.Println(events)
	if len(events) > 0 {
		log.Println("sqlite: counting")
		c := analysis.Summarize(cfg, events)
		log.Println("sqlite:", c)
		storeSummary(cfg, c)
	} else {
		log.Println("sqlite: no events to summarize")
	}
}

// FIXME
// On reset, we need to process and clear the sqlite tables. This should ping once daily.
func StoreToSqlite(ka *Keepalive, cfg *config.Config, ch_data <-chan []map[string]string, ch_reset <-chan Ping, ch_kill <-chan Ping) {
	log.Println("Starting StoreToSqlite")

	var ping, pong chan interface{} = nil, nil
	// ch_kill will be nil in production
	if ch_kill == nil {
		ping, pong = ka.Subscribe("StoreToSqlite", 30)
	}

	// If we aren't logging events...
	event_ndx := 0
	// We'll use an in-memory DB for the recording of data throughout the day.
	db := newInMemoryDB()

	for {
		select {
		case <-ping:
			pong <- "StoreToSqlite"
			// This is the [ uid -> ticks ] map (uid looks like "Next:0")
		case <-ch_kill:
			db.Close()
			log.Println("Exiting StoreToSqlite")
			return

		case <-ch_reset:
			// Process the data from the day.
			log.Println("sqlite: processing data from the day")
			processDataFromDay(cfg, db)
			log.Println("sqlite: resetting the in-memory db")
			clearInMemoryDB(db)

		case arr := <-ch_data:
			log.Println("sqlite: storing data into memory")
			// This has to be done on the db that is currently open.
			// Cannot pre-prepare for the entire process.
			insertS, err := db.Prepare(`INSERT INTO wifi (event_id, fcfs_seq_id, device_tag, localtime, session_id, manufacturer_index, patron_index) VALUES (?,?,?,?,?,?,?)`)
			if err != nil {
				log.Fatal("sqlite: could not prepare insert statement.")
			}

			for _, h := range arr {
				//log.Println("inserting...")
				// log.Println(h["event_id"], h["fcfs_seq_id"], h["device_tag"], h["localtime"], h["session_id"], h["manufacturer_index"], h["patron_index"])
				res, err := insertS.Exec(h["event_id"], h["fcfs_seq_id"], h["device_tag"],
					h["localtime"], h["session_id"], h["manufacturer_index"], h["patron_index"])
				if err != nil {
					log.Println("sqlite: could not insert into memory db")
					log.Println(res)
					log.Fatal(err.Error())
				}
			}

			event_ndx += 1
		}
	}
}
