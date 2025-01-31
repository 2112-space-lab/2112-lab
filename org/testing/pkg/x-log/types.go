package xlog

// LoggerType enum value for supportedLoggerTypes
type loggerTypeStr string

// LoggerType is safe wrapper around valid enum values
type LoggerType struct{ source loggerTypeStr }

type loggerTypes struct{}

// LoggerTypes is a reference object holding helpers for allowed values
var LoggerTypes = loggerTypes{}

// Unknown is helper for default fallback value
func (dm loggerTypes) Unknown() LoggerType { return LoggerType{""} }

// LoggerType helper for value
func (dm loggerTypes) SLog() LoggerType { return LoggerType{"slog"} }

// LoggerType helper for value
func (dm loggerTypes) Logrus() LoggerType { return LoggerType{"logrus"} }