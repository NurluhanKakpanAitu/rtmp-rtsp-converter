package logger

import (
    "os"
    
    "github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func Init(level, format string) {
    Log = logrus.New()
    Log.SetOutput(os.Stdout)
    
    // Set log level
    switch level {
    case "debug":
        Log.SetLevel(logrus.DebugLevel)
    case "info":
        Log.SetLevel(logrus.InfoLevel)
    case "warn":
        Log.SetLevel(logrus.WarnLevel)
    case "error":
        Log.SetLevel(logrus.ErrorLevel)
    default:
        Log.SetLevel(logrus.InfoLevel)
    }
    
    // Set log format
    if format == "json" {
        Log.SetFormatter(&logrus.JSONFormatter{})
    } else {
        Log.SetFormatter(&logrus.TextFormatter{
            FullTimestamp: true,
        })
    }
}

func GetLogger() *logrus.Logger {
    if Log == nil {
        Init("info", "text")
    }
    return Log
}