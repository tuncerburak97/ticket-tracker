package logger

import (
	"github.com/sirupsen/logrus"
	"os"
)

var Logger *logrus.Logger

func GetLogger() *logrus.Logger {
	return Logger
}

func init() {
	Logger = logrus.New()
	Logger.SetOutput(os.Stdout)

	level, err := logrus.ParseLevel("info")
	if err != nil {
		level = logrus.InfoLevel
	}
	Logger.SetLevel(level)

	Logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05", // Daha okunabilir format
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
		Logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05", // JSON formatında da uygulanıyor
		})
	} else {
		Logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05", // Daha okunabilir format
		})
	}
}
