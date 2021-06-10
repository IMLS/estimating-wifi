package analysis

import (
	"sort"
	"time"
)

const PATRONMINMINS = 30
const PATRONMAXMINS = 10 * 60
const (
	Transient = iota
	Patron
	Device
)

type Counter struct {
	patrons           int
	devices           int
	transients        int
	patron_minutes    int
	device_minutes    int
	transient_minutes int
}

func NewCounter() *Counter {
	return &Counter{0, 0, 0, 0, 0, 0}
}

func (c *Counter) add(field int, minutes int) {
	switch field {
	case Patron:
		c.patrons += 1
		c.patron_minutes += minutes
	case Device:
		c.devices += 1
		c.device_minutes += minutes
	case Transient:
		c.transients += 1
		c.transient_minutes += minutes
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

func getEventIdTime(events []WifiEvent, eventId int) (t time.Time) {
	for _, e := range events {
		if e.EventId == eventId {
			t = e.Localtime
			break
		}
	}
	return t
}

func doCounting(events []WifiEvent) *Counter {
	c := NewCounter()

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
func Summarize(events []WifiEvent) (c *Counter) {
	sort.Slice(events, func(i, j int) bool {
		return events[i].ID < events[j].ID
	})
	c = doCounting(events)
	return c
}
