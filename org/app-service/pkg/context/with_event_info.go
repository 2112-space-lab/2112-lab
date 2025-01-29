package xcontext

import (
	"context"
	"time"

	log "github.com/Elbujito/2112/src/app-service/pkg/log"
)

const (
	// EventInfoKey context key for eventInfo
	EventInfoKey contextKey = "__eventInfo"

	eventTimeLogKey      = "__eventTime"
	eventUIDLogKey       = "__eventUID"
	eventTypeLogKey      = "__eventType"
	eventProcessorLogKey = "__eventHandler"
)

// EventInfo holds information about event to enrich context and logger
type EventInfo struct {
	EventUID    string
	EventType   string
	HandlerName string
	EventTime   time.Time
}

// AppendToLogFields adds info fields into input log fields
func (info EventInfo) AppendToLogFields(logFields log.Fields) log.Fields {
	logFields[eventTimeLogKey] = info.EventTime
	logFields[eventUIDLogKey] = info.EventUID
	logFields[eventTypeLogKey] = info.EventType
	if info.HandlerName != "" {
		logFields[eventProcessorLogKey] = info.HandlerName
	}
	return logFields
}

// WithEventInfo adds EventInfo to context and logger
func WithEventInfo(eventUID string, eventType string, processorName string, eventTime time.Time) ContextEnhancer {
	return func(parentCtx context.Context, logFields log.Fields) (context.Context, log.Fields) {
		info := EventInfo{
			EventUID:    eventUID,
			EventType:   eventType,
			HandlerName: processorName,
			EventTime:   eventTime,
		}
		ctx := context.WithValue(parentCtx, EventInfoKey, info)
		logFields = info.AppendToLogFields(logFields)
		return ctx, logFields
	}
}

// ReadEventInfo tries to read the EventInfo from the context and returns a boolean indicating whether the Event info was found
func ReadEventInfo(ctx context.Context) (event EventInfo, found bool) {
	return ReadValue(ctx, EventInfoKey, EventInfo{})
}
