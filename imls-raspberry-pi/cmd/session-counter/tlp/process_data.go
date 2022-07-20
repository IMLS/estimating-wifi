package tlp

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"gsa.gov/18f/cmd/session-counter/state"
	"gsa.gov/18f/cmd/session-counter/structs"
	"gsa.gov/18f/internal/interfaces"
)

// https://stackoverflow.com/questions/71274361/go-error-cannot-use-generic-type-without-instantiation
// Instantiate generics.
func ProcessData(dDB interfaces.Database, sq *state.Queue[int64]) bool {
	// Queue up what needs to be sent still.
	thissession := state.GetCurrentSessionID()

	log.Debug().
		Int64("session", thissession).
		Msg("queueing to images and send")

	if thissession >= 0 {
		sq.Enqueue(thissession)
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
			// FIXME: All times should become UNIX epoch seconds...
			Start: se.Start,
			End:   se.End}

		//dDB.GetTableFromStruct(structs.Duration{}).InsertStruct(d)
		durations = append(durations, d)
		pidCounter += 1
	}

	dDB.GetTableFromStruct(structs.Duration{}).InsertMany(durations)
	return true
}
