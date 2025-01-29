package xlog

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"runtime"
	"sync"
)

const (
	// BufferSize to allocate a buffer to be used to dump strackTrace in case of Fatal error
	BufferSize = 1 << 16 // 65,536
)

var _ loggerHandler = &slogLogger{}

type slogLogger struct {
	mu   sync.Mutex
	log  *slog.Logger
	keys map[string]interface{}
}

// newSlogLogger slogger Constructor
func newSlogLogger(writer io.Writer, level Level) *slogLogger {
	return &slogLogger{
		keys: make(map[string]interface{}),
		log: slog.New(slog.NewJSONHandler(
			writer,
			&slog.HandlerOptions{
				Level: getSLogLevel(level),
			},
		)),
	}
}

func getSLogLevel(level Level) slog.Level {
	switch level {
	case InfoLevel:
		return slog.LevelInfo
	case WarnLevel:
		return slog.LevelWarn
	case DebugLevel:
		return slog.LevelDebug
	case TraceLevel:
		return slog.LevelDebug - 10
	case ErrorLevel:
		return slog.LevelError
	}
	//default is INFO Level
	return slog.LevelInfo
}

func (l *slogLogger) returnAllValues() []any {
	l.mu.Lock()
	defer l.mu.Unlock()

	kV := make([]any, 0, len(l.keys))
	for k, v := range l.keys {
		kV = append(kV, slog.Any(k, v))
	}
	return kV
}

func (l *slogLogger) appendToKeys(values Fields) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.keys == nil {
		l.keys = make(map[string]interface{})
	}
	for k, v := range values {
		l.keys[k] = v
	}
}

func (l *slogLogger) cloneKeys() map[string]interface{} {
	l.mu.Lock()
	defer l.mu.Unlock()

	c := make(map[string]interface{})
	for k, v := range l.keys {
		c[k] = v
	}
	return c
}

// Info logs at LevelInfo.
func (l *slogLogger) Info(msg string) {
	l.log.Info(msg, l.returnAllValues()...)
}

// Infof logs at LevelInfo.
func (l *slogLogger) Infof(msg string, args ...any) {
	l.Info(fmt.Sprintf(msg, args...))
}

// InfoContext logs at LevelInfo with the given context.
func (l *slogLogger) InfoContext(ctx context.Context, msg string) {
	l.log.InfoContext(ctx, msg, l.returnAllValues()...)
}

// Warn logs at LevelWarn.
func (l *slogLogger) Warn(msg string) {
	l.log.Warn(msg, l.returnAllValues()...)
}

// Warnf logs at LevelWarn.
func (l *slogLogger) Warnf(msg string, args ...any) {
	l.Warn(fmt.Sprintf(msg, args...))
}

// WarnContext logs at LevelWarn with the given context.
func (l *slogLogger) WarnContext(ctx context.Context, msg string) {
	l.log.WarnContext(ctx, msg, l.returnAllValues()...)
}

// Error logs at LevelError.
func (l *slogLogger) Error(msg string) {
	l.log.Error(msg, l.returnAllValues()...)
}

// Errorf logs at LevelError.
func (l *slogLogger) Errorf(msg string, args ...any) {
	l.Error(fmt.Sprintf(msg, args...))
}

// ErrorContext logs at LevelError with the given context.
func (l *slogLogger) ErrorContext(ctx context.Context, msg string) {
	l.log.ErrorContext(ctx, msg, l.returnAllValues()...)
}

// Fatal logs the error and exit
func (l *slogLogger) Fatal(msg string, err error) {
	buf := make([]byte, BufferSize)
	runtime.Stack(buf, true)

	l.log.Error(fmt.Sprintf(msg, err), "stackTrace", buf)
	panic(err)
}

// Fatalf logs the error and exit
func (l *slogLogger) Fatalf(msg string, args ...interface{}) {
	buf := make([]byte, BufferSize)
	runtime.Stack(buf, true)

	l.log.Error(fmt.Sprintf(msg, args...), "stackTrace", buf)
	panic(fmt.Sprintf(msg, args...))
}

// Debug logs at LevelDebug.
func (l *slogLogger) Debug(msg string) {
	l.log.Debug(msg, l.returnAllValues()...)
}

// Debugf logs at LevelDebug.
func (l *slogLogger) Debugf(msg string, args ...any) {
	l.Debug(fmt.Sprintf(msg, args...))
}

// DebugContext logs at LevelDebug with the given context.
func (l *slogLogger) DebugContext(ctx context.Context, msg string) {
	l.log.DebugContext(ctx, msg, l.returnAllValues()...)
}

// Trace logs at LevelDebug.
func (l *slogLogger) Trace(msg string) {
	l.log.Log(context.Background(), getSLogLevel(TraceLevel), msg, l.returnAllValues()...)
}

// Tracef logs at LevelDebug.
func (l *slogLogger) Tracef(msg string, args ...any) {
	l.Trace(fmt.Sprintf(msg, args...))
}

// With returns a Logger that includes the given attributes in each output operation. Arguments are converted to attributes as if by [Logger.Log].
func (l *slogLogger) with(values Fields) loggerHandler {
	instance := &slogLogger{
		log:  l.log.With(), // clones the internal config
		keys: l.cloneKeys(),
	}

	instance.appendToKeys(values)
	return instance
}

// WithGroup returns a Logger that starts a group, if name is non-empty. The keys of all attributes added to the Logger will be qualified by the given name.
func (l *slogLogger) WithGroup(name string) loggerHandler {
	return &slogLogger{
		log:  l.log.WithGroup(name),
		keys: l.cloneKeys(),
	}
}

// WithField alias to With() function
func (l *slogLogger) WithField(key string, value interface{}) loggerHandler {
	return l.with(Fields{key: value})
}

// WithFields alias to With() function
func (l *slogLogger) WithFields(fields Fields) loggerHandler {
	return l.with(fields)
}

// WithError alias to With() function with key value "error"
func (l *slogLogger) WithError(err error) loggerHandler {
	return l.with(Fields{"error": err})
}

// WithContext alias to With() function with key value "context"
func (l *slogLogger) WithContext(ctx context.Context) loggerHandler {
	return l.with(Fields{"context": ctx})
}
