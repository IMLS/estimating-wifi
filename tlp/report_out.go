package tlp

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"gsa.gov/18f/session-counter/api"
	"gsa.gov/18f/session-counter/config"
)

func report(service string, cfg *config.Config, session_id int, h map[string]int) (http_error_count int, err error) {
	tok, errGT := config.ReadAuth()
	http_error_count = 0

	if errGT != nil {
		log.Println("report:", service, "error in token fetch")
		log.Println(errGT)
		http_error_count = http_error_count + 1
	} else {
		for uid, last_seen := range h {
			// Only report devices that we saw *this past minute*
			// That means we saw them `zero minutes ago`.
			if last_seen == 0 {
				go func(id string) {
					time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
					err := api.StoreDeviceCount(cfg, tok, session_id, id)
					if err != nil {
						log.Println("report:", service, "results POST failure")
						log.Println(err)
						http_error_count = http_error_count + 1
					}
				}(uid)
			}
		}
	}

	var resultErr error = nil
	if http_error_count > 0 {
		resultErr = fmt.Errorf("error count is now %d", http_error_count)
	}
	return http_error_count, resultErr
}

func report2(service string, cfg *config.Config, session_id int, h map[string]int) (http_error_count int, err error) {
	tok, errGT := config.ReadAuth()
	http_error_count = 0

	if errGT != nil {
		log.Println("report2:", service, "error in token fetch")
		log.Println(errGT)
		http_error_count = http_error_count + 1
	} else {
		err := api.StoreDevicesCount(cfg, tok, session_id, h)
		if err != nil {
			log.Println("report2:", service, "results POST failure")
			log.Println(err)
			http_error_count = http_error_count + 1
		}
	}

	var resultErr error = nil
	if http_error_count > 0 {
		resultErr = fmt.Errorf("error count is now %d", http_error_count)
	}
	return http_error_count, resultErr
}

func ReportOut(ka *Keepalive, cfg *config.Config, ch_uidmap <-chan map[string]int) {
	log.Println("Starting ReportOut")
	ping, pong := ka.Subscribe("ReportOut", 30)
	http_error_count := 0

	// For event logging
	el := api.NewEventLogger(cfg)

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
		// This is the [ uid -> ticks ] map (uid looks like "Next:0")
		case h := <-ch_uidmap:
			event_ndx, logerr := el.Log("logging_devices", nil)
			if logerr != nil {
				log.Println("reportout: error in event logging: ", logerr)
				http_error_count += 1
				log.Println("reportout: HTTP_ERROR_COUNT", http_error_count)
			}

			// This used to loop over "directus" and "reval"
			// We decided we will log only to reval, and it will handle validation and logging.
			errCount, err := report2("reval", cfg, event_ndx, h)
			if err != nil {
				log.Println("reportout: error in reporting to reval")
				log.Println(err)
				http_error_count += errCount
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
