package xcontext

import (
	"context"
	"strings"

	log "github.com/org/2112-space-lab/org/testing/pkg/x-log"
)

const (
	// CallInfoKey context key for callInfo
	CallInfoKey contextKey = "callInfo"

	callerLogKey    string = "__caller"
	callStackLogKey string = "__callStack"
)

// CallInfo holds function call tracing info for context
type CallInfo struct {
	Caller    Caller
	CallStack CallStack
}

// AppendToLogFields adds info fields into input log fields
func (info CallInfo) AppendToLogFields(logFields log.Fields) log.Fields {
	logFields[callerLogKey] = info.Caller
	logFields[callStackLogKey] = info.CallStack
	return logFields
}

func withCallInfo(caller Caller) ContextEnhancer {
	return func(parentCtx context.Context, logFields log.Fields) (context.Context, log.Fields) {
		callInfo, found := ReadValue(parentCtx, CallInfoKey, CallInfo{
			Caller:    caller,
			CallStack: ">" + CallStack(caller),
		})
		if found {
			callInfo.CallStack = computeCallStack(caller, callInfo.CallStack)
			callInfo.Caller = caller
		}
		newCtx := context.WithValue(parentCtx, CallInfoKey, callInfo)
		logFields = callInfo.AppendToLogFields(logFields)
		return newCtx, logFields
	}
}

func computeCallStack(caller Caller, currentStack CallStack) CallStack {
	if strings.HasSuffix(string(currentStack), string(caller)) {
		return currentStack
	}
	return currentStack + ">" + CallStack(caller)
}
