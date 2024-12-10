package logger

import (
	"github.com/sirupsen/logrus"
	"os"
	"sync"
)

var (
	Logger *logrus.Logger
	once   sync.Once
)

// ReportCallerHook controls ReportCaller behavior based on log level
type ReportCallerHook struct{}

// Levels specifies the log levels for the hook
func (hook *ReportCallerHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// Fire toggles ReportCaller for error level logs
func (hook *ReportCallerHook) Fire(entry *logrus.Entry) error {
	entry.Logger.SetReportCaller(entry.Level == logrus.ErrorLevel)
	return nil
}

// GetLogger initializes or returns the singleton Logger instance
func GetLogger() *logrus.Logger {
	once.Do(func() {
		initLogger()
	})
	return Logger
}

// initLogger initializes the Logger with default settings
func initLogger() {
	Logger = logrus.New()
	Logger.SetOutput(os.Stdout)

	// Default log level
	Logger.SetLevel(logrus.InfoLevel)

	// Default formatter
	Logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// Add ReportCallerHook
	Logger.AddHook(&ReportCallerHook{})
}

// ConfigureLogger updates the Logger configuration dynamically
func ConfigureLogger(logFile string, logLevel string, jsonFormat bool) error {
	if logFile != "" {
		file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return err
		}
		Logger.SetOutput(file)
	}

	if logLevel != "" {
		level, err := logrus.ParseLevel(logLevel)
		if err != nil {
			return err
		}
		Logger.SetLevel(level)
	}

	if jsonFormat {
		Logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		})
	} else {
		Logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
		})
	}

	return nil
}
