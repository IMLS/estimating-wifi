package tlp

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/jmoiron/sqlx"
	"gsa.gov/18f/analysis"
	"gsa.gov/18f/config"
	"gsa.gov/18f/logwrapper"
)

func getSummaryDB(cfg *config.Config) *sqlx.DB {
	lw := logwrapper.NewLogger(nil)

	if _, err := os.Stat(cfg.Local.SummaryDB); os.IsNotExist(err) {
		file, err := os.Create(cfg.Local.SummaryDB)
		if err != nil {
			lw.Info("could not create sqlite summary db file: %v", cfg.Local.SummaryDB)
			lw.Fatal(err.Error())
		}
		file.Close()
	}

	db, err := sqlx.Open("sqlite3", cfg.Local.SummaryDB)
	if err != nil {
		lw.Fatal("could not open summary db.")
	}

	// Create tables if it doesn't exist
	createTableStatement := `
	CREATE TABLE IF NOT EXISTS counts (
		id integer PRIMARY KEY AUTOINCREMENT,
		pi_serial text,
		fcfs_seq_id text,
		device_tag text,
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
		lw.Fatal("could not create counts table in summary db.")
	}

	createTableStatement = `
	CREATE TABLE IF NOT EXISTS durations (
		id integer PRIMARY KEY AUTOINCREMENT,
		pi_serial text,
		fcfs_seq_id text,
		device_tag text,
		session_id text,
		pid integer,
		mfg_id integer,
		start text,
		end text
	);`

	_, err = db.Exec(createTableStatement)
	if err != nil {
		lw.Fatal("could not create durations table in summary db.")
	}

	return db
}

func newTemporaryDB(cfg *config.Config) *sqlx.DB {
	lw := logwrapper.NewLogger(nil)

	t := time.Now()
	todaysDB := fmt.Sprintf("%v-%02d-%02d-wifi.sqlite", t.Year(), int(t.Month()), int(t.Day()))
	lw.Info("Created temporary db: %v", todaysDB)
	path := filepath.Join(cfg.Local.WebDirectory, todaysDB)
	db, err := sqlx.Open("sqlite3", path)
	if err != nil {
		lw.Fatal("could not open temporary db: %v", path)
	}

	createWifiTable(cfg, db)
	return db
}

func createWifiTable(cfg *config.Config, db *sqlx.DB) {
	createTableStatement := `
	CREATE TABLE IF NOT EXISTS wifi (
		id integer PRIMARY KEY AUTOINCREMENT,
		event_id integer,
		fcfs_seq_id text,
		device_tag text,
		"localtime" date,
		session_id text,
		manufacturer_index integer,
		patron_index integer
	);`

	_, err := db.Exec(createTableStatement)
	if err != nil {
		lw := logwrapper.NewLogger(cfg)
		lw.Info("could not create wifi table in temporary db.")
	}
}

// When in-memory temp dbs are reintroduced, we might want this.
// func clearTemporaryDB(cfg *config.Config, db *sqlx.DB) {
// 	// Create tables.
// 	createTableStatement := `
// 	DROP TABLE IF EXISTS wifi;
// 	CREATE TABLE wifi (
// 		id integer PRIMARY KEY AUTOINCREMENT,
// 		event_id integer,
// 		fcfs_seq_id text,
// 		device_tag text,
// 		"localtime" date,
// 		session_id text,
// 		manufacturer_index integer,
// 		patron_index integer
// 	);`

// 	_, err := db.Exec(createTableStatement)
// 	if err != nil {
// 		lw := logwrapper.NewLogger(cfg)
// 		lw.Fatal("could not re-create wifi table in temporary db.")
// 	}
// }

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

func storeSummary(cfg *config.Config, c *analysis.Counter, d map[int]*analysis.Duration) {
	summarydb := getSummaryDB(cfg)
	defer summarydb.Close()
	lw := logwrapper.NewLogger(nil)

	insertS, err := summarydb.Prepare(`INSERT INTO counts (pi_serial, fcfs_seq_id, device_tag, session_id, minimum_minutes, maximum_minutes, patron_count, patron_minutes, device_count, device_minutes, transient_count, transient_minutes) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)`)
	if err != nil {
		lw.Info("could not prepare `counts` insert statement")
		lw.Fatal(err.Error())
	}

	_, err = insertS.Exec(config.GetSerial(), cfg.Auth.FCFSId, cfg.Auth.DeviceTag, cfg.SessionId, cfg.Monitoring.MinimumMinutes, cfg.Monitoring.MaximumMinutes, c.Patrons, c.PatronMinutes, c.Devices, c.DeviceMinutes, c.Transients, c.TransientMinutes)
	if err != nil {
		lw.Info("could not insert into `counts` db")
		lw.Fatal(err.Error())
	}

	insertS, err = summarydb.Prepare(fmt.Sprintf(`INSERT INTO durations (%v,%v,%v,%v,%v,%v,%v,%v) VALUES (?,?,?,?,?,?,?,?)`,
		"pi_serial",
		"fcfs_seq_id",
		"device_tag",
		"session_id",
		"pid",
		"mfg_id",
		"start",
		"end"))

	if err != nil {
		lw.Info("could not prepare `durations` insert statement.")
		lw.Fatal(err.Error())
	}
	for pid, duration := range d {
		_, err := insertS.Exec(config.GetSerial(), cfg.Auth.FCFSId, cfg.Auth.DeviceTag, cfg.SessionId, pid, duration.MfgId, duration.Start, duration.End)
		if err != nil {
			lw.Info("could not insert into `durations` db")
			lw.Fatal(err.Error())
		}
	}
}

func processDataFromDay(cfg *config.Config, memdb *sqlx.DB) {
	events := extractWifiEvents(memdb)
	lw := logwrapper.NewLogger(nil)

	if len(events) > 0 {
		lw.Length("events", events)
		c, d := analysis.Summarize(cfg, events)
		writeImages(cfg, events)
		writeSummaryCSV(cfg, events)
		storeSummary(cfg, c, d)
	} else {
		lw.Info("no `events` to summarize")
	}
}

//This must happen after the data is updated for the day.
func writeImages(cfg *config.Config, events []analysis.WifiEvent) {

	yesterday := time.Now().Add(-1 * time.Hour)
	if _, err := os.Stat(cfg.Local.WebDirectory); os.IsNotExist(err) {
		err := os.Mkdir(cfg.Local.WebDirectory, 0777)
		if err != nil {
			if config.Verbose {
				log.Println("could not create web directory:", cfg.Local.WebDirectory)
			}
		}
	}
	imagedir := filepath.Join(cfg.Local.WebDirectory, "images")
	if _, err := os.Stat(imagedir); os.IsNotExist(err) {
		err := os.Mkdir(imagedir, 0777)
		if err != nil {
			if config.Verbose {
				log.Println("could not create image directory")
			}
		}
	}

	path := filepath.Join(imagedir, fmt.Sprintf("%04d-%02d-%02d-%v-%v-summary.png", yesterday.Year(), int(yesterday.Month()), int(yesterday.Day()), cfg.Auth.FCFSId, cfg.Auth.DeviceTag))
	analysis.DrawPatronSessionsFromWifi(cfg, events, path)
}

func writeSummaryCSV(cfg *config.Config, events []analysis.WifiEvent) {

	_, durations := analysis.Summarize(cfg, events)
	if _, err := os.Stat(cfg.Local.WebDirectory); os.IsNotExist(err) {
		err := os.Mkdir(cfg.Local.WebDirectory, 0777)
		if err != nil {
			if config.Verbose {
				log.Println("could not create web directory:", cfg.Local.WebDirectory)
			}
		}
	}
	path := filepath.Join(cfg.Local.WebDirectory,
		fmt.Sprintf("%v-%v-durations.csv", cfg.Auth.FCFSId, cfg.Auth.DeviceTag))

	// Open for appending
	f, err := os.OpenFile(path,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		if config.Verbose {
			log.Println("could not open durations CSV for writing")
		}
	}
	defer f.Close()

	// Write the header only when the file is opened for the first time.
	fi, err := f.Stat()
	if err == nil {
		if fi.Size() == 0 {
			f.WriteString("pi_serial,fcfs_seq_id,device_tag,session_id,patron_id,mfg_id,start,end,minutes\n")
		}
	}

	for _, d := range durations {
		st, _ := time.Parse(time.RFC3339, d.Start)
		et, _ := time.Parse(time.RFC3339, d.End)
		diff := int(et.Sub(st).Minutes())
		str := fmt.Sprintf("%v,%v,%v,%v,%v,%v,%v,%v,%v\n",
			d.PiSerial,
			d.FCFSSeqId,
			d.DeviceTag,
			d.SessionId,
			d.PatronId,
			d.MfgId,
			d.Start,
			d.End,
			diff)
		f.WriteString(str)
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
	db := newTemporaryDB(cfg)
	lw := logwrapper.NewLogger(nil)

	for {
		select {
		case <-ping:
			pong <- "StoreToSqlite"
			// This is the [ uid -> ticks ] map (uid looks like "Next:0")
		case <-ch_kill:
			db.Close()
			lw.Debug("exiting StoreToSqlite")
			return

		case <-ch_reset:
			lw.Info("receieved reset message")
			// Process the data from the day.
			processDataFromDay(cfg, db)
			//clearTemporaryDB(cfg, db)
			db.Close()
			db = newTemporaryDB(cfg)
			// After clearing, it is a new session.
			cfg.SessionId = config.CreateSessionId()

		case arr := <-ch_data:
			lw.Info("storing temporary data")
			// This has to be done on the db that is currently open.
			// Cannot pre-prepare for the entire process.
			insertS, err := db.Prepare(`INSERT INTO wifi (event_id, fcfs_seq_id, device_tag, localtime, session_id, manufacturer_index, patron_index) VALUES (?,?,?,?,?,?,?)`)
			if err != nil {
				log.Fatal("sqlite: could not prepare insert statement.")
				lw.Fatal("could not prepare wifi insert statement")
			}

			for _, h := range arr {
				_, err := insertS.Exec(h["event_id"], h["fcfs_seq_id"], h["device_tag"],
					h["localtime"], h["session_id"], h["manufacturer_index"], h["patron_index"])
				if err != nil {
					lw.Info("could not insert into temporary db")
					lw.Fatal(err.Error())
				}
			}

			event_ndx += 1
		}
	}
}
