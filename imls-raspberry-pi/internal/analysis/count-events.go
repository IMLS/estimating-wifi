package analysis

import (
	"fmt"
	"log"
	"reflect"
	"sort"
	"strings"
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
	Id        int    `db:"id",sqlite:"INTEGER PRIMARY KEY AUTOINCREMENT"`
	PiSerial  string `db:"pi_serial",sqlite:"TEXT"`
	SessionId string `db:"session_id",sqlite:"TEXT"`
	FCFSSeqId string `db:"fcfs_seq_id",sqlite:"TEXT"`
	DeviceTag string `db:"device_tag",sqlite:"TEXT"`
	PatronId  int    `db:"patron_index",sqlite:"INTEGER"`
	MfgId     int    `db:"manufacturer_index",sqlite:"INTEGER"`
	Start     string `db:"start",sqlite:"DATE"`
	End       string `db:"end",sqlite:"DATE"`
}

func (d Duration) AsMap() map[string]interface{} {
	m := make(map[string]interface{})
	rt := reflect.TypeOf(d)
	if rt.Kind() != reflect.Struct {
		panic("bad type")
	}
	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		col := strings.Split(f.Tag.Get("db"), ",")[0]
		//t := strings.Split(f.Tag.Get("sqlite"), ",")[0]
		m[col] = reflect.ValueOf(f)
	}
	return m
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
			if e.ID < first {
				first = e.ID
			}
			if e.ID > last {
				last = e.ID
			}
		}
	}

	return first, last
}

func getEventIdTime(events []WifiEvent, eventId int) (t time.Time) {
	for _, e := range events {
		if e.ID == eventId {
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
		if e.ID != prevEvent.ID {
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

func durationSummary(cfg *config.Config, events []WifiEvent) map[int]*Duration {

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

			durations[e.PatronIndex] = &Duration{
				PiSerial:  cfg.Serial,
				SessionId: e.SessionId,
				FCFSSeqId: e.FCFSSeqId,
				DeviceTag: e.DeviceTag,
				PatronId:  e.PatronIndex,
				MfgId:     e.ManufacturerIndex,
				Start:     firstTime.Format(time.RFC3339),
				End:       lastTime.Format(time.RFC3339)}
		}
	}

	return durations
}

// func eod(t time.Time) time.Time {
// 	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 0, t.Location())
// }

func bod(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func tomorrow(t time.Time) time.Time {
	tomorrow := t.Add(24 * time.Hour)
	return tomorrow
}

// func nextd (t time.Time) time.Time {
// 	return t.Add(24 * time.Hours)
// }

// FIXME: note the swap of times when things are borked in the DB...
// is that the best way to fix things?
func MultiDayDurations(cfg *config.Config, swap bool, newPid int, events []WifiEvent) (map[int]*Duration, int) {

	// We want, for every patron_id, to know when the device started/ended.
	checked := make(map[int]bool)
	durations := make(map[int]*Duration)
	uniqsessions := make(map[string]string)

	// Get the largest patron ID in this set. Use it for new,
	// unique patron IDs.
	maxPatronId := -1
	for _, e := range events {
		if e.PatronIndex > maxPatronId {
			maxPatronId = e.PatronIndex
		}
	}
	maxPatronId += 1

	for _, e := range events {
		// For later
		uniqsessions[e.SessionId] = e.SessionId

		//log.Println("Patron index:", e.PatronIndex)
		if _, ok := checked[e.PatronIndex]; ok {
			// Skip if we already checked this patron
		} else {
			checked[e.PatronIndex] = true
			first, last := getPatronFirstLast(e.PatronIndex, events)
			firstTime := getEventIdTime(events, first)
			lastTime := getEventIdTime(events, last)
			if lastTime.Before(firstTime) {
				log.Println("start", firstTime, "end", lastTime)
				log.Println("cannot start after end! swapping...")
				if swap {
					tmp := lastTime
					lastTime = firstTime
					firstTime = tmp
				}
			}
			// Only process these if the times are in the right order...
			if firstTime.Before(lastTime) {
				firstDay := firstTime.Day()
				lastDay := lastTime.Day()

				if lastDay > firstDay {
					// If this patron spans multiple days, they need to be split into multiple durations, and
					// they need to be given new, unique patron IDs for each day/duration.
					for lastDay > firstDay {
						//duration := lastTime.Sub(firstTime)
						//slog.Println("splitting", int(float64(duration)/float64(time.Hour)), "hour device", firstDay, lastDay)

						//log.Println("id", e.PatronIndex, "ft", firstTime.Format(time.RFC3339), "lt", lastTime.Format(time.RFC3339))
						// First, bump the end of this session to the end of today.
						endOfToday := eod(firstTime)
						// Insert the duration between the firstTime and the endOfToday, with a unique id.
						//log.Println("splitting", e.SessionId, e.PatronIndex, "to", maxPatronId)
						durations[maxPatronId] = &Duration{
							PiSerial:  cfg.Serial,
							SessionId: e.SessionId,
							FCFSSeqId: e.FCFSSeqId,
							DeviceTag: e.DeviceTag,
							PatronId:  maxPatronId,
							MfgId:     e.ManufacturerIndex,
							Start:     firstTime.Format(time.RFC3339),
							End:       endOfToday.Format(time.RFC3339)}
						maxPatronId += 1
						firstTime = bod(tomorrow(firstTime))
						firstDay = firstTime.Day()

					}

					//duration := lastTime.Sub(firstTime)
					//log.Println("last split", int(float64(duration)/float64(time.Hour)), "hour device", e.PatronIndex, maxPatronId, firstDay, lastDay)
					endOfToday := lastTime
					firstTime = bod(lastTime)
					// When done looping, insert the remainder...
					durations[maxPatronId] = &Duration{
						PiSerial:  cfg.Serial,
						SessionId: e.SessionId,
						FCFSSeqId: e.FCFSSeqId,
						DeviceTag: e.DeviceTag,
						PatronId:  maxPatronId,
						MfgId:     e.ManufacturerIndex,
						Start:     firstTime.Format(time.RFC3339),
						End:       endOfToday.Format(time.RFC3339)}
					maxPatronId += 1

				} else {
					durations[e.PatronIndex] = &Duration{
						PiSerial:  cfg.Serial,
						SessionId: e.SessionId,
						FCFSSeqId: e.FCFSSeqId,
						DeviceTag: e.DeviceTag,
						PatronId:  e.PatronIndex,
						MfgId:     e.ManufacturerIndex,
						Start:     firstTime.Format(time.RFC3339),
						End:       lastTime.Format(time.RFC3339)}
				}
			}

		} // end else
	}

	// These are now in need of patron renumbering, because new patrons were introduced.
	// This leaves patron ID gaps. Lets keep those monotonic and gapless.
	pid := 0
	remapped := make(map[int]*Duration)
	sorted := make([]*Duration, 0)
	for _, v := range durations {
		sorted = append(sorted, v)
	}
	sort.Slice(sorted[:], func(i, j int) bool {
		st, _ := time.Parse(time.RFC3339, sorted[i].Start)
		et, _ := time.Parse(time.RFC3339, sorted[i].End)
		return st.Before(et)
	})

	for _, v := range sorted {
		remapped[pid] = v
		pid += 1
	}

	// No durations span days at this point. Now, it would be nice if we could rewrite sessions into days.
	// This way, a "session" is a day. This tracks with how libraries think of things.
	// This requires, as a first step, globally unique PIDs.

	// We don't care what the session ID is, because at this point, a given pid only occurs once in a
	// given session. That means we can just increment a global PID for renumbering.
	pid = newPid
	newmap := make(map[int]*Duration)
	for _, v := range remapped {
		newv := Duration{}
		newv.DeviceTag = v.DeviceTag
		newv.End = v.End
		newv.FCFSSeqId = v.FCFSSeqId
		newv.MfgId = v.MfgId
		//log.Println(v.PatronId, "becomes", pid)
		newv.PatronId = pid
		newv.PiSerial = v.PiSerial
		st, _ := time.Parse(time.RFC3339, v.Start)
		day := fmt.Sprintf("%v%v%v", st.Year(), fmt.Sprintf("%02d", int(st.Month())), fmt.Sprintf("%02d", int(st.Day())))
		newv.SessionId = day
		newv.Start = v.Start
		newmap[pid] = &newv
		pid = pid + 1
	}

	return newmap, pid
}

// Return the drawing context where the image is drawn.
// This can then be written to disk.
func Summarize(cfg *config.Config, events []WifiEvent) (c *Counter, d map[int]*Duration) {
	sort.Slice(events, func(i, j int) bool {
		return events[i].ID < events[j].ID
	})
	c = doCounting(cfg, events)
	d = durationSummary(cfg, events)
	return c, d
}
