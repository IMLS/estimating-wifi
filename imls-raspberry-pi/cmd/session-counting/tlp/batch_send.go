package tlp

import (
	"log"

	"gsa.gov/18f/analysis"
	"gsa.gov/18f/config"
	"gsa.gov/18f/http"
	"gsa.gov/18f/logwrapper"
	"gsa.gov/18f/session-counter/model"
)

func BatchSend(ka *Keepalive, cfg *config.Config, kb *KillBroker,
	ch_durations_db_in <-chan *model.TempDB) {

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
			sq := model.NewQueue(cfg, "sent")
			nextSessionIdToSend := sq.Peek()

			for nextSessionIdToSend != nil {
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
					sq.Dequeue()
				} else if cfg.IsStoringToApi() {
					lw.Info("attempting to send batch [", nextSessionIdToSend, "][", len(durations), "] to the API server")
					// convert []Duration to an array of map[string]interface{}
					data := make([]map[string]interface{}, 0)
					for _, duration := range durations {
						data = append(data, duration.AsMap())
					}
					// After writing images, we come back and try and send the data remotely.
					lw.Debug("PostJSONing ", len(data), " duration datas")
					_, err = http.PostJSON(cfg, cfg.GetDataUri(), data)
					if err != nil {
						log.Println("could not log to API")
						log.Println(err.Error())
					} else {
						// If we successfully sent the data remotely, we can now mark it is as sent.
						sq.Dequeue()
					}
				} else {
					// Always dequeue. We're storing locally "for free" into the
					// durations table before trying to do the send.
					sq.Dequeue()
				}
				// See if we have something else to send...
				nextSessionIdToSend = sq.Peek()
			}

		}
	}
}
