package logwrapper

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
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
	*logrus.Logger
}

var standardLogger *StandardLogger = nil
var once sync.Once

const LOGDIR = "/var/log/session-counter"

var logLevel = logrus.InfoLevel

func (l *StandardLogger) SetLogLevel(level string) {
	// log.Println("setting log level to", level)
	switch strings.ToLower(level) {
	case "debug":
		logLevel = logrus.DebugLevel
	case "info":
		logLevel = logrus.InfoLevel
	case "warn":
		logLevel = logrus.WarnLevel
	default:
		logLevel = logrus.ErrorLevel
	}
	standardLogger.SetLevel(logLevel)
	// log.Println("log level is now", l.GetLogLevelName())
}

func (l *StandardLogger) GetLogLevelName() string {
	switch logLevel {
	case logrus.DebugLevel:
		return "DEBUG"
	case logrus.InfoLevel:
		return "INFO"
	case logrus.WarnLevel:
		return "WARN"
	case logrus.ErrorLevel:
		return "ERROR"
	case logrus.FatalLevel:
		return "FATAL"
	}
	return "UNKNOWN"
}

func NewLogger(cfg *config.Config) *StandardLogger {
	once.Do(func() {
		initLogger(cfg)
		if cfg != nil {
			standardLogger.SetLogLevel(cfg.GetLogLevel())
		} else {
			standardLogger.SetLogLevel("FATAL")
		}
	})

	if standardLogger != nil {
		//log.Println("returning non-nil sl")
		return standardLogger
	}
	if cfg == nil && standardLogger != nil {
		//log.Println("config nil, returning non-nil sl")
		return standardLogger
	} else {
		//log.Println("Falling back on UnsafeNewLogger...")
		return UnsafeNewLogger(cfg)
	}

}

// For unit testing only
func UnsafeNewLogger(cfg *config.Config) (sl *StandardLogger) {
	if standardLogger == nil {
		initLogger(cfg)
	}
	return standardLogger
}

// Convoluted for use within libraries...
func initLogger(cfg *config.Config) {
	var baseLogger = logrus.New()

	if baseLogger == nil {
		log.Println("baseLogger is nil")
	}
	// If we have a config file, grab the loggers defined there.
	// Otherwise, use stderr.
	loggers := []string{"local:stderr"}

	if cfg != nil {
		loggers = cfg.GetLoggers()
	}
	log.Println("loggers", loggers)

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
	log.Println("writers", writers)
	mw := io.MultiWriter(writers...)
	baseLogger.SetOutput(mw)
	baseLogger.SetReportCaller(true)
	standardLogger = &StandardLogger{baseLogger}
	standardLogger.Formatter = &logrus.JSONFormatter{}
}

const (
	FATAL = iota
	ERROR
	WARN
	INFO
	DEBUG
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

// func (l *StandardLogger) Base(e Event, file string, line int, args ...interface{}) {
// 	fields := logrus.Fields{
// 		"file": file,
// 		"line": line,
// 	}
// 	if l == nil {
// 		log.Println("logger `l` not initialized.")
// 		log.Println("NewLogger() does not work in stateless contexts.")
// 		log.Println("You may be looking for UnsafeNewLogger().")
// 		log.Println("Creating a new (default) logger for you...")
// 		initLogger(nil)
// 	}

// 	// l.WithFields(fields).Debugf(e.message, args...)
// 	log.Printf(e.message, args...)

// 	switch e.level {
// 	case DEBUG:
// 		if logLevel >= DEBUG {
// 			// log.Println("DEBUG", fmt.Sprintf(e.message, args...))
// 			l.WithFields(fields).Debugf(e.message, args...)
// 		}
// 	case INFO:
// 		if logLevel >= INFO {
// 			l.WithFields(fields).Infof(e.message, args...)
// 		}
// 	case WARN:
// 		if logLevel >= WARN {
// 			l.WithFields(fields).Warnf(e.message, args...)
// 		}
// 	case ERROR:
// 		// Always log ERROR level log events.
// 		l.WithFields(fields).Errorf(e.message, args...)
// 	case FATAL:
// 		// Always log FATAL level log events.
// 		l.WithFields(fields).Fatalf(e.message, args...)
// 		// We're leaving, on a jet plane...
// 		os.Exit(-1)
// 	}
// }

// func (l *StandardLogger) Debug(msg string, args ...interface{}) {
// 	_, file, line, _ := runtime.Caller(1)
// 	l.Base(debugMsg, filepath.Base(file), line, fmt.Sprintf(msg, args...))
// }

// // InvalidArg is a standard error message
// func (l *StandardLogger) Info(msg string, args ...interface{}) {
// 	_, file, line, _ := runtime.Caller(1)
// 	l.Base(infoMsg, filepath.Base(file), line, fmt.Sprintf(msg, args...))
// }

// func (l *StandardLogger) Warn(msg string, args ...interface{}) {
// 	_, file, line, _ := runtime.Caller(1)
// 	l.Base(warnMsg, filepath.Base(file), line, fmt.Sprintf(msg, args...))
// }

// func (l *StandardLogger) Error(msg string, args ...interface{}) {
// 	_, file, line, _ := runtime.Caller(1)
// 	l.Base(errorMsg, filepath.Base(file), line, fmt.Sprintf(msg, args...))
// }

// func (l *StandardLogger) Fatal(msg string, args ...interface{}) {
// 	_, file, line, _ := runtime.Caller(1)
// 	l.Base(fatalMsg, filepath.Base(file), line, fmt.Sprintf(msg, args...))
// }

// func (l *StandardLogger) ExeNotFound(path string) {
// 	_, file, line, _ := runtime.Caller(1)
// 	l.Base(notFound, filepath.Base(file), line, path)
// }

// func (l *StandardLogger) Length(arrname string, arr ...interface{}) {
// 	_, file, line, _ := runtime.Caller(1)
// 	l.Base(lengthMsg, filepath.Base(file), line, arrname, len(arr))
// }
