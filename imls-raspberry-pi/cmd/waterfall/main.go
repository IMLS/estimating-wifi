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
	"gsa.gov/18f/analysis"
	"gsa.gov/18f/config"
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

func modColor(v int, p []color.Color) color.Color {
	return p[v%len(p)]
}

func drawLines(events []analysis.WifiEvent, dc *gg.Context, width int, height int) {

	// Draw hour lines.
	hoursSeen := make(map[int]bool)

	y := 0

	prevEvent := events[0]
	for _, e := range events {
		if e.EventId != prevEvent.EventId {
			y += 1
			prevEvent = e
		}
		currentHour := e.Localtime.Hour()
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

func drawPoints(events []analysis.WifiEvent, dc *gg.Context, width int, height int) {

	// Draw the points
	y := 0
	x := 0
	prevEvent := events[0]
	for _, e := range events {
		// If the event id changes, bump our y pointer down.
		if e.EventId != prevEvent.EventId {
			y += 1
			prevEvent = e
		}
		x = e.PatronIndex
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
			dc.SetColor(modColor(e.PatronIndex, palette.WebSafe))
		}
	}
}

func drawPatronLines(events []analysis.WifiEvent, dc *gg.Context, c *analysis.Counter, width int, height int) {
	x := 0
	y := 0

	// Draw patron lines
	// Map event IDs to y
	eventIdToY := make(map[int]int)
	prevEvent := events[0]
	eventIdToY[prevEvent.EventId] = y
	for _, e := range events {
		if e.EventId != prevEvent.EventId {
			y += 1
			prevEvent = e
			eventIdToY[prevEvent.EventId] = y
		}
	}

	prevEvent = events[0]
	y = 0
	drawnPatrons := make(map[int]bool)
	for _, e := range events {
		// If the event id changes, bump our y pointer down.
		if e.EventId != prevEvent.EventId {
			y += 1
			prevEvent = e
		}
		if _, ok := drawnPatrons[e.PatronIndex]; ok {
			// Skip if we already checked this patron
		} else {
			drawnPatrons[e.PatronIndex] = true
			isP := isPatron(e, events)
			switch isP {
			case analysis.Patron:
				x = e.PatronIndex
				first, last := getPatronFirstLast(e.PatronIndex, events)
				firstTime := getEventIdTime(events, first)
				lastTime := getEventIdTime(events, last)
				minutes := int(lastTime.Sub(firstTime).Minutes())
				c.Add(analysis.Patron, minutes)
				dc.SetColor(modColor(e.PatronIndex, palette.WebSafe))
				dc.DrawLine(float64(x), float64(eventIdToY[first]), float64(x), float64(eventIdToY[last]))
				//dc.DrawPoint(float64(x), float64(y), 2)
				//dc.Fill()
				dc.SetLineWidth(3)
				dc.Stroke()
			case analysis.Device:
				first, last := getPatronFirstLast(e.PatronIndex, events)
				firstTime := getEventIdTime(events, first)
				lastTime := getEventIdTime(events, last)
				minutes := int(lastTime.Sub(firstTime).Minutes())
				c.Add(analysis.Device, minutes)
			case analysis.Transient:
				first, last := getPatronFirstLast(e.PatronIndex, events)
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
func drawWaterfall(events []analysis.WifiEvent) (c *analysis.Counter, dc *gg.Context) {

	// Event ids are our measure of y, patron ids are x.
	// Colors come from mfg?
	width := countPatrons(events) + 1
	height := countEvents(events) + 1
	log.Println("width", width, "height", height)
	// allIds := allPatronIds(events)
	sort.Slice(events, func(i, j int) bool {
		return events[i].ID < events[j].ID
	})

	// This creates an infinite sized image of uniform color.
	// img := &image.Uniform(color.RGBA(0x00, 0x00, 0x00, 0x00))
	dc = gg.NewContext(width, height)
	c = analysis.NewCounter(5, 300)

	dc.SetRGBA(0.5, 0.5, 0, 0.5)
	dc.SetLineWidth(1)
	dc.DrawStringAnchored(fmt.Sprint(events[0].FCFSSeqId, " ", events[0].DeviceTag, " ", events[0].Localtime.Format("2006-01-02")), float64(width-5), float64(1), 1, 1)

	drawLines(events, dc, width, height)
	drawPoints(events, dc, width, height)
	// Get our minutes data from the line drawing, which is where
	// start and end duration points can be found/calculated.
	drawPatronLines(events, dc, c, width, height)

	return c, dc
}

func drawOldWaterfalls(csvPtr *string, dataPtr *string) {

	b, err := ioutil.ReadFile(*csvPtr)
	if err != nil {
		log.Fatal("could not open CSV file.")
	}

	var events []analysis.WifiEvent
	if err := csvutil.Unmarshal(b, &events); err != nil {
		log.Println(err)
		log.Fatal("could not unmarshal CSV file as wifi events.")
	}

	// Capture the data about the session while running in a `counter` structure.
	c, dc := drawWaterfall(events)

	sid := events[0].SessionId
	seqId := events[0].FCFSSeqId
	dt := events[0].DeviceTag
	_ = os.Mkdir("output", 0777)
	fcfs_tag := fmt.Sprintf("%v-%v", seqId, dt)
	outdir := filepath.Join("output", fcfs_tag)
	_ = os.Mkdir(outdir, 0777)
	baseFilename := fmt.Sprint(filepath.Join(outdir, fmt.Sprintf("%v-%v-%v", sid, seqId, dt)))
	pngFilename := fmt.Sprintf("%v.png", baseFilename)
	err = dc.SavePNG(pngFilename)
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

// sqlx makes it easy to pull things in.
func getDurations(db *sqlx.DB, sid string) []analysis.Duration {
	durations := []analysis.Duration{}
	err := db.Select(&durations, `SELECT * FROM durations WHERE session_id=?`, sid)
	if err != nil {
		log.Println("sqlite: error in getDurations query.")
		log.Fatal(err.Error())
	}
	return durations
}

func getUniqueSessions(db *sqlx.DB) (sessions []string) {
	err := db.Select(&sessions, `SELECT DISTINCT(session_id) FROM durations`)
	if err != nil {
		log.Println("Failed to get unique sessions.")
	}

	return sessions
}

func _drawNewWaterfalls(cfg *config.Config, db *sqlx.DB) {
	WIDTH := 1200
	HEIGHT := 600

	// sid := durations[0].SessionId
	// seqId := durations[0].FCFSSeqId
	// dt := durations[0].DeviceTag
	// _ = os.Mkdir("output", 0777)
	// fcfs_tag := fmt.Sprintf("%v-%v", seqId, dt)
	// outdir := filepath.Join("output", fcfs_tag)
	// _ = os.Mkdir(outdir, 0777)
	// baseFilename := fmt.Sprint(filepath.Join(outdir, fmt.Sprintf("%v-%v-%v", sid, seqId, dt)))
	//pngFilename := fmt.Sprintf("%v.png", baseFilename)

	// This creates an infinite sized image of uniform color.
	// img := &image.Uniform(color.RGBA(0x00, 0x00, 0x00, 0x00))
	dc := gg.NewContext(WIDTH, HEIGHT)
	dc.SetRGBA(0.5, 0.5, 0, 0.5)
	dc.SetLineWidth(1)

	for _, sid := range getUniqueSessions(db) {
		for _, d := range getDurations(db, sid) {
			st, _ := time.Parse(time.RFC3339, d.Start)
			et, _ := time.Parse(time.RFC3339, d.End)
			diff := int(et.Sub(st).Minutes())
			// log.Println("st", st, "et", et, "diff", diff)
			if (diff > cfg.Monitoring.MinimumMinutes) && (diff < cfg.Monitoring.MaximumMinutes) {
				log.Println("d", d)

			}
		}
	}
}

func main() {
	csvPtr := flag.String("csv", "", "A CSV datafile.")
	summaryDbPtr := flag.String("summary", "", "A summary DB (SQLite)")
	cfgPath := flag.String("config", "", "Path to valid config file. REQUIRED.")
	dataPtr := flag.String("data", "", "Data output.")
	newstyleFlag := flag.Bool("newstyle", false, "Draw new style waterfalls.")
	flag.Parse()

	if *csvPtr != "" && *summaryDbPtr != "" {
		log.Fatal("Provide EITHER a CSV file OR a summary DB.")
		os.Exit(-1)
	}

	if *csvPtr == "" && *summaryDbPtr == "" {
		log.Fatal("Provide EITHER a CSV file OR a summary DB.")
		os.Exit(-1)
	}

	if *cfgPath == "" {
		log.Fatal("Must provide valid config file.")
	}

	if *csvPtr != "" && !*newstyleFlag {
		drawOldWaterfalls(csvPtr, dataPtr)
	}

	if *csvPtr != "" && *newstyleFlag {

		// FOR PROCESSING NEW-STYLE SUMMARY DBS... NOT YET READY

		// The new waterfalls use the summary DB.
		// 	insertS, err = summarydb.Prepare(`INSERT INTO durations (pi_serial, fcfs_seq_id, device_tag, session_id, pid, type, start, end) VALUES (?,?,?,?,?,?,?,?)`)
		// As a CSV, those are the headers.
		// if _, err := os.Stat(*summaryDbPtr); os.IsNotExist(err) {
		// 	log.Println("Could not find", *summaryDbPtr)
		// 	os.Exit(-1)
		// }

		// db, err := sqlx.Open("sqlite3", fmt.Sprintf("%v?parseTime=true", *summaryDbPtr))
		// if err != nil {
		// 	log.Fatal("sqlite: could not open summary db.")
		// }

		cfg, _ := config.ReadConfig(*cfgPath)
		//drawNewWaterfalls(cfg, db)
		//analysis.DrawPatronSessions(cfg, csvPtr, "test.png")
	}
}
