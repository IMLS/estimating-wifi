package tlp

import (
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"gsa.gov/18f/internal/interfaces"
	"gsa.gov/18f/internal/state"
	"gsa.gov/18f/internal/structs"
)

func ProcessData(dDB interfaces.Database, sq *state.Queue, iq *state.Queue) bool {
	// Queue up what needs to be sent still.
	thissession := state.GetCurrentSessionID()

	log.Debug().
		Int64("session", thissession).
		Msg("queueing to images and send")

	if thissession >= 0 {
		sq.Enqueue(fmt.Sprint(thissession))
		iq.Enqueue(fmt.Sprint(thissession))
	}

	pidCounter := 0
	durations := make([]interface{}, 0)

	for _, se := range state.GetMACs() {

		d := structs.Duration{
			PiSerial:  state.GetSerial(),
			SessionID: fmt.Sprint(state.GetCurrentSessionID()),
			FCFSSeqID: state.GetFCFSSeqID(),
			DeviceTag: state.GetDeviceTag(),
			PatronID:  pidCounter,
			Start:     time.Unix(se.Start, 0).Format(time.RFC3339),
			End:       time.Unix(se.End, 0).Format(time.RFC3339),
		}

		//dDB.GetTableFromStruct(structs.Duration{}).InsertStruct(d)
		durations = append(durations, d)
		pidCounter += 1
	}

	dDB.GetTableFromStruct(structs.Duration{}).InsertMany(durations)
	return true
}
