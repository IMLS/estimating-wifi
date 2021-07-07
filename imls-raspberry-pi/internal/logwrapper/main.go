package logwrapper

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
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
	logger *logrus.Logger
}

var standardLogger *StandardLogger = nil
var once sync.Once

const LOGDIR = "/var/log/session-counter"

var logLevel int = ERROR

func (l *StandardLogger) SetLogLevel(level string) {
	switch strings.ToLower(level) {
	case "debug":
		logLevel = DEBUG
	case "info":
		logLevel = INFO
	case "warn":
		logLevel = WARN
	default:
		logLevel = ERROR
	}
}

func (l *StandardLogger) GetLogLevelName() string {
	switch logLevel {
	case 0:
		return "DEBUG"
	case 1:
		return "INFO"
	case 2:
		return "WARN"
	case 3:
		return "ERROR"
	case 4:
		return "FAIL"
	}
	return "UNKNOWN"
}

func NewLogger(cfg *config.Config) *StandardLogger {
	once.Do(func() {
		standardLogger = newLogger(cfg)
	})

	if standardLogger != nil {
		return standardLogger
	}
	if cfg == nil && standardLogger != nil {
		return standardLogger
	} else {
		log.Println("Falling back on UnsafeNewLogger...")
		return UnsafeNewLogger(cfg)
	}

}

// For unit testing only
func UnsafeNewLogger(cfg *config.Config) (sl *StandardLogger) {
	sl = standardLogger
	if sl != nil {
		return sl
	} else {
		sl = newLogger(cfg)
	}
	return sl
}

// Convoluted for use within libraries...
func newLogger(cfg *config.Config) *StandardLogger {
	var baseLogger = logrus.New()
	if baseLogger == nil {
		log.Println("baseLogger is nil")
	}
	// If we have a config file, grab the loggers defined there.
	// Otherwise, use stderr.
	loggers := make([]string, 0)
	level := "ERROR"
	if cfg == nil {
		loggers = append(loggers, "local:stderr")
	} else {
		loggers = cfg.GetLoggers()
		level = cfg.GetLogLevel()
	}

	// log.Println("loggers", loggers)
	writers := make([]io.Writer, 0)

	for _, l := range loggers {
		switch l {
		case "local:stderr":
			writers = append(writers, os.Stderr)
		case "local:stdout":
			writers = append(writers, os.Stdout)
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
			api := NewApiLogger(cfg)
			writers = append(writers, api)
		}
	}

	mw := io.MultiWriter(writers...)
	baseLogger.SetOutput(mw)

	// If we have a valid config file, and lw is not already configured...
	standardLogger = &StandardLogger{baseLogger}
	standardLogger.logger.Formatter = &logrus.JSONFormatter{}
	standardLogger.SetLogLevel(level)

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
	debugMsg  = Event{0, DEBUG, "%s"}
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
	if l == nil {
		log.Println("logger `l` not initialized.")
		log.Println("NewLogger() does not work in stateless contexts.")
		log.Println("You may be looking for UnsafeNewLogger().")
		log.Println("Creating a new (default) logger for you...")
		l = newLogger(nil)
	}
	switch e.level {
	case DEBUG:
		if logLevel >= DEBUG {
			l.logger.WithFields(fields).Debugf(e.message, args...)
		}
	case INFO:
		if logLevel >= INFO {
			l.logger.WithFields(fields).Infof(e.message, args...)
		}
	case WARN:
		if logLevel >= WARN {
			l.logger.WithFields(fields).Warnf(e.message, args...)
		}
	case ERROR:
		// Always log ERROR level log events.
		l.logger.WithFields(fields).Errorf(e.message, args...)
	case FATAL:
		// Always log FATAL level log events.
		l.logger.WithFields(fields).Fatalf(e.message, args...)
		// We're leaving, on a jet plane...
		os.Exit(-1)
	}
}

func (l *StandardLogger) Debug(msg string, args ...interface{}) {
	_, file, _, _ := runtime.Caller(1)
	if l == nil {
		log.Println("nil here", msg)
	}
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
