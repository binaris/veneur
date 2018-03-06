package veneur

import (
	"github.com/sirupsen/logrus"
)

// EntryLogger logs to some level on entry, possibly entry.Info.
type EntryLogger func(entry *logrus.Entry, args ...interface{})

// MakeLevelLogger returns a function that logs level on entry, possible entry.Info.
func MakeLevelLogger(level logrus.Level) EntryLogger {
	switch level {
	case logrus.PanicLevel:
		return func(entry *logrus.Entry, args ...interface{}) { entry.Panic(args) }
	case logrus.FatalLevel:
		return func(entry *logrus.Entry, args ...interface{}) { entry.Fatal(args) }
	case logrus.ErrorLevel:
		return func(entry *logrus.Entry, args ...interface{}) { entry.Error(args) }
	case logrus.WarnLevel:
		return func(entry *logrus.Entry, args ...interface{}) { entry.Warn(args) }
	case logrus.InfoLevel:
		return func(entry *logrus.Entry, args ...interface{}) { entry.Info(args) }
	case logrus.DebugLevel:
		return func(entry *logrus.Entry, args ...interface{}) { entry.Debug(args) }
	default:
		log.Errorf("Unexpected verbose log level %v (using info)", level)
		return func(entry *logrus.Entry, args ...interface{}) { entry.Info(args) }
	}
}

// GetVerboseLogLevel parses logLevelStr or Info if empty.
func GetVerboseLogLevel(logLevelStr string) logrus.Level {
	if logLevelStr == "" {
		return logrus.InfoLevel
	}
	// ParseLevel logs errors.
	level, _ := logrus.ParseLevel(logLevelStr)
	return level
}
