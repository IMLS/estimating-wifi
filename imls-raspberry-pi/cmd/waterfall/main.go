package main

import (
	"flag"
	"fmt"
	"image/color"
	"image/color/palette"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/fogleman/gg"
	"github.com/jmoiron/sqlx"
	"github.com/jszwec/csvutil"
	_ "github.com/mattn/go-sqlite3"
	"gsa.gov/18f/internal/analysis"
	"gsa.gov/18f/internal/config"
	"gsa.gov/18f/internal/structs"
)

const PATRONMINMINS = 30
const PATRONMAXMINS = 10 * 60

func countEvents(events []structs.Duration) int {
	prev := events[0]
	counter := 1

	for _, e := range events {
		if prev.Id != e.Id {
			prev = e
			counter += 1
		}
	}

	return counter
}

func allPatronIds(events []structs.Duration) []int {
	d := make(map[int]bool)
	for _, e := range events {
		d[e.PatronId] = true
	}
	a := make([]int, 0)
	for k := range d {
		a = append(a, k)
	}
	return a
}

func countPatrons(events []structs.Duration) int {
	max := 0

	for _, e := range events {
		if e.PatronId > max {
			max = e.PatronId
		}
	}

	return max
}

func isPatron(p structs.Duration, es []structs.Duration) int {
	var earliest time.Time
	var latest time.Time

	earliest = time.Date(3000, 1, 1, 0, 0, 0, 0, time.UTC)
	latest = time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)

	for _, e := range es {
		if p.PatronId == e.PatronId {
			start, _ := time.Parse(time.RFC3339, e.Start)
			if start.Before(earliest) {
				earliest = start
			}
			if start.After(latest) {
				latest = start
			}
		}
	}

	diff := latest.Sub(earliest).Minutes()
	if diff < PATRONMINMINS {
		return analysis.Transient
	} else if diff > PATRONMAXMINS {
		// log.Println("id", p.PatronId, "diff", diff)
		return analysis.Device
	} else {
		// log.Println("patron", p)
		return analysis.Patron
	}
}

func getPatronFirstLast(patronId int, events []structs.Duration) (int, int) {
	first := 1000000000
	last := -1000000000

	for _, e := range events {
		if e.PatronId == patronId {
			if e.Id < first {
				first = e.Id
			}
			if e.Id > last {
				last = e.Id
			}
		}
	}

	return first, last
}

func getEventIdTime(events []structs.Duration, eventId int) (t time.Time) {
	for _, e := range events {
		if e.Id == eventId {
			start, _ := time.Parse(time.RFC3339, e.Start)
			t = start
			break
		}
	}
	return t
}

func modColor(v int, p []color.Color) color.Color {
	return p[v%len(p)]
}

func drawLines(events []structs.Duration, dc *gg.Context, width int, height int) {

	// Draw hour lines.
	hoursSeen := make(map[int]bool)

	y := 0

	prevEvent := events[0]
	for _, e := range events {
		if e.Id != prevEvent.Id {
			y += 1
			prevEvent = e
		}
		start, _ := time.Parse(time.RFC3339, e.Start)
		currentHour := start.Hour()
		if _, ok := hoursSeen[currentHour]; ok {
			// Skip
		} else {
			hoursSeen[currentHour] = true
			if currentHour%12 == 0 {
				dc.SetRGBA(0.5, 0.0, 0, 0.5)
				dc.SetLineWidth(2)
			} else {
				dc.SetRGBA(0.5, 0.5, 0, 0.5)
				dc.SetLineWidth(1)
				dc.DrawStringAnchored(fmt.Sprint(currentHour), float64(width-5), float64(y), 1, 1)
				dc.DrawStringAnchored(fmt.Sprint(currentHour), float64(width)*0.5, float64(y), 1, 1)
				dc.DrawStringAnchored(fmt.Sprint(currentHour), float64(50), float64(y), 1, 1)
			}
			dc.DrawLine(0, float64(y), float64(width), float64(y))
			dc.Stroke()
		}

	}
}

func drawPoints(events []structs.Duration, dc *gg.Context, width int, height int) {

	// Draw the points
	y := 0
	x := 0
	prevEvent := events[0]
	for _, e := range events {
		// If the event id changes, bump our y pointer down.
		if e.Id != prevEvent.Id {
			y += 1
			prevEvent = e
		}
		x = e.PatronId
		isP := isPatron(e, events)
		if isP == analysis.Device {
			dc.SetRGBA(0.5, 0.5, 0.5, 0.5)
			dc.DrawPoint(float64(x), float64(y), 2)
			dc.Fill()
		} else if isP == analysis.Transient {
			dc.SetRGBA(0.75, 0.75, 0.75, 0.5)
			dc.DrawPoint(float64(x), float64(y), 2)
			dc.Fill()
		} else {
			// Don't draw patron points... draw lines?
			dc.SetColor(modColor(e.PatronId, palette.WebSafe))
		}
	}
}

func drawPatronLines(events []structs.Duration, dc *gg.Context, c *analysis.Counter, width int, height int) {
	x := 0
	y := 0

	// Draw patron lines
	// Map event IDs to y
	eventIdToY := make(map[int]int)
	prevEvent := events[0]
	eventIdToY[prevEvent.Id] = y
	for _, e := range events {
		if e.Id != prevEvent.Id {
			y += 1
			prevEvent = e
			eventIdToY[prevEvent.Id] = y
		}
	}

	prevEvent = events[0]
	y = 0
	drawnPatrons := make(map[int]bool)
	for _, e := range events {
		// If the event id changes, bump our y pointer down.
		if e.Id != prevEvent.Id {
			y += 1
			prevEvent = e
		}
		if _, ok := drawnPatrons[e.PatronId]; ok {
			// Skip if we already checked this patron
		} else {
			drawnPatrons[e.PatronId] = true
			isP := isPatron(e, events)
			switch isP {
			case analysis.Patron:
				x = e.PatronId
				first, last := getPatronFirstLast(e.PatronId, events)
				firstTime := getEventIdTime(events, first)
				lastTime := getEventIdTime(events, last)
				minutes := int(lastTime.Sub(firstTime).Minutes())
				c.Add(analysis.Patron, minutes)
				dc.SetColor(modColor(e.PatronId, palette.WebSafe))
				dc.DrawLine(float64(x), float64(eventIdToY[first]), float64(x), float64(eventIdToY[last]))
				//dc.DrawPoint(float64(x), float64(y), 2)
				//dc.Fill()
				dc.SetLineWidth(3)
				dc.Stroke()
			case analysis.Device:
				first, last := getPatronFirstLast(e.PatronId, events)
				firstTime := getEventIdTime(events, first)
				lastTime := getEventIdTime(events, last)
				minutes := int(lastTime.Sub(firstTime).Minutes())
				c.Add(analysis.Device, minutes)
			case analysis.Transient:
				first, last := getPatronFirstLast(e.PatronId, events)
				firstTime := getEventIdTime(events, first)
				lastTime := getEventIdTime(events, last)
				minutes := int(lastTime.Sub(firstTime).Minutes())
				if minutes <= 0 {
					minutes = 1
				}
				c.Add(analysis.Transient, minutes)
			}
		}
	}
}

// Return the drawing context where the image is drawn.
// This can then be written to disk.
func drawWaterfall(events []structs.Duration) (c *analysis.Counter, dc *gg.Context) {

	// Event ids are our measure of y, patron ids are x.
	// Colors come from mfg?
	width := countPatrons(events) + 1
	height := countEvents(events) + 1
	log.Println("width", width, "height", height)
	// allIds := allPatronIds(events)
	sort.Slice(events, func(i, j int) bool {
		return events[i].Id < events[j].Id
	})

	// This creates an infinite sized image of uniform color.
	// img := &image.Uniform(color.RGBA(0x00, 0x00, 0x00, 0x00))
	dc = gg.NewContext(width, height)
	c = analysis.NewCounter(5, 300)

	dc.SetRGBA(0.5, 0.5, 0, 0.5)
	dc.SetLineWidth(1)

	start, _ := time.Parse(time.RFC3339, events[0].Start)
	str := fmt.Sprint(events[0].FCFSSeqId, " ", events[0].DeviceTag, " ", start.Format("2006-01-02"))
	dc.DrawStringAnchored(str, float64(width-5), float64(1), 1, 1)

	drawLines(events, dc, width, height)
	drawPoints(events, dc, width, height)
	// Get our minutes data from the line drawing, which is where
	// start and end duration points can be found/calculated.
	drawPatronLines(events, dc, c, width, height)

	return c, dc
}

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
	fcfs_tag := fmt.Sprintf("%v-%v", fcfs, deviceTag)
	outdir := filepath.Join("output", fcfs_tag)
	_ = os.Mkdir(outdir, 0777)
	baseFilename := fmt.Sprint(filepath.Join(outdir, pngName))
	fullPath := fmt.Sprintf("%v.png", baseFilename)

	return fullPath
}

func drawOldWaterfalls(events []structs.Duration, dataPtr *string) {

	// Capture the data about the session while running in a `counter` structure.
	c, dc := drawWaterfall(events)

	sid := events[0].SessionId
	seqId := events[0].FCFSSeqId
	dt := events[0].DeviceTag
	pngName := fmt.Sprintf("%v-%v-%v", sid, seqId, dt)
	fullPath := buildImagePath(seqId, dt, pngName)

	err := dc.SavePNG(fullPath)
	if err != nil {
		log.Println("failed to save png")
		log.Fatal(err)
	}

	// Write the count data, possibly.
	if *dataPtr != "" {
		outCSVFilename := *dataPtr
		writeHeader := false
		if _, err := os.Stat(outCSVFilename); os.IsNotExist(err) {
			writeHeader = true
		}
		f, err := os.OpenFile(outCSVFilename,
			os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		if writeHeader {
			f.WriteString("fcfs_seq_id,device_tag,session_id,patrons,patron_minutes,devices,device_minutes,transients,transient_minutes\n")
		}
		f.WriteString(fmt.Sprintf("%v,%v,%v,%v,%v,%v,%v,%v,%v\n",
			seqId, dt, sid, c.Patrons, c.PatronMinutes, c.Devices, c.DeviceMinutes, c.Transients, c.TransientMinutes))
	} else {
		fmt.Printf("%+v\n", c)
	}
}

func readDurationsFromSqlite(path string) []*structs.Duration {
	db, err := sqlx.Open("sqlite3", path)
	if err != nil {
		log.Fatal("could not open sqlite file.")
	}

	events := []*structs.Duration{}
	rows, err := db.Query("SELECT *, cast((JulianDay(end) - JulianDay(start)) * 24 * 60 as integer) as minutes FROM durations")
	if err != nil {
		log.Println("error in read query")
		log.Fatal(err)
	}
	for rows.Next() {
		d := structs.Duration{}
		var id int
		var minutes int
		err = rows.Scan(&id, &d.PiSerial, &d.SessionId, &d.FCFSSeqId, &d.DeviceTag, &d.PatronId, &d.MfgId, &d.Start, &d.End, &minutes)
		if err != nil {
			log.Println("could not scan")
			log.Fatal(err)
		}
		events = append(events, &d)
	}

	return events
}

func main() {
	srcPtr := flag.String("src", "", "A source datafile (sqlite or csv).")
	cfgPath := flag.String("config", "", "Path to valid config file. REQUIRED.")
	dataPtr := flag.String("data", "", "Data output.")
	typeFlag := flag.String("type", "sqlite", "Either 'csv' or 'sqlite' for source data")
	flag.Parse()

	if *cfgPath == "" {
		log.Fatal("Must provide valid config file.")
	}

	cfg, _ := config.NewConfigFromPath(*cfgPath)

	if *typeFlag == "sqlite" {
		durations := readDurationsFromSqlite(*srcPtr)
		sessions := make(map[string]string)
		for _, d := range durations {
			sessions[d.SessionId] = d.SessionId
		}

		for _, s := range sessions {
			subset := make([]structs.Duration, 0)
			for _, d := range durations {
				if d.SessionId == s {
					subset = append(subset, *d)
				}
			}

			fcfs := subset[0].FCFSSeqId
			dt := subset[0].DeviceTag
			d := subset[0].Start
			dtime, _ := time.Parse(time.RFC3339, d)
			// This is necessary... in case we're testing with a
			// bogus config.yaml file. Better to pull the identifiers from
			// the actual event stream than trust the file passed.
			cfg.Auth.FCFSId = fcfs
			cfg.Auth.DeviceTag = dt
			pngName := fmt.Sprintf("%v-%.2v-%.2v-%v-%v-patron-sessions", dtime.Year(), int(dtime.Month()), int(dtime.Day()), fcfs, dt)
			//log.Println("writing to", pngName)
			analysis.DrawPatronSessions(cfg, subset, buildImagePath(fcfs, dt, pngName))
		}

	} else {
		events := readWifiEventsFromCSV(*srcPtr)
		drawOldWaterfalls(events, dataPtr)
	}

}
