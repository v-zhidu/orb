package logging

//This is a package that make logger easier to use in your application
//based on logrus repository. Add context logging feature.

//TODO - Set logger formatter
import (
	"context"

	"github.com/sirupsen/logrus"
)

type loggerKeyType int

const loggerKey loggerKeyType = iota

var logEntry *logrus.Entry

var logger *logrus.Logger

//Fields return a map that contains fields to logging
type Fields map[string]interface{}

func init() {
	logger = logrus.New()
	logEntry = logrus.NewEntry(logger)
}

//SetLevel set logger output level
func SetLevel(strlevel string) {
	level, err := logrus.ParseLevel(strlevel)
	if err != nil {
		panic(err)
	}

	logger.SetLevel(level)
}

//Debugln returns message logging by debug level
func Debugln(message string) {
	logEntry.Debug(message)
}

//Infoln returns message logging by info level
func Infoln(message string) {
	logEntry.Info(message)
}

//Warnln returns message logging by warn level
func Warnln(message string) {
	logEntry.Warn(message)
}

//Errorln returns message logging by error level
func Errorln(message string) {
	logEntry.Error(message)
}

//WithError returns message logging by error level
func WithError(message string, err error) {
	logEntry.WithError(err).Error(message)
}

//Fatalln returns message logging by fatal level and exit the application
func Fatalln(message string) {
	logEntry.Fatal(message)
}

//Debug returns message logging by debug level
func Debug(message string, fields Fields) {
	logEntry.WithFields(logrus.Fields(fields)).Debug(message)
}

//Info returns message logging by info level
func Info(message string, fields Fields) {
	logEntry.WithFields(logrus.Fields(fields)).Info(message)
}

//Warn returns message logging by warn level
func Warn(message string, fields Fields) {
	logEntry.WithFields(logrus.Fields(fields)).Warn(message)
}

//Error returns message logging by error level
func Error(message string, fields Fields, err error) {
	logEntry.WithFields(logrus.Fields(fields)).WithError(err).Error(message)
}

//Fatal returns message logging by fatal level and exit the application
func Fatal(message string, fields Fields) {
	logEntry.WithFields(logrus.Fields(fields)).Fatal(message)
}

//NewContext returns a context has a logrus entry with default fields added
func NewContext(ctx context.Context, fields Fields) context.Context {
	return context.WithValue(ctx, loggerKey, WithContext(ctx).WithFields(logrus.Fields(fields)))
}

//WithContext returns a logrus entry with fields in a same context
func WithContext(ctx context.Context) *logrus.Entry {
	if ctx == nil {
		return logEntry
	}

	if ctxLogEntry, ok := ctx.Value(loggerKey).(*logrus.Entry); ok {
		return ctxLogEntry
	}

	return logEntry
}
