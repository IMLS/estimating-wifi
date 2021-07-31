package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/jszwec/csvutil"
	_ "github.com/mattn/go-sqlite3"
	"gsa.gov/18f/internal/analysis"
	"gsa.gov/18f/internal/state"
	"gsa.gov/18f/internal/structs"
)

func readWifiEventsFromCSV(path string) []structs.Duration {

	b, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal("could not open CSV file.")
	}

	var events []structs.Duration
	if err := csvutil.Unmarshal(b, &events); err != nil {
		log.Println(err)
		log.Fatal("could not unmarshal CSV file as wifi events.")
	}

	return events
}

func buildImagePath(fcfs string, deviceTag string, pngName string) string {

	_ = os.Mkdir("output", 0777)
	fcfsTag := fmt.Sprintf("%v-%v", fcfs, deviceTag)
	outdir := filepath.Join("output", fcfsTag)
	_ = os.Mkdir(outdir, 0777)
	baseFilename := fmt.Sprint(filepath.Join(outdir, pngName))
	fullPath := fmt.Sprintf("%v.png", baseFilename)

	return fullPath
}

func readDurationsFromSqlite(path string) []structs.Duration {
	// db, err := sqlx.Open("sqlite3", path)
	// if err != nil {
	// 	log.Fatal("could not open sqlite file.")
	// }
	db := state.NewSqliteDB(path)
	durations := []structs.Duration{}
	err := db.GetPtr().Select(&durations, "SELECT * FROM durations")
	if err != nil {
		log.Fatal(err.Error())
	}
	// rows, err := db.Query("SELECT *, cast((JulianDay(end) - JulianDay(start)) * 24 * 60 as integer) as minutes FROM durations")

	return durations
}

func main() {
	srcPtr := flag.String("src", "", "A source datafile (sqlite or csv).")
	cfgPath := flag.String("config", "", "Path to valid config file. REQUIRED.")
	typeFlag := flag.String("type", "sqlite", "Either 'csv' or 'sqlite' for source data")
	flag.Parse()

	if *cfgPath == "" {
		log.Fatal("Must provide valid config file.")
	}

	state.SetConfigAtPath(*cfgPath)
	cfg := state.GetConfig()

	var durations []structs.Duration
	if *typeFlag == "sqlite" {
		durations = readDurationsFromSqlite(*srcPtr)
	} else {
		durations = readWifiEventsFromCSV(*srcPtr)
	}

	sessions := make(map[string]string)
	for _, d := range durations {
		sessions[d.SessionID] = d.SessionID
	}

	for _, s := range sessions {
		subset := make([]structs.Duration, 0)
		for _, d := range durations {
			if d.SessionID == s {
				subset = append(subset, d)
			}
		}

		fcfs := subset[0].FCFSSeqID
		dt := subset[0].DeviceTag
		// This is necessary... in case we're testing with a
		// bogus config.sqlite file. Better to pull the identifiers from
		// the actual event stream than trust the file passed.
		cfg.SetFCFSSeqID(fcfs)
		cfg.SetDeviceTag(dt)
		pngName := fmt.Sprintf("%v-%v-%v-patron-sessions", subset[0].SessionID, fcfs, dt)
		log.Println("writing to", pngName)
		analysis.DrawPatronSessions(subset, buildImagePath(fcfs, dt, pngName))
	}
}
