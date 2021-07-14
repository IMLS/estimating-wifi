package tlp

import (
	"log"

	"gsa.gov/18f/internal/analysis"
	"gsa.gov/18f/internal/config"
	"gsa.gov/18f/internal/http"
	"gsa.gov/18f/internal/logwrapper"
	"gsa.gov/18f/internal/tempdb"
)

func BatchSend(ka *Keepalive, cfg *config.Config, kb *KillBroker,
	ch_durations_db_in <-chan *tempdb.TempDB) {

	lw := logwrapper.NewLogger(nil)
	lw.Debug("Starting BatchSend")
	var ping, pong chan interface{} = nil, nil
	var ch_kill chan interface{} = nil
	if kb != nil {
		ch_kill = kb.Subscribe()
	} else {
		ping, pong = ka.Subscribe("GenerateDurations", 30)
	}

	for {
		select {
		case <-ping:
			pong <- "BatchSend"
		case <-ch_kill:
			lw.Debug("exiting BatchSend")
			return
		case db := <-ch_durations_db_in:
			// This only comes in on reset...
			sq := tempdb.NewQueue(cfg, "sent")
			sessionsToSend := sq.AsList()

			for _, nextSessionIdToSend := range sessionsToSend {
				durations := []analysis.Duration{}
				// lw.Debug("looking for session ", unsent.Session, " in durations table")
				db.Open()
				err := db.Ptr.Select(&durations, "SELECT * FROM durations WHERE session_id=?", nextSessionIdToSend)
				db.Close()

				if err != nil {
					lw.Info("error in extracting durations for session", nextSessionIdToSend)
					lw.Error(err.Error())
				}
				lw.Debug("found ", len(durations), " durations to send.")

				if len(durations) == 0 {
					lw.Info("found zero durations to send/draw. dequeing session [", nextSessionIdToSend, "]")
					sq.Remove(nextSessionIdToSend)
				} else if cfg.IsStoringToApi() {
					lw.Info("attempting to send batch [", nextSessionIdToSend, "][", len(durations), "] to the API server")
					// convert []Duration to an array of map[string]interface{}
					data := make([]map[string]interface{}, 0)
					for _, duration := range durations {
						data = append(data, duration.AsMap())
					}
					// After writing images, we come back and try and send the data remotely.
					lw.Debug("PostJSONing ", len(data), " duration datas")
					err = http.PostJSON(cfg, cfg.GetDataUri(), data)
					if err != nil {
						log.Println("could not log to API")
						log.Println(err.Error())
					} else {
						// If we successfully sent the data remotely, we can now mark it is as sent.
						sq.Remove(nextSessionIdToSend)
					}
				} else {
					// Always dequeue. We're storing locally "for free" into the
					// durations table before trying to do the send.
					sq.Remove(nextSessionIdToSend)
				}
			}

		}
	}
}
