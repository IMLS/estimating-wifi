package tlp

import (
	"github.com/rs/zerolog/log"
	"gsa.gov/18f/cmd/session-counter/state"
	"gsa.gov/18f/cmd/session-counter/structs"
	"gsa.gov/18f/internal/config"
)

// https://stackoverflow.com/questions/71274361/go-error-cannot-use-generic-type-without-instantiation
// Instantiate generics.
func ProcessData(dDB *state.DurationsDB, sq *state.Queue[int64]) bool {
	// Queue up what needs to be sent still.
	thissession := state.GetCurrentSessionID()

	log.Debug().
		Int64("session", thissession).
		Msg("queueing to images and send")

	if thissession >= 0 {
		sq.Enqueue(thissession)
	}

	pidCounter := 0
	durations := make([]*structs.Duration, 0)

	for _, se := range state.GetMACs() {

		d := &structs.Duration{
			PiSerial:  state.GetCachedSerial(),
			SessionID: state.GetCurrentSessionID(),
			FCFSSeqID: config.GetFCFSSeqID(),
			DeviceTag: config.GetDeviceTag(),
			PatronID:  pidCounter,
			Start:     se.Start,
			End:       se.End,
		}

		durations = append(durations, d)
		pidCounter += 1
	}

	dDB.InsertMany(durations)
	return true
}
