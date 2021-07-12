package tlp

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"gsa.gov/18f/config"
	"gsa.gov/18f/http"
)

func report(service string, cfg *config.Config, session_id int, arr []map[string]interface{}) (http_error_count int, err error) {
	http_error_count = 0

	uri := http.FormatUri(cfg.Umbrella.Scheme, cfg.Umbrella.Host, cfg.Umbrella.Data)
	err2 := http.PostJSON(cfg, uri, arr)
	if err2 != nil {
		log.Println("report2:", service, "results POST failure")
		log.Println(err2)
		http_error_count = http_error_count + 1
	}

	var resultErr error = nil
	if http_error_count > 0 {
		resultErr = fmt.Errorf("error count is now %d", http_error_count)
	}
	return http_error_count, resultErr
}

func StoreToCloud(ka *Keepalive, cfg *config.Config, ch_data <-chan []map[string]interface{}, ch_reset <-chan Ping, ch_kill <-chan Ping) {

	log.Println("Starting ReportOut")

	http_error_count := 0

	//ch_kill will be nil in production
	var ping, pong chan interface{} = nil, nil
	if ch_kill == nil {
		ping, pong = ka.Subscribe("ReportOut", 30)
	}

	// For event logging
	// el := http.NewEventLogger(cfg)

	// We never reset anything when storing to the cloud; that is only used by the SQLite version.
	// Spawn a concurrent process to consume everything that comes in on the reset channel.
	go Blackhole(ch_reset)

	for {
		select {
		case <-ping:
			// Every [cfg.Monitoring.HTTPErrorIntervalMins] this value
			// is reset to zero. If we get too many errors in that number of
			// minutes, then we should fail the next pong request. This will
			// kill the program, and we'll restart.
			if http_error_count < cfg.Monitoring.MaxHTTPErrorCount {
				pong <- "ReportOut"
			} else {
				log.Printf("reportout: http_error_count threshold of %d reached\n", http_error_count)
			}

		case <-ch_kill:
			log.Println("Exiting StoreToCloud")
			return

		// This is the [ uid -> ticks ] map (uid looks like "Next:0")
		case arr := <-ch_data:
			// TODO: event ids are broken and we need a better approach.
			event_ndx := 1
			// event_ndx, logerr := el.Log("logging_devices", nil)

			// Overwrite the existing event IDs in the prepared data.
			// We want it to connect to the event logged in the DB.
			for _, m := range arr {
				m["event_id"] = strconv.Itoa(event_ndx)
			}

			errCount, err := report("reval", cfg, event_ndx, arr)
			if err != nil {
				http_error_count += errCount
				log.Println("reportout: error in reporting to reval")
				log.Println(err)
				log.Println("reportout: HTTP_ERROR_COUNT", http_error_count)
			}

		case <-time.After(time.Duration(cfg.Monitoring.HTTPErrorIntervalMins) * time.Minute):
			// If this much time has gone by, go ahead and reset the error count.
			if http_error_count != 0 {
				log.Printf("reportout: RESETTING http_error_count FROM %d TO 0\n", http_error_count)
				http_error_count = 0
			}
		}
	}
}
