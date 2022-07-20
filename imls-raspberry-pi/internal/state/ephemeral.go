package state

import (
	"crypto/sha256"
	"fmt"
	"time"
)

type StartEnd struct {
	Start int64
	End   int64
}

type EphemeralDB map[string]StartEnd

var ed EphemeralDB = make(EphemeralDB)

func GetMACs() EphemeralDB {
	return ed
}

func ClearEphemeralDB() {
	ed = make(EphemeralDB)
}

// NOTE: Do not log MAC addresses.
func RecordMAC(mac string) {
	now := GetClock().Now().In(time.Local).Unix()
	// cfg := GetConfig()
	// cfg.Log().Debug("THE TIME IS NOW ", GetClock().Now().In(time.Local), " or ", now)

	// Check if we already have the MAC address in the ephemeral table.
	if p, ok := ed[mac]; ok {
		//cfg.Log().Debug(mac, " exists, updating")
		// Has this device been away for more than 2 hours?
		// Start by grabbing the start/end times.
		se := ed[mac]
		if (now > se.End) && ((now - se.End) > MAC_MEMORY_DURATION_SEC) {
			// If it has been, we need to "forget" the old device.
			// Do this by hashing the mac with the current time, store the original data
			// unchanged, and create a new entry for the current mac address, in case we
			// see it again (in less than 2h).
			// cfg.Log().Debug(mac, " is an old mac, refreshing/changing")
			hash := sha256.Sum256([]byte(mac + fmt.Sprint(now)))
			ed[fmt.Sprintf("%x", hash)] = se
			ed[mac] = StartEnd{Start: now, End: now}
		} else {
			// Just update the mac address. It has been less than 2h.
			ed[mac] = StartEnd{Start: p.Start, End: now}
		}
	} else {
		// We have never seen the MAC address.
		//cfg.Log().Debug(mac, " is new, inserting")
		ed[mac] = StartEnd{Start: now, End: now}
	}
}
