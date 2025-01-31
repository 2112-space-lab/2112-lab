package xtestlog

import (
	"errors"
	"log/slog"
)

const (
	AttrErrorKey string = "error"
)

var (
	ErrRunLoggerAlreadyInit = errors.New("attempting to re-init runLogger")
)

const LevelTrace slog.Level = slog.LevelDebug - 10
