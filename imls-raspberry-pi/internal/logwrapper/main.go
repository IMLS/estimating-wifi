package logwrapper

import (
	"fmt"
	"os"
	"runtime"

	"github.com/sirupsen/logrus"
	"gsa.gov/18f/config"
)

// Code for this wrapper inspired by
// https://www.datadoghq.com/blog/go-logging/

// Event stores messages to log later, from our standard interface
type Event struct {
	id      int
	message string
}

// StandardLogger enforces specific log message formats
type StandardLogger struct {
	*logrus.Logger
}

// NewLogger initializes the standard logger
func NewLogger(cfg *config.Config) *StandardLogger {
	var baseLogger = logrus.New()
	var standardLogger = &StandardLogger{baseLogger}
	standardLogger.Formatter = &logrus.JSONFormatter{}

	iow, err := os.OpenFile(cfg.Local.Logfile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Printf("could not open logfile %v for wriiting\n", cfg.Local.Logfile)
		os.Exit(-1)
	}
	defer iow.Close()
	logrus.SetOutput(iow)

	return standardLogger
}

// Declare variables to store log messages as new Events
var (
	infoMsg = Event{1, "INFO: %s"}
)

func (l *StandardLogger) Base(e Event, loc string, args ...interface{}) {
	l.Errorf(fmt.Sprintf("[%v] %v", loc, fmt.Sprintf(e.message, args...)))
}

// InvalidArg is a standard error message
func (l *StandardLogger) Info(msg string) {
	_, file, _, _ := runtime.Caller(1)
	l.Base(infoMsg, file, msg)
}
