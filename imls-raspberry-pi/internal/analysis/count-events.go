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
	Transient = iota
	Patron
	Device
)

type Counter struct {
	Patrons          int
	Devices          int
	Transients       int
	PatronMinutes    int
	DeviceMinutes    int
	TransientMinutes int
}

func NewCounter(cfg *config.Config) *Counter {
	patron_min_mins = float64(cfg.Monitoring.MinimumMinutes)
	patron_max_mins = float64(cfg.Monitoring.MaximumMinutes)
	return &Counter{0, 0, 0, 0, 0, 0}
}

func (c *Counter) add(field int, minutes int) {
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
	c := NewCounter(cfg)

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
			isP := isPatron(e, events)
			switch isP {
			case Patron:
				first, last := getPatronFirstLast(e.PatronIndex, events)
				firstTime := getEventIdTime(events, first)
				lastTime := getEventIdTime(events, last)
				minutes := int(lastTime.Sub(firstTime).Minutes())
				c.add(Patron, minutes)
			case Device:
				first, last := getPatronFirstLast(e.PatronIndex, events)
				firstTime := getEventIdTime(events, first)
				lastTime := getEventIdTime(events, last)
				minutes := int(lastTime.Sub(firstTime).Minutes())
				c.add(Device, minutes)
			case Transient:
				first, last := getPatronFirstLast(e.PatronIndex, events)
				firstTime := getEventIdTime(events, first)
				lastTime := getEventIdTime(events, last)
				minutes := int(lastTime.Sub(firstTime).Minutes())
				if minutes <= 0 {
					minutes = 1
				}
				c.add(Transient, minutes)
			}
		}
	}

	return c
}

// Return the drawing context where the image is drawn.
// This can then be written to disk.
func Summarize(cfg *config.Config, events []WifiEvent) (c *Counter) {
	sort.Slice(events, func(i, j int) bool {
		return events[i].ID < events[j].ID
	})
	c = doCounting(cfg, events)
	return c
}
