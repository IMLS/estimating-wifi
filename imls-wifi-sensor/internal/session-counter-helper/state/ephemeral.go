package state

import (
	"crypto/sha1"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"gsa.gov/18f/internal/config"
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

	// Check if we already have the MAC address in the ephemeral table.
	if p, ok := ed[mac]; ok {
		// Has this device been away for more than 2 hours?
		// Start by grabbing the start/end times.
		se := ed[mac]

		// XXX Not for production. Stop logging MACs outside of testing
		log.Debug().Msg("Known MAC address--> " + mac)

		if (now > se.End) && ((now - se.End) > int64(config.GetDeviceMemory())) {
			// If it has been, we need to "forget" the old device.
			// Do this by hashing the mac with the current time, store the original data
			// unchanged, and create a new entry for the current mac address, in case we
			// see it again (in less than 2h).
			//log.Debug().Msg("Been away for 2h. Forgetting old MAC address.")
			sha1 := sha1.Sum([]byte(mac + fmt.Sprint(now)))
			ed[fmt.Sprintf("%x", sha1)] = se
			ed[mac] = StartEnd{Start: now, End: now}
		} else {
			// Just update the mac address. It has been less than 2h.
			//log.Debug().Msg("Updating end time for MAC address.")
			ed[mac] = StartEnd{Start: p.Start, End: now}
		}
	} else {
		// We have never seen the MAC address.

		// XXX Not for production. Stop logging MACs outside of testing
		log.Debug().Msg("New MAC address--> " + mac)

		ed[mac] = StartEnd{Start: now, End: now}
	}
}
