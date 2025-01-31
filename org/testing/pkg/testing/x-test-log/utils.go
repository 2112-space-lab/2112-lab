package xtestlog

import (
	"log/slog"
)

func LogAndPanicIfError(logger *slog.Logger, msg string, err error, args ...any) {
	if err == nil {
		return
	}
	logger.Error("invalid init sequence",
		slog.Any(AttrErrorKey, err),
		slog.Group("errorContext", args...),
	)
	panic(err)
}
