package state

type StartEnd struct {
	Start int64
	End   int64
}

type EphemeralDB map[string]StartEnd

var ed EphemeralDB = make(EphemeralDB)

// NOTE: Do not log MAC addresses.
func RecordMAC(mac string) {
	now := GetClock().Now().Unix()
	// If we already have it
	if p, ok := ed[mac]; ok {
		//cfg.Log().Debug(mac, " exists, updating")
		ed[mac] = StartEnd{Start: p.Start, End: now}
	} else {
		//cfg.Log().Debug(mac, " is new, inserting")
		ed[mac] = StartEnd{Start: now, End: now}
	}
}
func GetMACs() EphemeralDB {
	return ed
}

func ClearEphemeralDB() {
	ed = make(EphemeralDB)
}
