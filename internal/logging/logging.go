package logging

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
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
	level := logrus.InfoLevel
	if _, ok := os.LookupEnv("DEBUG"); ok {
		level = logrus.DebugLevel
	}
	logger.SetLevel(level)
	logger.SetOutput(os.Stdout)
	//logger.SetFormatter(&logrus.JSONFormatter{})

	return &logrusLogger{
		logger: logger,
	}
}

func NewConsoleLogger() Logger {
	logger := logrus.New()
	level := logrus.InfoLevel
	timestamp := false
	if _, ok := os.LookupEnv("DEBUG"); ok {
		level = logrus.DebugLevel
		timestamp = true
	}
	logger.SetLevel(level)
	logger.SetOutput(os.Stdout)
	logger.SetFormatter(&consoleFormatter{timestamp})

	return &logrusLogger{
		logger: logger,
	}
}

type consoleFormatter struct {
	timestamp bool
}

var colorMap = map[logrus.Level]string{
	logrus.TraceLevel: "90",
	logrus.DebugLevel: "90",
	logrus.InfoLevel:  "39",
	logrus.WarnLevel:  "33",
	logrus.ErrorLevel: "31",
	logrus.FatalLevel: "31",
	logrus.PanicLevel: "31",
}

func (f *consoleFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	sb := &strings.Builder{}
	sb.WriteString(fmt.Sprintf("\033[%sm", colorMap[entry.Level]))

	if f.timestamp {
		sb.WriteString(fmt.Sprintf("[%s] ", entry.Time.Format("2006-01-02 03:04:05")))
	}
	sb.WriteString(entry.Message)
	sb.WriteString("\033[0m\n")

	return []byte(sb.String()), nil
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
