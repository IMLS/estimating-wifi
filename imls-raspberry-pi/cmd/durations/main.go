package main

import (
	"database/sql/driver"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"gsa.gov/18f/internal/analysis"
	"gsa.gov/18f/internal/config"
)

const PATRONMINMINS = 30
const PATRONMAXMINS = 10 * 60

func countEvents(events []analysis.WifiEvent) int {
	prev := events[0]
	counter := 1

	for _, e := range events {
		if prev.EventId != e.EventId {
			prev = e
			counter += 1
		}
	}

	return counter
}

func allPatronIds(events []analysis.WifiEvent) []int {
	d := make(map[int]bool)
	for _, e := range events {
		d[e.PatronIndex] = true
	}
	a := make([]int, 0)
	for k := range d {
		a = append(a, k)
	}
	return a
}

func countPatrons(events []analysis.WifiEvent) int {
	max := 0

	for _, e := range events {
		if e.PatronIndex > max {
			max = e.PatronIndex
		}
	}

	return max
}

func isPatron(p analysis.WifiEvent, es []analysis.WifiEvent) int {
	var earliest time.Time
	var latest time.Time

	earliest = time.Date(3000, 1, 1, 0, 0, 0, 0, time.UTC)
	latest = time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)

	for _, e := range es {
		if p.PatronIndex == e.PatronIndex {
			if e.Localtime.Before(earliest) {
				earliest = e.Localtime
			}
			if e.Localtime.After(latest) {
				latest = e.Localtime
			}
		}
	}

	diff := latest.Sub(earliest).Minutes()
	if diff < PATRONMINMINS {
		return analysis.Transient
	} else if diff > PATRONMAXMINS {
		// log.Println("id", p.PatronIndex, "diff", diff)
		return analysis.Device
	} else {
		// log.Println("patron", p)
		return analysis.Patron
	}
}

func getPatronFirstLast(patronId int, events []analysis.WifiEvent) (int, int) {
	first := 1000000000
	last := -1000000000

	for _, e := range events {
		if e.PatronIndex == patronId {
			if e.EventId < first {
				first = e.EventId
			}
			if e.EventId > last {
				last = e.EventId
			}
		}
	}

	return first, last
}

func getEventIdTime(events []analysis.WifiEvent, eventId int) (t time.Time) {
	for _, e := range events {
		if e.EventId == eventId {
			t = e.Localtime
			break
		}
	}
	return t
}

type NullTime struct {
	Time  time.Time
	Valid bool // Valid is true if Time is not NULL
}

// Scan implements the Scanner interface.
func (nt *NullTime) Scan(value interface{}) error {
	nt.Time, nt.Valid = value.(time.Time)
	return nil
}

// Value implements the driver Valuer interface.
func (nt NullTime) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nil
	}
	return nt.Time, nil
}

func readWifiEventsFromSqlite(path string, tzoffset int) []analysis.WifiEvent {
	db, err := sqlx.Open("sqlite3", fmt.Sprintf("%v?parseTime=true", path)) //?parseTime=true
	if err != nil {
		log.Fatal("could not open sqlite file.")
	}

	events := []analysis.WifiEvent{}
	rows, err := db.Query("SELECT * FROM wifi")
	if err != nil {
		log.Println("error in read query")
		log.Fatal(err)
	}
	for rows.Next() {
		e := analysis.WifiEvent{}
		var id int
		var lt string
		err = rows.Scan(&id, &e.EventId, &e.FCFSSeqId, &e.DeviceTag,
			&lt, &e.SessionId, &e.ManufacturerIndex, &e.PatronIndex)
		if err != nil {
			var st string

			err = rows.Scan(&id, &e.EventId, &e.FCFSSeqId, &e.DeviceTag,
				&lt, &st, &e.SessionId, &e.ManufacturerIndex, &e.PatronIndex)
			if err != nil {
				log.Println("failed to scan with 8 and 9 args. Exiting.")
				log.Fatal(err)
			}
		}
		e.Localtime, _ = time.Parse(time.RFC3339, lt)
		e.Localtime = e.Localtime.Add(time.Duration(tzoffset) * time.Hour)
		events = append(events, e)
	}

	return events
}

/*
type Duration struct {
	Id        int    `db:"id"`
	PiSerial  string `db:"pi_serial"`
	SessionId string `db:"session_id"`
	FCFSSeqId string `db:"fcfs_seq_id"`
	DeviceTag string `db:"device_tag"`
	PatronId  int    `db:"pid"`
	MfgId     int    `db:"mfgid"`
	Start     string `db:"start"`
	End       string `db:"end"`
}
*/

func createDurationsTable(db *sqlx.DB) {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS durations
			(id INTEGER PRIMARY KEY,
			pi_serial text,
			session_id text,
			fcfs_seq_id text,
			device_tag text,
			pid integer,
			mfgid integer,
			start text,
			end text,
			minutes integer
			)`)
	if err != nil {
		log.Println("error creating table")
		log.Fatal(err)
	}
}

func writeDurations(path string, durations []*analysis.Duration) {
	if len(durations) <= 0 {
		return
	}

	thedb := filepath.Join(path, fmt.Sprintf("%v-%v-durations.sqlite", durations[0].FCFSSeqId, durations[0].DeviceTag))
	out, err := sqlx.Open("sqlite3", thedb)
	if err != nil {
		log.Fatal("could not open durations sqlite file")
	}

	createDurationsTable(out)
	tx, _ := out.Begin()
	stat, err := tx.Prepare(`INSERT INTO durations
			(pi_serial, session_id, fcfs_seq_id, device_tag, pid, mfgid, start, end, minutes)
			VALUES  (?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		log.Println("cannot prepare transactional insert")
		log.Fatal(err)
	}
	defer stat.Close()

	for _, v := range durations {
		st, e := time.Parse(time.RFC3339, v.Start)
		if e != nil {
			log.Println("error parsing start time", st)
			log.Fatal(e)
		}
		et, e := time.Parse(time.RFC3339, v.End)
		if e != nil {
			log.Println("error parsing end time", st)
			log.Fatal(e)
		}
		minutes := int(et.Sub(st).Minutes())
		stat.Exec(v.PiSerial, v.SessionId, v.FCFSSeqId, v.DeviceTag, v.PatronId, v.MfgId, v.Start, v.End, minutes)
	}
	err = tx.Commit()
	if err != nil {
		log.Println("failure to commit transaction")
		log.Fatal(err)
	}
}

func reviseDurations(cfg *config.Config, path string, swap bool, newPid int, sessionId string, events []analysis.WifiEvent) ([]*analysis.Duration, int) {
	filtered := make([]analysis.WifiEvent, 0)
	for _, e := range events {
		if e.SessionId == sessionId {
			filtered = append(filtered, e)
		}
	}

	// A session may span multiple days in the old data.
	// This means we have one session id, but a device id might span multiple days.
	// We want this to be broken up so a single session becomes multiple sessions, one for each day.
	// As we roll, a device needs to be reset into multiple (new) devices.
	// Devices that span from one day to the next also need to be split.

	d, newPid := analysis.MultiDayDurations(cfg, swap, newPid, filtered)
	sorted := make([]*analysis.Duration, 0)

	for _, v := range d {
		sorted = append(sorted, v)
	}
	sort.Slice(sorted[:], func(i, j int) bool {
		return sorted[i].PatronId < sorted[j].PatronId
	})

	return sorted, newPid
}

func getSessions(events []analysis.WifiEvent) []string {
	uniq := make(map[string]string)
	sessions := make([]string, 0)
	for _, e := range events {
		uniq[e.SessionId] = string(e.SessionId)
	}

	for k, _ := range uniq {
		sessions = append(sessions, k)
	}
	return sessions
}

func getDurationSessions(in []*analysis.Duration) []string {
	esesmap := make(map[string]string)
	for _, v := range in {
		esesmap[v.SessionId] = v.SessionId
	}
	esses := make([]string, 0)
	for _, v := range esesmap {
		esses = append(esses, v)
	}
	return esses
}

func remapPidsPerSession(in []*analysis.Duration) []*analysis.Duration {
	sessions := getDurationSessions(in)
	new := make([]*analysis.Duration, 0)
	for _, s := range sessions {
		pid := 0
		durations := make([]*analysis.Duration, 0)
		for _, d := range in {
			if d.SessionId == s {
				durations = append(durations, d)
			}
		}
		sort.Slice(durations[:], func(i, j int) bool {
			return durations[i].PatronId < durations[j].PatronId
		})
		for _, v := range durations {
			// Now, finally, remap back down to 0-indexed.
			newv := analysis.Duration{}
			newv.DeviceTag = v.DeviceTag
			newv.End = v.End
			newv.FCFSSeqId = v.FCFSSeqId
			newv.MfgId = v.MfgId
			newv.PatronId = pid
			newv.PiSerial = v.PiSerial
			newv.SessionId = v.SessionId
			newv.Start = v.Start
			pid = pid + 1
			new = append(new, &newv)
		}
	}
	return new
}

func main() {
	dataPtr := flag.String("sqlite", "", "A raw SQLite datafile.")
	cfgPath := flag.String("config", "", "Path to valid config file. REQUIRED.")
	outPath := flag.String("dest", "", "Path to output directory.")
	swapPtr := flag.Bool("swap", true, "Swap Start/End times that are out of order... (O_o)")

	// Dynamically grab the timezone offset for the machine we're running on.
	// It will be more... useful than a fixed constant.
	t := time.Now()
	zone, offset := t.Zone()
	// Get the offset in hours, not seconds.
	offset = offset / (60 * 60)
	tzPtr := flag.Int("tz", offset, fmt.Sprintf("timezone offset (%v is %v)", zone, offset))

	// newstyleFlag := flag.Bool("new", false, "Draw new style waterfalls.")
	flag.Parse()

	if *cfgPath == "" || *outPath == "" {
		log.Fatal("Must provide valid config file and dest path.")
	}

	cfg, _ := config.ReadConfig(*cfgPath)

	events := readWifiEventsFromSqlite(*dataPtr, *tzPtr)

	if len(events) > 0 {
		thedb := filepath.Join(*outPath, fmt.Sprintf("%v-%v-durations.sqlite", events[0].FCFSSeqId, events[0].DeviceTag))
		if _, err := os.Stat(thedb); err == nil {
			os.Remove(thedb)
		}

		sessions := getSessions(events)
		pid := 0
		var revised []*analysis.Duration
		for _, s := range sessions {
			revised, pid = reviseDurations(cfg, *outPath, *swapPtr, pid, s, events)
			revised = remapPidsPerSession(revised)
			writeDurations(*outPath, revised)
		}
	} else {
		log.Fatal("No wifi events found. Exiting.")
		os.Exit(-1)
	}
}
