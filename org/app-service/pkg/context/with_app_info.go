package xcontext

import (
	"context"

	log "github.com/Elbujito/2112/src/app-service/pkg/log"
)

const (
	// AppInfoKey context key
	AppInfoKey contextKey = "appInfo"
)

// AppInfo holds application general info. Used to enrich logs
type AppInfo struct {
	AppName      string
	InstanceName string
	AppVersion   string
}

// AppendToLogFields adds info fields into input log fields
func (info AppInfo) AppendToLogFields(logFields log.Fields) log.Fields {
	logFields["__appName"] = info.AppName
	if info.InstanceName != "" {
		logFields["__instanceName"] = info.InstanceName
	}
	logFields["__appVersion"] = info.AppVersion
	return logFields
}

// WithAppInfo adds AppInfo to context and logger
func WithAppInfo(appName string, instanceName string, appVersion string) ContextEnhancer {
	return func(parentCtx context.Context, logFields log.Fields) (context.Context, log.Fields) {
		info := AppInfo{
			AppName:      appName,
			InstanceName: instanceName,
			AppVersion:   appVersion,
		}
		ctx := context.WithValue(parentCtx, AppInfoKey, info)
		logFields = info.AppendToLogFields(logFields)
		return ctx, logFields
	}
}
