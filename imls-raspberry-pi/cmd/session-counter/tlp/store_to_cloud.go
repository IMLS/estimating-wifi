package tlp

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"gsa.gov/18f/internal/config"
	"gsa.gov/18f/internal/http"
)

func report(service string, cfg *config.Config, sessionID int, arr []map[string]interface{}) (httpErrorCount int, err error) {
	httpErrorCount = 0

	uri := http.FormatUri(cfg.Umbrella.Scheme, cfg.Umbrella.Host, cfg.Umbrella.Data)
	err2 := http.PostJSON(cfg, uri, arr)
	if err2 != nil {
		log.Println("report2:", service, "results POST failure")
		log.Println(err2)
		httpErrorCount = httpErrorCount + 1
	}

	var resultErr error = nil
	if httpErrorCount > 0 {
		resultErr = fmt.Errorf("error count is now %d", httpErrorCount)
	}
	return httpErrorCount, resultErr
}

func StoreToCloud(ka *Keepalive, cfg *config.Config, chData <-chan []map[string]interface{}, chReset <-chan Ping, chKill <-chan Ping) {

	log.Println("Starting ReportOut")

	httpErrorCount := 0

	//chKill will be nil in production
	var ping, pong chan interface{} = nil, nil
	if chKill == nil {
		ping, pong = ka.Subscribe("ReportOut", 30)
	}

	// For event logging
	// el := http.NewEventLogger(cfg)

	// We never reset anything when storing to the cloud; that is only used by the SQLite version.
	// Spawn a concurrent process to consume everything that comes in on the reset channel.
	go Blackhole(chReset)

	for {
		select {
		case <-ping:
			// Every [cfg.Monitoring.HTTPErrorIntervalMins] this value
			// is reset to zero. If we get too many errors in that number of
			// minutes, then we should fail the next pong request. This will
			// kill the program, and we'll restart.
			if httpErrorCount < cfg.Monitoring.MaxHTTPErrorCount {
				pong <- "ReportOut"
			} else {
				log.Printf("reportout: http_error_count threshold of %d reached\n", httpErrorCount)
			}

		case <-chKill:
			log.Println("Exiting StoreToCloud")
			return

		// This is the [ uid -> ticks ] map (uid looks like "Next:0")
		case arr := <-chData:
			// TODO: event ids are broken and we need a better approach.
			eventNdx := 1
			// event_ndx, logerr := el.Log("logging_devices", nil)

			// Overwrite the existing event IDs in the prepared data.
			// We want it to connect to the event logged in the DB.
			for _, m := range arr {
				m["event_id"] = strconv.Itoa(eventNdx)
			}

			errCount, err := report("reval", cfg, eventNdx, arr)
			if err != nil {
				httpErrorCount += errCount
				log.Println("reportout: error in reporting to reval")
				log.Println(err)
				log.Println("reportout: HTTP_ERROR_COUNT", httpErrorCount)
			}

		case <-cfg.Clock.After(time.Duration(cfg.Monitoring.HTTPErrorIntervalMins) * time.Minute):
			// If this much time has gone by, go ahead and reset the error count.
			if httpErrorCount != 0 {
				log.Printf("reportout: RESETTING http_error_count FROM %d TO 0\n", httpErrorCount)
				httpErrorCount = 0
			}
		}
	}
}
