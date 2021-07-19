package analysis

import (
	"fmt"
	"sort"
	"time"

	"gsa.gov/18f/internal/state"
	"gsa.gov/18f/internal/structs"
)

// These defaults get overridden by the config.Config file.
var patronMinMins float64 = 30
var patronMaxMins float64 = 10 * 60

const (
	Patron = iota
	Device
	Transient
)

func getDeviceType(p structs.WifiEvent, es []structs.WifiEvent) int {
	var earliest time.Time
	var latest time.Time

	earliest = time.Date(3000, 1, 1, 0, 0, 0, 0, time.UTC)
	latest = time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)

	for _, e := range es {
		elt, _ := time.Parse(time.RFC3339, e.Localtime)
		if p.PatronIndex == e.PatronIndex {
			if elt.Before(earliest) {
				earliest = elt
			}
			if elt.After(latest) {
				latest = elt
			}
		}
	}

	diff := latest.Sub(earliest).Minutes()
	if diff < patronMinMins {
		return Transient
	} else if diff > patronMaxMins {
		// log.Println("id", p.PatronIndex, "diff", diff)
		return Device
	} else {
		// log.Println("patron", p)
		return Patron
	}
}

func getPatronFirstLast(patronID int, events []structs.WifiEvent) (int, int) {
	first := 1000000000
	last := -1000000000

	for _, e := range events {
		if e.PatronIndex == patronID {
			if e.ID < first {
				first = e.ID
			}
			if e.Id > last {
				last = e.Id
			}
		}
	}

	return first, last
}

func getEventIDTime(events []structs.WifiEvent, eventID int) (t time.Time) {
	for _, e := range events {
		if e.ID == eventID {
			t, _ = time.Parse(time.RFC3339, e.Localtime)
			break
		}
	}

	return t
}

func doCounting(events []structs.WifiEvent) *Counter {
	cfg := state.GetConfig()
	c := NewCounter(cfg.GetMinimumMinutes(), cfg.GetMaximumMinutes())

	prevEvent := events[0]
	checked := make(map[int]bool)
	for _, e := range events {
		// If the event id changes, bump our y pointer down.
		if e.Id != prevEvent.Id {
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
				firstTime := getEventIDTime(events, first)
				lastTime := getEventIDTime(events, last)
				minutes := int(lastTime.Sub(firstTime).Minutes())
				c.Add(Patron, minutes)
			case Device:
				first, last := getPatronFirstLast(e.PatronIndex, events)
				firstTime := getEventIDTime(events, first)
				lastTime := getEventIDTime(events, last)
				minutes := int(lastTime.Sub(firstTime).Minutes())
				c.Add(Device, minutes)
			case Transient:
				first, last := getPatronFirstLast(e.PatronIndex, events)
				firstTime := getEventIDTime(events, first)
				lastTime := getEventIDTime(events, last)
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

func durationSummary(events []structs.WifiEvent) map[int]structs.Duration {
	cfg := state.GetConfig()
	// We want, for every patron_id, to know when the device started/ended.
	checked := make(map[int]bool)
	durations := make(map[int]structs.Duration)

	for _, e := range events {
		//log.Println("Patron index:", e.PatronIndex)
		if _, ok := checked[e.PatronIndex]; ok {
			// Skip if we already checked this patron
		} else {
			checked[e.PatronIndex] = true
			first, last := getPatronFirstLast(e.PatronIndex, events)
			firstTime := getEventIDTime(events, first)
			lastTime := getEventIDTime(events, last)

			durations[e.PatronIndex] = structs.Duration{
				PiSerial:  cfg.GetSerial(),
				SessionId: e.SessionID,
				FCFSSeqId: e.FCFSSeqID,
				DeviceTag: e.DeviceTag,
				PatronID:  e.PatronIndex,
				MfgID:     e.ManufacturerIndex,
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

func MultiDayDurations(swap bool, newPid int, events []structs.WifiEvent) (map[int]*structs.Duration, int) {
	// FIXME: note the swap of times when things are borked in the DB...
	// is that the best way to fix things?
	cfg := state.GetConfig()
	// We want, for every patron_id, to know when the device started/ended.
	checked := make(map[int]bool)
	durations := make(map[int]*structs.Duration)
	uniqsessions := make(map[string]string)

	// Get the largest patron ID in this set. Use it for new,
	// unique patron IDs.
	maxPatronID := -1
	for _, e := range events {
		if e.PatronIndex > maxPatronID {
			maxPatronID = e.PatronIndex
		}
	}
	maxPatronID += 1

	for _, e := range events {
		// For later
		uniqsessions[e.SessionID] = e.SessionID

		//log.Println("Patron index:", e.PatronIndex)
		if _, ok := checked[e.PatronIndex]; ok {
			// Skip if we already checked this patron
		} else {
			checked[e.PatronIndex] = true
			first, last := getPatronFirstLast(e.PatronIndex, events)
			firstTime := getEventIDTime(events, first)
			lastTime := getEventIDTime(events, last)
			if lastTime.Before(firstTime) {
				cfg.Log().Error("start", firstTime, "end", lastTime)
				cfg.Log().Error("cannot start after end! swapping...")
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
						durations[maxPatronID] = &structs.Duration{
							PiSerial:  cfg.GetSerial(),
							SessionID: e.SessionID,
							FCFSSeqID: e.FCFSSeqID,
							DeviceTag: e.DeviceTag,
							PatronID:  maxPatronID,
							MfgID:     e.ManufacturerIndex,
							Start:     firstTime.Format(time.RFC3339),
							End:       endOfToday.Format(time.RFC3339)}
						maxPatronID += 1
						firstTime = bod(tomorrow(firstTime))
						firstDay = firstTime.Day()

					}

					//duration := lastTime.Sub(firstTime)
					//log.Println("last split", int(float64(duration)/float64(time.Hour)), "hour device", e.PatronIndex, maxPatronId, firstDay, lastDay)
					endOfToday := lastTime
					firstTime = bod(lastTime)
					// When done looping, insert the remainder...
					durations[maxPatronID] = &structs.Duration{
						PiSerial:  cfg.GetSerial(),
						SessionID: e.SessionID,
						FCFSSeqID: e.FCFSSeqID,
						DeviceTag: e.DeviceTag,
						PatronID:  maxPatronID,
						MfgID:     e.ManufacturerIndex,
						Start:     firstTime.Format(time.RFC3339),
						End:       endOfToday.Format(time.RFC3339)}
					maxPatronID += 1

				} else {
					durations[e.PatronIndex] = &structs.Duration{
						PiSerial:  cfg.GetSerial(),
						SessionID: e.SessionID,
						FCFSSeqID: e.FCFSSeqID,
						DeviceTag: e.DeviceTag,
						PatronID:  e.PatronIndex,
						MfgID:     e.ManufacturerIndex,
						Start:     firstTime.Format(time.RFC3339),
						End:       lastTime.Format(time.RFC3339)}
				}
			}

		} // end else
	}

	// These are now in need of patron renumbering, because new patrons were introduced.
	// This leaves patron ID gaps. Lets keep those monotonic and gapless.
	pid := 0
	remapped := make(map[int]*structs.Duration)
	sorted := make([]*structs.Duration, 0)
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
	newmap := make(map[int]*structs.Duration)
	for _, v := range remapped {
		newv := structs.Duration{}
		newv.DeviceTag = v.DeviceTag
		newv.End = v.End
		newv.FCFSSeqID = v.FCFSSeqID
		newv.MfgID = v.MfgID
		//log.Println(v.PatronId, "becomes", pid)
		newv.PatronID = pid
		newv.PiSerial = v.PiSerial
		st, _ := time.Parse(time.RFC3339, v.Start)
		newv.SessionID = fmt.Sprintf("%v%02d%02d", st.Year(), st.Month(), st.Day())
		newv.Start = v.Start
		newmap[pid] = &newv
		pid = pid + 1
	}

	return newmap, pid
}

// Summarize returns the drawing context where the image is drawn. This can then
// be written to disk.
func Summarize(events []structs.WifiEvent) (c *Counter, d map[int]structs.Duration) {
	sort.Slice(events, func(i, j int) bool {
		return events[i].Id < events[j].Id
	})
	c = doCounting(events)
	d = durationSummary(events)
	return c, d
}
