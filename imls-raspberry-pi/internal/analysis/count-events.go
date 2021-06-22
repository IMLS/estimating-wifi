package analysis

import (
	"sort"
	"time"

	"gsa.gov/18f/config"
)

// These defaults get overridden by the config.Config file.
var patron_min_mins float64 = 30
var patron_max_mins float64 = 10 * 60

const (
	Patron = iota
	Device
	Transient
)

type Counter struct {
	Patrons          int
	Devices          int
	Transients       int
	PatronMinutes    int
	DeviceMinutes    int
	TransientMinutes int
}
type ByStart []*Duration

func (a ByStart) Len() int { return len(a) }
func (a ByStart) Less(i, j int) bool {
	it, _ := time.Parse(time.RFC3339, a[i].Start)
	jt, _ := time.Parse(time.RFC3339, a[j].Start)
	return it.Before(jt)
}
func (a ByStart) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

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

func NewCounter(minMinutes int, maxMinutes int) *Counter {
	patron_min_mins = float64(minMinutes)
	patron_max_mins = float64(maxMinutes)
	return &Counter{0, 0, 0, 0, 0, 0}
}

func (c *Counter) Add(field int, minutes int) {
	switch field {
	case Patron:
		c.Patrons += 1
		c.PatronMinutes += minutes
	case Device:
		c.Devices += 1
		c.DeviceMinutes += minutes
	case Transient:
		c.Transients += 1
		c.TransientMinutes += minutes
	}
}

func getDeviceType(p WifiEvent, es []WifiEvent) int {
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
	if diff < patron_min_mins {
		return Transient
	} else if diff > patron_max_mins {
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

func getEventIdTime(events []WifiEvent, eventId int) (t time.Time) {
	for _, e := range events {
		if e.EventId == eventId {
			t = e.Localtime
			break
		}
	}

	return t
}

func doCounting(cfg *config.Config, events []WifiEvent) *Counter {
	c := NewCounter(cfg.Monitoring.MinimumMinutes, cfg.Monitoring.MaximumMinutes)

	prevEvent := events[0]
	checked := make(map[int]bool)
	for _, e := range events {
		// If the event id changes, bump our y pointer down.
		if e.EventId != prevEvent.EventId {
			prevEvent = e
		}
		if _, ok := checked[e.PatronIndex]; ok {
			// Skip if we already checked this patron
		} else {
			checked[e.PatronIndex] = true
			isP := getDeviceType(e, events)
			switch isP {
			case Patron:
				first, last := getPatronFirstLast(e.PatronIndex, events)
				firstTime := getEventIdTime(events, first)
				lastTime := getEventIdTime(events, last)
				minutes := int(lastTime.Sub(firstTime).Minutes())
				c.Add(Patron, minutes)
			case Device:
				first, last := getPatronFirstLast(e.PatronIndex, events)
				firstTime := getEventIdTime(events, first)
				lastTime := getEventIdTime(events, last)
				minutes := int(lastTime.Sub(firstTime).Minutes())
				c.Add(Device, minutes)
			case Transient:
				first, last := getPatronFirstLast(e.PatronIndex, events)
				firstTime := getEventIdTime(events, first)
				lastTime := getEventIdTime(events, last)
				minutes := int(lastTime.Sub(firstTime).Minutes())
				if minutes <= 0 {
					minutes = 1
				}
				c.Add(Transient, minutes)
			}
		}
	}

	return c
}

func durationSummary(events []WifiEvent) map[int]*Duration {

	// We want, for every patron_id, to know when the device started/ended.
	checked := make(map[int]bool)
	durations := make(map[int]*Duration)

	for _, e := range events {
		//log.Println("Patron index:", e.PatronIndex)
		if _, ok := checked[e.PatronIndex]; ok {
			// Skip if we already checked this patron
		} else {
			checked[e.PatronIndex] = true
			first, last := getPatronFirstLast(e.PatronIndex, events)
			firstTime := getEventIdTime(events, first)
			lastTime := getEventIdTime(events, last)
			durations[e.PatronIndex] = &Duration{PatronId: e.PatronIndex, MfgId: e.ManufacturerIndex, Start: firstTime.Format(time.RFC3339), End: lastTime.Format(time.RFC3339)}
		}
	}

	return durations
}

// Return the drawing context where the image is drawn.
// This can then be written to disk.
func Summarize(cfg *config.Config, events []WifiEvent) (c *Counter, d map[int]*Duration) {
	sort.Slice(events, func(i, j int) bool {
		return events[i].ID < events[j].ID
	})
	c = doCounting(cfg, events)
	d = durationSummary(events)
	return c, d
}
