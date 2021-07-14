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
	"gsa.gov/18f/internal/config"
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
		return standardLogger
	}
	if cfg == nil && standardLogger != nil {
		return standardLogger
	} else {
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
			api := NewAPILogger(cfg)
			writers = append(writers, api)
		}
	}
	mw := io.MultiWriter(writers...)
	baseLogger.SetOutput(mw)
	baseLogger.SetReportCaller(true)
	standardLogger = &StandardLogger{baseLogger}
	standardLogger.Formatter = &logrus.JSONFormatter{}
}
