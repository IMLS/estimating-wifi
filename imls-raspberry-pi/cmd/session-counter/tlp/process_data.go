package tlp

import (
	"fmt"
	"log"

	"gsa.gov/18f/internal/interfaces"
	"gsa.gov/18f/internal/state"
	"gsa.gov/18f/internal/structs"
)

func ProcessData(db interfaces.Database, dDB interfaces.Database, sq *state.Queue, iq *state.Queue) bool {
	cfg := state.GetConfig()
	// Queue up what needs to be sent still.
	thissession := cfg.GetCurrentSessionID()
	cfg.Log().Debug("queueing current session [ ", thissession, " ] to images and send queue... ")
	if thissession >= 0 {
		sq.Enqueue(fmt.Sprint(thissession))
		iq.Enqueue(fmt.Sprint(thissession))
	}
	// Grab the ephemeral durations
	var eds []structs.EphemeralDuration
	db.GetPtr().Select(&eds, "SELECT * FROM ephemeraldurations")
	log.Println(eds)
	// Copy them over to the durations DB, with additional data
	pidCounter := 0

	for _, ed := range eds {

		d := structs.Duration{
			PiSerial:  cfg.GetSerial(),
			SessionID: fmt.Sprint(cfg.GetCurrentSessionID()),
			FCFSSeqID: cfg.GetFCFSSeqID(),
			DeviceTag: cfg.GetDeviceTag(),
			PatronID:  pidCounter,
			// FIXME: All times should become UNIX epoch seconds...
			Start: fmt.Sprint(ed.Start),
			End:   fmt.Sprint(ed.End)}

		dDB.GetTableFromStruct(structs.Duration{}).InsertStruct(d)
	}
	return true
}
