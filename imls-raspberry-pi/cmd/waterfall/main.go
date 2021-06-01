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
	"github.com/jszwec/csvutil"
	. "gsa.gov/18f/analysis"
)

const PATRONMINMINS = 30
const PATRONMAXMINS = 10 * 60
const (
	Transient = iota
	Patron
	Device
)

func countEvents(events []WifiEvent) int {
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

func allPatronIds(events []WifiEvent) []int {
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

func countPatrons(events []WifiEvent) int {
	max := 0

	for _, e := range events {
		if e.PatronIndex > max {
			max = e.PatronIndex
		}
	}

	return max
}

func isPatron(p WifiEvent, es []WifiEvent) int {
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
		return Transient
	} else if diff > PATRONMAXMINS {
		// log.Println("id", p.PatronIndex, "diff", diff)
		return Device
	} else {
		// log.Println("patron", p)
		return Patron
	}
}

func getPatronFirstLast(patronId int, events []WifiEvent) (int, int) {
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

func modColor(v int, p []color.Color) color.Color {
	return p[v%len(p)]
}

func main() {
	csvPtr := flag.String("csv", "", "A CSV datafile.")
	flag.Parse()

	if *csvPtr == "" {
		log.Fatal("no CSV file provided.")
		os.Exit(-1)
	}

	b, err := ioutil.ReadFile(*csvPtr)
	if err != nil {
		log.Fatal("could not open CSV file.")
	}

	var events []WifiEvent
	if err := csvutil.Unmarshal(b, &events); err != nil {
		log.Println(err)
		log.Fatal("could not unmarshal CSV file as wifi events.")
	}

	// Event ids are our measure of y, patron ids are x.
	// Colors come from mfg?
	y := 0
	x := 0

	width := countPatrons(events) + 1
	height := countEvents(events) + 1
	log.Println("width", width, "height", height)
	// allIds := allPatronIds(events)
	sort.Slice(events, func(i, j int) bool {
		return events[i].ID < events[j].ID
	})

	// This creates an infinite sized image of uniform color.
	// img := &image.Uniform(color.RGBA(0x00, 0x00, 0x00, 0x00))
	dc := gg.NewContext(width, height)

	dc.SetRGBA(0.5, 0.5, 0, 0.5)
	dc.SetLineWidth(1)
	dc.DrawStringAnchored(fmt.Sprint(events[0].FCFSSeqId, " ", events[0].DeviceTag, " ", events[0].Localtime.Format("2006-01-02")), float64(width-5), float64(1), 1, 1)

	// Draw hour lines.
	hoursSeen := make(map[int]bool)

	y = 0
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

	// Draw the points
	y = 0
	prevEvent = events[0]
	for _, e := range events {
		// If the event id changes, bump our y pointer down.
		if e.EventId != prevEvent.EventId {
			y += 1
			prevEvent = e
		}
		x = e.PatronIndex
		isP := isPatron(e, events)
		if isP == Device {
			dc.SetRGBA(0.5, 0.5, 0.5, 0.5)
			dc.DrawPoint(float64(x), float64(y), 2)
			dc.Fill()
		} else if isP == Transient {
			dc.SetRGBA(0.75, 0.75, 0.75, 0.5)
			dc.DrawPoint(float64(x), float64(y), 2)
			dc.Fill()
		} else {
			// Don't draw patron points... draw lines?
			dc.SetColor(modColor(e.PatronIndex, palette.WebSafe))
		}

	}

	// Draw patron lines
	// Map event IDs to y
	eventIdToY := make(map[int]int)
	prevEvent = events[0]
	y = 0
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
			if isP == Patron {
				x = e.PatronIndex
				first, last := getPatronFirstLast(e.PatronIndex, events)

				dc.SetColor(modColor(e.PatronIndex, palette.WebSafe))
				dc.DrawLine(float64(x), float64(eventIdToY[first]), float64(x), float64(eventIdToY[last]))
				//dc.DrawPoint(float64(x), float64(y), 2)
				//dc.Fill()
				dc.SetLineWidth(3)
				dc.Stroke()
			}
		}
	}

	sid := events[0].SessionId
	seqId := events[0].FCFSSeqId
	dt := events[0].DeviceTag
	_ = os.Mkdir("output", 0777)
	fcfs_tag := fmt.Sprintf("%v-%v", seqId, dt)
	outdir := filepath.Join("output", fcfs_tag)
	_ = os.Mkdir(outdir, 0777)
	err = dc.SavePNG(fmt.Sprint(filepath.Join(outdir, fmt.Sprintf("%v-%v-%v.png", sid, seqId, dt))))
	if err != nil {
		log.Println("failed to save png")
		log.Fatal(err)
	}
}
