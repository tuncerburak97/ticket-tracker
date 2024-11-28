package logger

import (
	"github.com/sirupsen/logrus"
	"os"
	_ "os"
)

var Logger *logrus.Logger

func init() {
	Logger = logrus.New()
	Logger.SetOutput(os.Stdout)

	level, err := logrus.ParseLevel("info")
	if err != nil {
		level = logrus.InfoLevel
	}
	Logger.SetLevel(level)

	Logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
}

func ConfigureLogger(logFile string, logLevel string, jsonFormat bool) {
	if logFile != "" {
		file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			Logger.SetOutput(file)
		} else {
			Logger.Warn("Failed to log to file, using default stderr")
		}
	}

	if logLevel != "" {
		level, err := logrus.ParseLevel(logLevel)
		if err == nil {
			Logger.SetLevel(level)
		}
	}

	if jsonFormat {
		Logger.SetFormatter(&logrus.JSONFormatter{})
	}
}
