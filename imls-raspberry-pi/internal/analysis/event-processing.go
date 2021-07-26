package analysis

import (
	"fmt"
	"sort"
	"time"

	"gsa.gov/18f/internal/structs"
)

func GetSessions(events []structs.WifiEvent) []string {
	eventset := make([]string, 0)
	for _, e := range events {
		found := false
		for _, uniq := range eventset {
			if e.SessionID == uniq {
				found = true
			}
		}
		if !found {
			eventset = append(eventset, e.SessionID)
		}
	}
	return eventset
}

func RemapEvents(events []structs.WifiEvent) []structs.WifiEvent {
	// Some sessions start and end on the same day. Because we're rewriting the session id,
	// this means that it is possible to see a mapping get reused. The way to fix that is
	// to clear the mapping tables, but keep the counters going up. That way, sessions that start/end
	// on the same day will have unique devices.

	manufacturerNdx := 0
	patronNdx := 0
	// log.Println("len(events)", len(events))

	for _, pass := range []string{"first", "second"} {
		// Get the unique sessions in the dataset
		sessions := GetSessions(events)
		// log.Println("pass", pass, "sessions", sessions)
		// We need things in order. This matters for remapping
		// the manufacturer and patron indicies.
		sort.Slice(events, func(i, j int) bool {
			// return events[i].Localtime.Before(events[j].Localtime)
			return events[i].ID < events[j].ID
		})

		manufacturerNdx = 0
		patronNdx = 0
		var manufacturerMap map[int]int
		var patronMap map[int]int

		for _, s := range sessions {
			// For each session, create a new patron/mfg mapping.
			manufacturerMap = make(map[int]int)
			patronMap = make(map[int]int)
			// The second time through, we will have everything sessioned into days.
			// But, we have some big indicies. Reset them so they're small.
			if pass == "second" {
				// log.Println(s, "resetting map counters")
				manufacturerNdx = 0
				patronNdx = 0
			}
			// Now, remap every event.
			// This means putting it in a new session (based on days instead of event ids)
			// and renumbering the patron/
			for ndx, e := range events {
				if e.SessionID == s {
					// We will rewrite all session fields to the current day.
					// Need to modify the array, not the local variable `e`
					if pass == "first" {
						t, _ := time.Parse(time.RFC3339, e.Localtime)
						events[ndx].SessionID = fmt.Sprintf("%v%02d%02d", t.Year(), t.Month(), t.Day())
					}

					// If we have already mapped this manufacturer, then update
					// the event with the stored value.
					if val, ok := manufacturerMap[e.ManufacturerIndex]; ok {
						events[ndx].ManufacturerIndex = val
					} else {
						// Otherwise, update the map and the event.
						manufacturerMap[e.ManufacturerIndex] = manufacturerNdx
						events[ndx].ManufacturerIndex = manufacturerNdx
						manufacturerNdx += 1
					}

					// Now, check the patron index.
					if val, ok := patronMap[e.PatronIndex]; ok {
						events[ndx].PatronIndex = val
					} else {
						// fmt.Printf("[%v] Remapping %v to %v\n", s, e.PatronIndex, patronNdx)
						patronMap[e.PatronIndex] = patronNdx
						events[ndx].PatronIndex = patronNdx
						patronNdx += 1
					}
				} // if e.sessionId == s
			} // end for e := events

			// log.Println("max patronIndex", patronNdx, "max MfgIndex", manufacturerNdx)
		} // end for s := sessions
	}

	return events
}

func GetEventsFromSession(events []structs.WifiEvent, session string) []structs.WifiEvent {
	filtered := make([]structs.WifiEvent, 0)
	for _, e := range events {
		if e.SessionID == session {
			filtered = append(filtered, e)
		}
	}
	return filtered
}
