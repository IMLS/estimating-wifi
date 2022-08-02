package zls

import (
	"io"
	"os"
	"strings"
	"time"

	"github.com/buger/jsonparser"
	"github.com/getsentry/sentry-go"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type ZeroLogSentry struct {
	name   string
	client *sentry.Client
}

func getAppStacktrace() *sentry.Stacktrace {
	trace := sentry.NewStacktrace()
	appTraces := []sentry.Frame{}
	for count := len(trace.Frames) - 1; count >= 0; count-- {
		frame := trace.Frames[count]
		// do a rough filter pass on this frame
		if frame.InApp && (!strings.HasPrefix(frame.Module, "github.com")) {
			appTraces = append(appTraces, frame)
		}
	}
	trace.Frames = appTraces
	return trace
}

func (zls *ZeroLogSentry) Write(data []byte) (int, error) {
	event := sentry.Event{
		Logger: zls.name,
		Extra:  make(map[string]interface{}),
	}
	err := jsonparser.ObjectEach(data, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
		switch string(key) {
		case "level":
			switch string(value) {
			case "debug":
				event.Level = sentry.LevelDebug
			case "info":
				event.Level = sentry.LevelInfo
			case "warn":
				event.Level = sentry.LevelWarning
			case "error":
				event.Level = sentry.LevelError
			case "fatal", "panic":
				event.Level = sentry.LevelFatal
			}
		case "error":
			exc := sentry.Exception{
				Type:       string(value),
				Stacktrace: getAppStacktrace(),
			}
			event.Exception = []sentry.Exception{exc}
		case "time":
			ts, err := time.Parse(time.RFC3339, string(value))
			if err != nil {
				ts = time.Now()
			}
			event.Timestamp = ts
		case "message":
			event.Message = string(value)
		default:
			event.Extra[string(key)] = string(value)
		}
		return nil
	})
	if err == nil {
		scope := sentry.CurrentHub().Scope()
		zls.client.CaptureEvent(&event, nil, scope)
		defer zls.client.Flush(2 * time.Second)
	}
	return len(data), nil
}

func SetTags(m map[string]string) {
	sentry.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetTags(m)
	})
}

func SetupZeroLogSentry(name string, dsn string) {
	options := sentry.ClientOptions{
		Dsn:        dsn,
		SampleRate: 1.0, // no sampling; send all events
	}
	client, err := sentry.NewClient(options)
	if err != nil {
		log.Error().Err(err).Msg("could not initialize sentry")
	} else {
		customWriter := ZeroLogSentry{
			client: client,
			name:   name,
		}
		writer := io.MultiWriter(&customWriter, os.Stdout)
		log.Logger = zerolog.New(writer).With().Timestamp().Logger()
	}
}
