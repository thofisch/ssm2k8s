package logging

import (
	"github.com/sirupsen/logrus"
	"os"
)

type Logger interface {
	Printf(format string, args ...interface{})
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
}

type logrusLogger struct {
	logger *logrus.Logger
}

func NewLogger() Logger {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	logger.SetOutput(os.Stdout)
	//logger.SetFormatter(&logrus.JSONFormatter{})

	return &logrusLogger{
		logger: logger,
	}
}

//func NewContextLogger(c logrus.Fields) func(f logrus.Fields) *logrus.Entry {
//	return func(f logrus.Fields) *logrus.Entry {
//		for k, v := range c {
//			f[k] = v
//		}
//		return logrus.WithFields(f)
//	}
//}

func (l *logrusLogger) Printf(format string, args ...interface{}) {
	l.logger.Printf(format, args...)
}

func (l *logrusLogger) Debug(args ...interface{}) {
	l.logger.Debugln(args...)
}

func (l *logrusLogger) Debugf(format string, args ...interface{}) {
	l.logger.Debugf(format, args...)
}

func (l *logrusLogger) Info(args ...interface{}) {
	l.logger.Infoln(args...)
}

func (l *logrusLogger) Infof(format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}

func (l *logrusLogger) Warn(args ...interface{}) {
	l.logger.Warnln(args...)
}

func (l *logrusLogger) Warnf(format string, args ...interface{}) {
	l.logger.Warnf(format, args...)
}

func (l *logrusLogger) Error(args ...interface{}) {
	l.logger.Errorln(args...)
}

func (l *logrusLogger) Errorf(format string, args ...interface{}) {
	l.logger.Errorf(format, args...)
}

type nullLogger struct {
}

func NewNullLogger() Logger {
	return &nullLogger{}
}

func (nullLogger) Printf(format string, args ...interface{}) {
}

func (nullLogger) Debug(args ...interface{}) {
}

func (nullLogger) Debugf(format string, args ...interface{}) {
}

func (nullLogger) Info(args ...interface{}) {
}

func (nullLogger) Infof(format string, args ...interface{}) {
}

func (nullLogger) Warn(args ...interface{}) {
}

func (nullLogger) Warnf(format string, args ...interface{}) {
}

func (nullLogger) Error(args ...interface{}) {
}

func (nullLogger) Errorf(format string, args ...interface{}) {
}
