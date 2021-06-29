package logwrapper

import (
	"fmt"
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

// Convoluted for use within libraries...
func NewLogger(cfg *config.Config) *StandardLogger {
	once.Do(func() {
		var baseLogger = logrus.New()

		// If we have a config object...
		if cfg != nil {
			// Set the output to a file if we have a config file to guide us.

			iow, err := os.OpenFile(cfg.Local.Logfile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
			if err != nil {
				fmt.Printf("could not open logfile %v for wriiting\n", cfg.Local.Logfile)
				os.Exit(-1)
			}
			baseLogger.SetOutput(iow)
		}

		// If we have a valid config file, and lw is not already configured...
		standardLogger = &StandardLogger{baseLogger}
		standardLogger.Formatter = &logrus.JSONFormatter{}

	})

	return standardLogger
}

const (
	INFO = iota
	WARN
	ERROR
	FATAL
)

// Declare variables to store log messages as new Events
var (
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
	case INFO:
		l.WithFields(fields).Info(fmt.Sprintf(e.message, args...))
	case WARN:
		l.WithFields(fields).Warn(fmt.Sprintf(e.message, args...))
	case ERROR:
		l.WithFields(fields).Error(fmt.Sprintf(e.message, args...))
	case FATAL:
		l.WithFields(fields).Fatal(fmt.Sprintf(e.message, args...))
	}

}

// InvalidArg is a standard error message
func (l *StandardLogger) Info(msg string, args ...interface{}) {
	_, file, _, _ := runtime.Caller(1)
	l.Base(infoMsg, filepath.Base(file), fmt.Sprintf(msg, args...))
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
