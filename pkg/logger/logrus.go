package logger

import (
	"github.com/sirupsen/logrus"
)

type logrusLogger struct {
	logger *logrus.Entry
}

var _ Logger = (*logrusLogger)(nil)

func NewLogrusLogger(level string) *logrusLogger {
	var logrusLevel logrus.Level

	switch level {
	case "error":
		logrusLevel = logrus.ErrorLevel
	case "warn":
		logrusLevel = logrus.WarnLevel
	case "info":
		logrusLevel = logrus.InfoLevel
	case "debug":
		logrusLevel = logrus.DebugLevel
	default:
		logrusLevel = logrus.InfoLevel
	}

	logger := logrus.New()
	logger.SetLevel(logrusLevel)

	return &logrusLogger{logger: logrus.NewEntry(logger)}
}

func (l *logrusLogger) Named(name string) Logger {
	return &logrusLogger{
		logger: l.logger.WithField("name", name),
	}
}

func (l *logrusLogger) Debug(message string, args M) {
	l.logger.WithFields(logrus.Fields(args)).Debug(message)
}

func (l *logrusLogger) Info(message string, args M) {
	l.logger.WithFields(logrus.Fields(args)).Info(message)
}

func (l *logrusLogger) Warn(message string, args M) {
	l.logger.WithFields(logrus.Fields(args)).Warn(message)
}

func (l *logrusLogger) Error(message string, args M) {
	l.logger.WithFields(logrus.Fields(args)).Error(message)
}

func (l *logrusLogger) Fatal(message string, args M) {
	l.logger.WithFields(logrus.Fields(args)).Fatal(message)
}
