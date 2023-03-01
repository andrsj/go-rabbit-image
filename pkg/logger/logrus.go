package logger

import (
	"github.com/sirupsen/logrus"
)

type logrusLogger struct {
	logger *logrus.Entry
}

var _ Logger = (*logrusLogger)(nil)

func NewLogrusLogger(level string) *logrusLogger {
	var l logrus.Level

	switch level {
	case "error":
		l = logrus.ErrorLevel
	case "warn":
		l = logrus.WarnLevel
	case "info":
		l = logrus.InfoLevel
	case "debug":
		l = logrus.DebugLevel
	default:
		l = logrus.InfoLevel
	}

	logger := logrus.New()
	logger.SetLevel(l)
	logger.SetFormatter(&logrus.JSONFormatter{})
	return &logrusLogger{logger: logrus.NewEntry(logger)}
}

func (l *logrusLogger) Named(name string) Logger {
	return &logrusLogger{
		logger: l.logger.WithField("name", name),
	}
}

func (l *logrusLogger) Debug(message string, args ...interface{}) {
	l.logger.WithFields(logrus.Fields{"args": args}).Debug(message)
}

func (l *logrusLogger) Info(message string, args ...interface{}) {
	l.logger.WithFields(logrus.Fields{"args": args}).Info(message)
}

func (l *logrusLogger) Warn(message string, args ...interface{}) {
	l.logger.WithFields(logrus.Fields{"args": args}).Warn(message)
}

func (l *logrusLogger) Error(message string, args ...interface{}) {
	l.logger.WithFields(logrus.Fields{"args": args}).Error(message)
}

func (l *logrusLogger) Fatal(message string, args ...interface{}) {
	l.logger.WithFields(logrus.Fields{"args": args}).Fatal(message)
}
