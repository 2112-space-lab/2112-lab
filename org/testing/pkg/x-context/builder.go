package xcontext

import (
	"context"

	log "github.com/org/2112-space-lab/org/testing/pkg/x-log"
)

// TraceUID type alias
type TraceUID string

// Caller type alias
type Caller string

// CallStack type alias
type CallStack string

// contextKey type alias
type contextKey string

// ContextEnhancer signature to be implemented by functions to enrich context and log fields
type ContextEnhancer func(parentCtx context.Context, logFields log.Fields) (context.Context, log.Fields)

// ContextBuilder takes in parent context and logger and enrich them with provided enhancers
func ContextBuilder(parentCtx context.Context, caller Caller, parentLogger *log.Entry, ctxOptions ...ContextEnhancer) (context.Context, *log.Entry) {
	logFields := log.Fields{}
	newCtx, logFields := withCallInfo(caller)(parentCtx, logFields)
	for _, enhance := range ctxOptions {
		newCtx, logFields = enhance(newCtx, logFields)
	}

	//! we re-inject from context info values to ensure all logs have it
	appInfo, found := ReadValue(newCtx, AppInfoKey, AppInfo{})
	if found {
		logFields = appInfo.AppendToLogFields(logFields)
	}

	eventInfo, found := ReadValue(newCtx, EventInfoKey, EventInfo{})
	if found {
		logFields = eventInfo.AppendToLogFields(logFields)
	}

	trackInfo, found := ReadValue(newCtx, TrackInfoKey, TrackInfo{})
	if found {
		logFields = trackInfo.AppendToLogFields(logFields)
	}

	traceUID, found := ReadValue(parentCtx, TraceUIDKey, "")
	if found && traceUID != "" {
		logFields[traceUIDLogKey] = traceUID
	}

	userInfo, found := ReadValue(parentCtx, userInfoKey, UserInfo{})
	if found {
		logFields = userInfo.AppendToLogFields(logFields)
	}

	var logger *log.Entry
	if parentLogger == nil {
		logger = log.WithContext(newCtx).WithFields(logFields)
	} else {
		logger = parentLogger.WithContext(newCtx).WithFields(logFields)
	}
	return newCtx, logger
}

// ReadValue safely extract value from context given a specific key
func ReadValue[TValue any](ctx context.Context, key contextKey, defaultValue TValue) (value TValue, found bool) {
	val := ctx.Value(key)
	if val != nil {
		typedVal, ok := val.(TValue)
		if ok {
			return typedVal, true
		}
	}
	return defaultValue, false
}

// PlanID, GmsTrackID, TrackInstanceUID, SatelliteID, DeviceUID or DeviceFullTag
// TraceUID
// Verb, URL/path ?, origin (ip), source, user agent
// Protocol (rest or grpc)
// }

// withGrpcInfo(ctx, funcName, map[], source, headers[])
// withRestInfo
// withDeviceInfo?
