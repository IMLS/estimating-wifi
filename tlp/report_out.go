package tlp

import (
	"log"
	"math/rand"
	"time"

	"gsa.gov/18f/session-counter/api"
	"gsa.gov/18f/session-counter/config"
	"gsa.gov/18f/session-counter/csp"
)

func reportToDirectus(cfg *config.Config, h map[string]int) (http_error_count int, err error) {
	http_error_count = 0

	svr := config.GetServer(cfg, "directus")
	tok, errGT := api.GetToken(svr)
	if errGT != nil {
		log.Println("report: error in token fetch")
		log.Println(errGT)
		http_error_count = http_error_count + 1
	} else {
		// If we had no problems getting a token, we can then report
		// the data to Directus.
		for uid, count := range h {
			go func(id string, cnt int) {
				time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
				err := api.StoreDeviceCount(svr, tok, id, cnt)
				if err != nil {
					log.Println("report: results POST failure")
					log.Println(err)
					http_error_count = http_error_count + 1
				}
			}(uid, count)
		}
	}

	return http_error_count, nil
}

func reportToReval(cfg *config.Config, h map[string]int) (http_error_count int, err error) {
	http_error_count = 0

	svr := config.GetServer(cfg, "reval")
	tok, errGT := api.GetToken(svr)
	if errGT != nil {
		log.Println("report: error in token fetch")
		log.Println(errGT)
		http_error_count = http_error_count + 1
	} else {
		// If we had no problems getting a token, we can then report
		// the data to Directus.
		for uid, count := range h {
			go func(id string, cnt int) {
				time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
				err := api.StoreDeviceCount(svr, tok, id, cnt)
				if err != nil {
					log.Println("report: results POST failure")
					log.Println(err)
					http_error_count = http_error_count + 1
				}
			}(uid, count)
		}
	}

	return http_error_count, nil
}

func ReportOut(ka *csp.Keepalive, cfg *config.Config, ch_uidmap <-chan map[string]int) {
	log.Println("Starting ReportOut")
	ping, pong := ka.Subscribe("ReportOut", 15)

	http_error_count := 0

	for {
		select {
		case <-ping:
			// Every [cfg.Monitoring.HTTPErrorIntervalMins] this value
			// is reset to zero. If we get too many errors in that number of
			// minutes, then we should fail the next pong request. This will
			// kill the program, and we'll restart.
			if http_error_count < cfg.Monitoring.MaxHTTPErrorCount {
				pong <- "reportMap"
			} else {
				log.Printf("report: http_error_count threshold of %d reached\n", http_error_count)
			}
		// This is the [ uid -> ticks ] map (uid looks like "Next:0")
		case h := <-ch_uidmap:
			// Send the data to directus
			derrCount, derr := reportToDirectus(cfg, h)
			if derr != nil {
				log.Println("report: error in reporting to directus")
				log.Println(derr)
				http_error_count += derrCount
			}
			rerrCount, rerr := reportToReval(cfg, h)
			if rerr != nil {
				log.Println("report: error in reporting to directus")
				log.Println(derr)
				http_error_count += rerrCount
			}

		case <-time.After(time.Duration(cfg.Monitoring.HTTPErrorIntervalMins) * time.Minute):
			// If this much time has gone by, go ahead and reset the error count.
			if http_error_count != 0 {
				log.Printf("report: resetting http_error_count from %d to 0\n", http_error_count)
				http_error_count = 0
			}
		}
	}
}
