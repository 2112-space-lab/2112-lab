package xlog

import (
	"context"
	"fmt"
	"io"
	"github.com/sirupsen/logrus"
)

var _ loggerHandler = &logrusLogger{}

type logrusLogger struct {
	log *logrus.Entry
}

// newSlogLogger slogger Constructor
func newLogrusLogger(w io.Writer, level Level) *logrusLogger {
	logger := logrus.New()
	logger.SetOutput(w)
	logger.SetLevel(getLogrusLevel(level))
	return &logrusLogger{
		log: logrus.NewEntry(logger),
	}
}

// GetLogrusEntry returns logrus entry if loger type is logrus
// else returns new entry that prints to console with an error
func GetLogrusEntry(h *Logger) (*logrus.Entry, error) {
	logger, ok := h.log.(*logrusLogger)

	if !ok {
		return logrus.NewEntry(logrus.New()), fmt.Errorf("failed to cast logger to logrus")
	}

	return logger.log, nil
}

// SetLevel sets the default log level. This can be overridden in each log Instance
func (l *logrusLogger) SetLevel(level Level) {
	l.log.Logger.SetLevel(getLogrusLevel(level))
}

func getLogrusLevel(level Level) logrus.Level {
	switch level {
	case InfoLevel:
		return logrus.InfoLevel
	case WarnLevel:
		return logrus.WarnLevel
	case DebugLevel:
		return logrus.DebugLevel
	case TraceLevel:
		return logrus.TraceLevel
	case ErrorLevel:
		return logrus.ErrorLevel
	}
	//default is infoLevel
	return logrus.InfoLevel
}

// Info logs at LevelInfo.
func (l *logrusLogger) Info(msg string) {
	l.log.Info(msg)
}

// Infof logs at LevelInfo.
func (l *logrusLogger) Infof(msg string, args ...any) {
	l.log.Infof(msg, args...)
}

// InfoContext logs at LevelInfo with the given context.
func (l *logrusLogger) InfoContext(ctx context.Context, msg string) {
	l.log.Info(ctx, msg)
}

// Warn logs at LevelWarn.
func (l *logrusLogger) Warn(msg string) {
	l.log.Warn(msg)
}

// Warnf logs at LevelWarn.
func (l *logrusLogger) Warnf(msg string, args ...any) {
	l.log.Warnf(msg, args...)
}

// WarnContext logs at LevelWarn with the given context.
func (l *logrusLogger) WarnContext(ctx context.Context, msg string) {
	l.log.Warn(ctx, msg)
}

// Error logs at LevelError.
func (l *logrusLogger) Error(msg string) {
	l.log.Error(msg)
}

// Errorf logs at LevelError.
func (l *logrusLogger) Errorf(msg string, args ...any) {
	l.log.Errorf(msg, args...)
}

// ErrorContext logs at LevelError with the given context.
func (l *logrusLogger) ErrorContext(ctx context.Context, msg string) {
	l.log.Error(ctx, msg)
}

// Fatal logs the error and exit
func (l *logrusLogger) Fatal(msg string, err error) {
	l.log.Fatal(msg, err)
}

// Fatalf logs the error and exit
func (l *logrusLogger) Fatalf(msg string, args ...interface{}) {
	l.log.Fatalf(msg, args...)
}

// Debug logs at LevelDebug.
func (l *logrusLogger) Debug(msg string) {
	l.log.Debug(msg)
}

// Debugf logs at LevelDebug.
func (l *logrusLogger) Debugf(msg string, args ...any) {
	l.log.Debugf(msg, args...)
}

// DebugContext logs at LevelDebug with the given context.
func (l *logrusLogger) DebugContext(ctx context.Context, msg string) {
	l.log.Debug(ctx, msg)
}

// Trace logs at LevelDebug.
func (l *logrusLogger) Trace(msg string) {
	l.log.Trace(msg)
}

// Tracef logs at LevelDebug.
func (l *logrusLogger) Tracef(msg string, args ...any) {
	l.log.Tracef(msg, args...)
}

// With alias to WithField with key "fields".
func (l *logrusLogger) With(args ...any) loggerHandler {
	return l.WithField("fields", args)
}

// WithGroup alias to WithField key "group"
func (l *logrusLogger) WithGroup(name string) loggerHandler {
	return l.WithField("group", name)
}

// WithField Add a single field to the Entry.
func (l *logrusLogger) WithField(key string, value interface{}) loggerHandler {
	return &logrusLogger{
		log: l.log.WithField(key, value),
	}
}

// WithFields Add a map of fields to the Entry.
func (l *logrusLogger) WithFields(fields Fields) loggerHandler {
	return &logrusLogger{
		log: l.log.WithFields(fields),
	}
}

// WithError Add an error as single field (using the key defined in ErrorKey) to the Entry.
func (l *logrusLogger) WithError(err error) loggerHandler {
	return &logrusLogger{
		log: l.log.WithError(err),
	}
}

// WithContext Add a context to the Entry.
func (l *logrusLogger) WithContext(ctx context.Context) loggerHandler {
	return &logrusLogger{
		log: l.log.WithContext(ctx),
	}
}
