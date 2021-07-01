package logwrapper

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/sirupsen/logrus"
	"gsa.gov/18f/config"
)

// Code for this wrapper inspired by
// https://www.datadoghq.com/blog/go-logging/

// Event stores messages to log later, from our standard interface
type Event struct {
	id      int
	level   int
	message string
}

// StandardLogger enforces specific log message formats
type StandardLogger struct {
	*logrus.Logger
}

var standardLogger *StandardLogger = nil
var once sync.Once

const LOGDIR = "/var/log/session-counter"

var logLevel int = ERROR

func (l *StandardLogger) SetLogLevel(lvl int) {
	// Don't allow broken logging levels.
	if lvl > DEBUG && lvl <= FATAL {
		logLevel = lvl
	}
}

func NewLogger(cfg *config.Config) (sl *StandardLogger) {
	once.Do(func() {
		sl = newLogger(cfg)
	})
	return sl
}

func UnsafeNewLogger(cfg *config.Config) (sl *StandardLogger) {
	return newLogger(cfg)
}

// Convoluted for use within libraries...
func newLogger(cfg *config.Config) *StandardLogger {
	var baseLogger = logrus.New()

	// If we have a config file, grab the loggers defined there.
	// Otherwise, use stderr.
	loggers := make([]string, 0)
	if cfg == nil {
		loggers = append(loggers, "local:stderr")
	} else {
		log.Println("cfg not nil")
		loggers = cfg.GetLoggers()
	}
	log.Println("loggers", loggers)
	writers := make([]io.Writer, 0)

	for _, l := range loggers {
		switch l {
		case "local:stderr":
			writers = append(writers, os.Stderr)
		case "local:file":
			os.Mkdir(LOGDIR, 0755)
			filename := filepath.Join(LOGDIR, "log.json")
			iow, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
			if err != nil {
				fmt.Printf("could not open logfile %v for writing\n", filename)
				os.Exit(-1)
			}
			writers = append(writers, iow)
		case "local:tmp":
			filename := filepath.Join("/tmp", "log.json")
			iow, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
			if err != nil {
				fmt.Printf("could not open logfile %v for writing\n", filename)
				os.Exit(-1)
			}
			writers = append(writers, iow)
		case "api:directus":
			uri := cfg.GetLoggingUri()
			api := NewApiLogger(uri)
			writers = append(writers, api)
		}
	}

	mw := io.MultiWriter(writers...)
	baseLogger.SetOutput(mw)

	// If we have a valid config file, and lw is not already configured...
	standardLogger = &StandardLogger{baseLogger}
	standardLogger.Formatter = &logrus.JSONFormatter{}

	return standardLogger
}

const (
	DEBUG = iota
	INFO
	WARN
	ERROR
	FATAL
)

// Declare variables to store log messages as new Events
// THe error level recorded here will impact what happens at runtime.
// For example, a FATAL message type will exit.
var (
	debugMsg  = Event{0, INFO, "%s"}
	infoMsg   = Event{1, INFO, "%s"}
	warnMsg   = Event{2, WARN, "%s"}
	errorMsg  = Event{3, ERROR, "%s"}
	fatalMsg  = Event{4, FATAL, "%s"}
	notFound  = Event{5, FATAL, "not found: %s"}
	lengthMsg = Event{6, INFO, "array %v length %d"}
)

func (l *StandardLogger) Base(e Event, loc string, args ...interface{}) {
	fields := logrus.Fields{
		"file": loc,
	}
	switch e.level {
	case DEBUG:
		if logLevel <= DEBUG {
			l.WithFields(fields).Debug(fmt.Sprintf(e.message, args...))
		}
	case INFO:
		if logLevel <= INFO {
			l.WithFields(fields).Info(fmt.Sprintf(e.message, args...))
		}
	case WARN:
		if logLevel <= WARN {
			l.WithFields(fields).Warn(fmt.Sprintf(e.message, args...))
		}
	case ERROR:
		// Always log ERROR level log events.
		l.WithFields(fields).Error(fmt.Sprintf(e.message, args...))
	case FATAL:
		// Always log FATAL level log events.
		l.WithFields(fields).Fatal(fmt.Sprintf(e.message, args...))
		// We're leaving, on a jet plane...
		os.Exit(-1)
	}

}

func (l *StandardLogger) Debug(msg string, args ...interface{}) {
	_, file, _, _ := runtime.Caller(1)
	l.Base(debugMsg, filepath.Base(file), fmt.Sprintf(msg, args...))
}

// InvalidArg is a standard error message
func (l *StandardLogger) Info(msg string, args ...interface{}) {
	_, file, _, _ := runtime.Caller(1)
	l.Base(infoMsg, filepath.Base(file), fmt.Sprintf(msg, args...))
}

func (l *StandardLogger) Warn(msg string, args ...interface{}) {
	_, file, _, _ := runtime.Caller(1)
	l.Base(warnMsg, filepath.Base(file), fmt.Sprintf(msg, args...))
}

func (l *StandardLogger) Error(msg string, args ...interface{}) {
	_, file, _, _ := runtime.Caller(1)
	l.Base(errorMsg, filepath.Base(file), fmt.Sprintf(msg, args...))
}

func (l *StandardLogger) Fatal(msg string, args ...interface{}) {
	_, file, _, _ := runtime.Caller(1)
	l.Base(fatalMsg, filepath.Base(file), fmt.Sprintf(msg, args...))
}

func (l *StandardLogger) ExeNotFound(path string) {
	_, file, _, _ := runtime.Caller(1)
	l.Base(notFound, filepath.Base(file), path)
}

func (l *StandardLogger) Length(arrname string, arr ...interface{}) {
	_, file, _, _ := runtime.Caller(1)
	l.Base(lengthMsg, filepath.Base(file), arrname, len(arr))
}
