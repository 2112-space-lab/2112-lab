package xlog

import (
	"context"
	"fmt"
	"io"
)

// Fields alias type
type Fields = map[string]interface{}

// Entry alias type
type Entry = Logger

// Logger Wrapper
type Logger struct {
	log loggerHandler
}

// Level the log levels. Panic is not currently supported
type Level string

const (
	// InfoLevel info log level
	InfoLevel Level = "info"
	// WarnLevel warning log level
	WarnLevel Level = "warn"
	// DebugLevel debug log level
	DebugLevel Level = "debug"
	// TraceLevel trace log level
	TraceLevel Level = "trace"
	// ErrorLevel error log level
	ErrorLevel Level = "error"
)

type loggerHandler interface {
	Info(msg string)
	Infof(msg string, args ...any)
	InfoContext(ctx context.Context, msg string)
	Warn(msg string)
	Warnf(msg string, args ...any)
	WarnContext(ctx context.Context, msg string)
	Error(msg string)
	Errorf(msg string, args ...any)
	ErrorContext(ctx context.Context, msg string)
	Fatal(msg string, err error)
	Fatalf(msg string, args ...interface{})
	Debug(msg string)
	Debugf(msg string, args ...any)
	DebugContext(ctx context.Context, msg string)
	Trace(msg string)
	Tracef(msg string, args ...any)
	WithGroup(name string) loggerHandler
	WithField(key string, value interface{}) loggerHandler
	WithFields(fields map[string]interface{}) loggerHandler
	WithError(err error) loggerHandler
	WithContext(ctx context.Context) loggerHandler
}

// NewLogger custom Logger Constructor
func NewLogger(w io.Writer, level Level, loggerType LoggerType) (*Logger, error) {
	var handler loggerHandler
	var err error

	switch loggerType {
	case LoggerTypes.Logrus():
		handler = newLogrusLogger(w, level)
	case LoggerTypes.SLog():
		handler = newSlogLogger(w, level)
	default:
		err = fmt.Errorf("invalid logger type")
	}

	logger := &Logger{
		log: handler,
	}
	return logger, err
}

// Info logs at LevelInfo.
func (l *Logger) Info(msg string) {
	l.log.Info(msg)
}

// Infof logs at LevelInfo.
func (l *Logger) Infof(msg string, args ...any) {
	l.log.Infof(msg, args...)
}

// InfoContext logs at LevelInfo with the given context.
func (l *Logger) InfoContext(ctx context.Context, msg string) {
	l.log.InfoContext(ctx, msg)
}

// Warn logs at LevelWarn.
func (l *Logger) Warn(msg string) {
	l.log.Warn(msg)
}

// Warnf logs at LevelWarn.
func (l *Logger) Warnf(msg string, args ...any) {
	l.log.Warnf(msg, args...)
}

// WarnContext logs at LevelWarn with the given context.
func (l *Logger) WarnContext(ctx context.Context, msg string) {
	l.log.WarnContext(ctx, msg)
}

// Error logs at LevelError.
func (l *Logger) Error(msg string) {
	l.log.Error(msg)
}

// Errorf logs at LevelError.
func (l *Logger) Errorf(msg string, args ...any) {
	l.log.Errorf(msg, args...)
}

// ErrorContext logs at LevelError with the given context.
func (l *Logger) ErrorContext(ctx context.Context, msg string) {
	l.log.ErrorContext(ctx, msg)
}

// Fatal logs the error and exit
func (l *Logger) Fatal(msg string, err error) {
	l.log.Fatal(msg, err)
}

// Fatalf logs the error and exit
func (l *Logger) Fatalf(msg string, args ...interface{}) {
	l.log.Fatalf(msg, args...)
}

// Debug logs at LevelDebug.
func (l *Logger) Debug(msg string) {
	l.log.Debug(msg)
}

// Debugf logs at LevelDebug.
func (l *Logger) Debugf(msg string, args ...any) {
	l.log.Debugf(msg, args...)
}

// DebugContext logs at LevelDebug with the given context.
func (l *Logger) DebugContext(ctx context.Context, msg string) {
	l.log.DebugContext(ctx, msg)
}

// Trace logs at LevelDebug.
func (l *Logger) Trace(msg string) {
	l.log.Trace(msg)
}

// Tracef logs at LevelDebug.
func (l *Logger) Tracef(msg string, args ...any) {
	l.log.Tracef(msg, args...)
}

// WithGroup returns a Logger that starts a group, if name is non-empty. The keys of all attributes added to the Logger will be qualified by the given name.
func (l *Logger) WithGroup(name string) *Logger {
	return &Logger{
		log: l.log.WithGroup(name),
	}
}

// WithField alias to With() func (l *Logger)tion
func (l *Logger) WithField(key string, value interface{}) *Logger {
	return &Logger{
		log: l.log.WithField(key, value),
	}
}

// WithFields alias to With() func (l *Logger)tion
func (l *Logger) WithFields(fields Fields) *Logger {
	return &Logger{
		log: l.log.WithFields(fields),
	}
}

// WithError alias to With() func (l *Logger)tion with key value "error"
func (l *Logger) WithError(err error) *Logger {
	return &Logger{
		log: l.log.WithError(err),
	}
}

// WithContext alias to With() func (l *Logger)tion with key value "context"
func (l *Logger) WithContext(ctx context.Context) *Logger {
	return &Logger{
		log: l.log.WithContext(ctx),
	}
}
