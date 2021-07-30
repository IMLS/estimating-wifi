package tlp

import (
	"fmt"

	"gsa.gov/18f/internal/interfaces"
	"gsa.gov/18f/internal/state"
	"gsa.gov/18f/internal/structs"
)

func ProcessData(dDB interfaces.Database, sq *state.Queue, iq *state.Queue) bool {
	cfg := state.GetConfig()
	// Queue up what needs to be sent still.
	thissession := cfg.GetCurrentSessionID()
	cfg.Log().Debug("queueing current session [ ", thissession, " ] to images and send queue... ")
	if thissession >= 0 {
		sq.Enqueue(fmt.Sprint(thissession))
		iq.Enqueue(fmt.Sprint(thissession))
	}

	pidCounter := 0
	durations := make([]interface{}, 0)

	for _, se := range state.GetMACs() {

		d := structs.Duration{
			PiSerial:  cfg.GetSerial(),
			SessionID: fmt.Sprint(cfg.GetCurrentSessionID()),
			FCFSSeqID: cfg.GetFCFSSeqID(),
			DeviceTag: cfg.GetDeviceTag(),
			PatronID:  pidCounter,
			// FIXME: All times should become UNIX epoch seconds...
			Start: fmt.Sprint(se.Start),
			End:   fmt.Sprint(se.End)}

		//dDB.GetTableFromStruct(structs.Duration{}).InsertStruct(d)
		durations = append(durations, d)
		pidCounter += 1
	}

	dDB.GetTableFromStruct(structs.Duration{}).InsertMany(durations)
	return true
}
