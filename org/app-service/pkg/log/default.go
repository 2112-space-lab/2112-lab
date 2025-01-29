package xlog

import (
	"context"
	"os"
)

var instance *Logger

func getInstance() *Logger {
	if instance == nil {
		instance, _ = NewLogger(os.Stdout, DebugLevel, LoggerTypes.Logrus())
	}
	return instance
}

// SetDefaultLogger set Global Default Logger
func SetDefaultLogger(logger *Logger) {
	instance = logger
}

// Info logs at LevelInfo.
func Info(msg string) {
	getInstance().Info(msg)
}

// Infof logs at LevelInfo.
func Infof(msg string, args ...any) {
	getInstance().Infof(msg, args...)
}

// InfoContext logs at LevelInfo with the given context.
func InfoContext(ctx context.Context, msg string) {
	getInstance().InfoContext(ctx, msg)
}

// Warn logs at LevelWarn.
func Warn(msg string) {
	getInstance().Warn(msg)
}

// Warnf logs at LevelWarn.
func Warnf(msg string, args ...any) {
	getInstance().Warnf(msg, args...)
}

// WarnContext logs at LevelWarn with the given context.
func WarnContext(ctx context.Context, msg string) {
	getInstance().WarnContext(ctx, msg)
}

// Error logs at LevelError.
func Error(msg string) {
	getInstance().Error(msg)
}

// Errorf logs at LevelError.
func Errorf(msg string, args ...any) {
	getInstance().Errorf(msg, args...)
}

// ErrorContext logs at LevelError with the given context.
func ErrorContext(ctx context.Context, msg string) {
	getInstance().ErrorContext(ctx, msg)
}

// Fatal logs the error and exit
func Fatal(msg string, err error) {
	getInstance().Fatal(msg, err)
}

// Fatalf logs the error and exit
func Fatalf(msg string, args ...interface{}) {
	getInstance().Fatalf(msg, args...)
}

// Debug logs at LevelDebug.
func Debug(msg string) {
	getInstance().Debug(msg)
}

// Debugf logs at LevelDebug.
func Debugf(msg string, args ...any) {
	getInstance().Debugf(msg, args...)
}

// DebugContext logs at LevelDebug with the given context.
func DebugContext(ctx context.Context, msg string) {
	getInstance().DebugContext(ctx, msg)
}

// Trace logs at LevelDebug.
func Trace(msg string) {
	getInstance().Trace(msg)
}

// Tracef logs at LevelDebug.
func Tracef(msg string, args ...any) {
	getInstance().Tracef(msg, args...)
}

// WithGroup returns a Logger that starts a group, if name is non-empty. The keys of all attributes added to the Logger will be qualified by the given name.
func WithGroup(name string) *Logger {
	return getInstance().WithGroup(name)
}

// WithField alias to With() function
func WithField(key string, value interface{}) *Logger {
	return getInstance().WithField(key, value)
}

// WithFields alias to With() function
func WithFields(fields map[string]interface{}) *Logger {
	return getInstance().WithFields(fields)
}

// WithError alias to With() function with key value "error"
func WithError(err error) *Logger {
	return getInstance().WithError(err)
}

// WithContext alias to With() function with key value "context"
func WithContext(ctx context.Context) *Logger {
	return getInstance().WithContext(ctx)
}
