package xcontext

import (
	"context"

	log "github.com/Elbujito/2112/src/app-service/pkg/log"
	"github.com/google/uuid"
)

const (
	// TraceUIDKey context key for traceUID
	TraceUIDKey contextKey = "traceUID"

	traceUIDLogKey string = "__traceUID"
)

// WithTrace adds traceUID to context and logger
func WithTrace(traceUID TraceUID, forceOverride bool) ContextEnhancer {
	return func(parentCtx context.Context, logFields log.Fields) (context.Context, log.Fields) {
		if traceUID == "" {
			newUUID, errUUID := uuid.NewV7()
			if errUUID != nil {
				panic(errUUID)
			}
			traceUID = TraceUID(newUUID.String())
		}
		traceUID, found := ReadValue(parentCtx, TraceUIDKey, traceUID)
		if found && !forceOverride {
			return parentCtx, logFields
		}
		newCtx := context.WithValue(parentCtx, TraceUIDKey, traceUID)
		logFields[traceUIDLogKey] = traceUID
		return newCtx, logFields
	}
}
