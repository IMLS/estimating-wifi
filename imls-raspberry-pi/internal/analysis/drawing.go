// Package analysis provides visualization primitives.
package analysis

import (
	"fmt"
	"sort"
	"time"

	"github.com/fogleman/gg"
	"gsa.gov/18f/internal/logwrapper"
	"gsa.gov/18f/internal/state"
	"gsa.gov/18f/internal/structs"
)

func DrawPatronSessionsFromWifi(events []structs.WifiEvent, outputPath string) {
	_, d := Summarize(events)
	durations := make([]structs.Duration, 0)
	for _, v := range d {
		durations = append(durations, v)
	}
	DrawPatronSessions(durations, outputPath)
}

func DrawPatronSessions(durations []structs.Duration, outputPath string) {
	cfg := state.GetConfig()
	lw := logwrapper.NewLogger(nil)

	if len(durations) == 0 {
		lw.Error("DrawPatronSessions was passed zero durations to draw.")
		lw.Error("Wanted to draw to the output path ", outputPath)
		return
	}

	// Capture the data about the session while running in a `counter` structure.
	durationsInRange := 0
	sort.Sort(structs.ByStart(durations))
	lw.Debug("about to iterate over [", len(durations), "] durations")

	for _, d := range durations {
		st, _ := time.Parse(time.RFC3339, d.Start)
		et, _ := time.Parse(time.RFC3339, d.End)
		diff := int(et.Sub(st).Minutes())
		// log.Println("st", st, "et", et, "diff", diff)
		if (diff > cfg.Monitoring.MinimumMinutes) && (diff < cfg.Monitoring.MaximumMinutes) {
			//log.Println("id", d.PatronId, "diff", diff)
			durationsInRange += 1
		}
	}

	lw.Debug("durationsInRange [", durationsInRange, "]")

	WIDTH := 1440
	hourWidth := WIDTH / 24

	HEIGHT := 24 * (durationsInRange + 2)

	lw.Debug("image dimensions (WxH) ", WIDTH, " x ", HEIGHT)

	dc := gg.NewContext(WIDTH, HEIGHT)
	dc.SetRGBA(0.5, 0.5, 0, 0.5)
	dc.SetLineWidth(1)
	ystep := 0

	totalMinutes := 0
	totalPatrons := 0
	dc.SetRGB(1, 1, 1)
	dc.Push()
	dc.DrawRectangle(0, 0, float64(WIDTH), float64(HEIGHT))
	dc.Fill()
	dc.Stroke()
	dc.Pop()

	for _, d := range durations {
		// lw.Debug("duration ", d)
		st, _ := time.Parse(time.RFC3339, d.Start)
		et, _ := time.Parse(time.RFC3339, d.End)
		diff := int(et.Sub(st).Minutes())

		totalPatrons += 1
		totalMinutes += diff

		if (diff > cfg.Monitoring.MinimumMinutes) && (diff < cfg.Monitoring.MaximumMinutes) {
			ystep += 1

			// Draw the hour lines
			for hour := 1; hour <= 24; hour++ {
				x := hourWidth * hour
				if hour == 12 {
					dc.SetRGBA(0.9, 0.1, 0.1, 0.2)
					dc.SetLineWidth(2)
					dc.DrawLine(float64(x), 0, float64(x), float64(HEIGHT))
					dc.DrawStringAnchored("noon", float64(x+10), float64(10), 0, 0)
				} else {
					dc.SetRGBA(0.9, 0.9, 0.9, 0.2)
					dc.SetLineWidth(0.5)
					dc.DrawLine(float64(x), 0, float64(x), float64(HEIGHT))
				}
				dc.Stroke()
			}

			// Draw the duration block
			// 1440 minutes in a day
			dc.SetRGB(0.7, 0.2, 0.2)
			dc.SetLineWidth(1)

			// Therefore...
			// log.Println("eod", eod(st))
			stInMinutes := 1440 - int(eod(st).Sub(st).Minutes())
			x := stInMinutes
			y := 20 + (ystep * 20)
			// log.Println("start time", st, "end time", et)
			// log.Println("rect", x, y, diff, 20)

			dc.DrawRectangle(float64(x), float64(y), float64(diff), 20)
			dc.Stroke()

			// Position the start time string
			dc.SetRGB(0.2, 0.2, 0.2)
			if st.Hour() < 1 {
				dc.DrawStringAnchored(fmt.Sprintf("%v:%v", st.Hour(), pad(st.Minute())), float64(x+diff), float64(y), -0.5, 1)
			} else {
				dc.DrawStringAnchored(fmt.Sprintf("%v:%v", st.Hour(), pad(st.Minute())), float64(x), float64(y), 1.1, 1)
			}

			// Position the duration string
			duration := ""
			if diff < 60 {
				duration = fmt.Sprintf("%vm", pad(diff))
			} else {
				// log.Println("diff", diff)
				hours := (diff / 60)
				minutes := diff - ((diff / 60) * 60)
				duration = fmt.Sprintf("%vh%vm", hours, pad(minutes))
				// log.Println(duration)
			}

			dc.DrawStringAnchored(duration, float64(x+diff), float64(y), 1.25, 1)

			dc.Stroke()
		}

	}

	day, _ := time.Parse(time.RFC3339, durations[0].Start)
	summaryD := fmt.Sprintf("Patron sessions from %v %v, %v - %v %v", day.Month(), day.Day(), day.Year(), cfg.GetFCFSSeqId(), cfg.GetDeviceTag())
	summaryA := fmt.Sprintf("%v devices seen", totalPatrons)
	summaryP := fmt.Sprintf("%v patron devices", durationsInRange)
	summaryM := fmt.Sprintf("%v minutes served", totalMinutes)

	dc.DrawStringAnchored(summaryD, float64(20), float64(20), 0, 0)

	firstLineY := float64(HEIGHT) - ((20 * 5) + 10)
	dc.DrawStringAnchored(summaryA, float64(20), float64(firstLineY+20), 0, 0)
	dc.DrawStringAnchored(summaryP, float64(20), float64(firstLineY+40), 0, 0)
	dc.DrawStringAnchored(summaryM, float64(20), float64(firstLineY+60), 0, 0)

	// LEGEND
	xpos := float64(WIDTH - 300)
	dc.SetRGB(0.9, 0.1, 0.1)
	dc.DrawRectangle(xpos, 7.5, 120, 20)
	dc.Stroke()
	dc.SetRGB(0.0, 0.0, 0.0)
	dc.DrawStringAnchored("LEGEND", xpos-100, 7.5, 1, 1)
	dc.DrawStringAnchored("LEGEND", xpos-99, 7.5, 1, 1)
	w, _ := dc.MeasureString("LEGEND")
	dc.DrawLine(xpos-100-w, 35, xpos+120, 35)
	dc.Stroke()

	dc.DrawStringAnchored("start time", xpos, 7.5, 1.15, 1)
	dc.DrawStringAnchored("duration", xpos, 7.5, -0.95, 1)

	// Hours
	dc.SetRGB(0.7, 0.7, 0.7)
	for hour := 1; hour <= 23; hour++ {
		x := float64(hourWidth * hour)
		dc.Push()
		//gg.Translate(x, float64(HEIGHT-20))
		///dc.Rotate(gg.Degrees(90))
		dc.DrawStringAnchored(fmt.Sprintf("%v:00", hour), x, float64(HEIGHT-10), 0.5, 0)
		dc.Pop()
	}

	//baseFilename := fmt.Sprint(filepath.Join(outdir, fmt.Sprintf("%v-%v-%v", sid, seqId, dt)))
	lw.Debug("writing summary image to ", outputPath)

	err := dc.SavePNG(outputPath)
	if err != nil {
		lw.Info("drawing: failed to save png")
		lw.Fatal(err.Error())
	}
}

func pad(n int) string {
	if n < 10 {
		return fmt.Sprintf("0%v", n)
	} else {
		return fmt.Sprint(n)
	}
}

func eod(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 23, 59, 59, 0, t.Location())
}
