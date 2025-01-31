package xcontext

import (
	"context"

	log "github.com/org/2112-space-lab/org/testing/pkg/x-log"
)

const (
	// TrackInfoKey context key for trackInfo
	TrackInfoKey contextKey = "TrackInfo"

	planIDLogKey           string = "__planID"
	trackSourceLogKey      string = "__trackSource"
	gmsTrackIDLogKey       string = "__gmsTrackID"
	trackInstanceUIDLogKey string = "__trackInstanceUID"
	satelliteIDLogKey      string = "__satelliteID"
)

// TrackInfo holds information about track to enrich context and logger
type TrackInfo struct {
	PlanID           string
	TrackSource      string
	GmsTrackID       string
	TrackInstanceUID string
	SatelliteID      string
	//TODO PlanID, source, GmsTrackID, TrackInstanceUID, SatelliteID, DeviceUID or DeviceFullTag
}

// NewTrackInfo standard constructor for TrackInfo
func NewTrackInfo(planID string, trackSource string, gmsTrackID string, trackInstanceUID string, satelliteID string) TrackInfo {
	return TrackInfo{
		PlanID:           planID,
		TrackSource:      trackSource,
		GmsTrackID:       gmsTrackID,
		TrackInstanceUID: trackInstanceUID,
		SatelliteID:      satelliteID,
	}
}

// AppendToLogFields adds info fields into input log fields
func (info TrackInfo) AppendToLogFields(logFields log.Fields) log.Fields {
	if info.PlanID != "" {
		logFields[planIDLogKey] = info.PlanID
	}
	if info.TrackSource != "" {
		logFields[trackSourceLogKey] = info.TrackSource
	}
	if info.GmsTrackID != "" {
		logFields[gmsTrackIDLogKey] = info.GmsTrackID
	}
	if info.TrackInstanceUID != "" {
		logFields[trackInstanceUIDLogKey] = info.TrackInstanceUID
	}
	if info.SatelliteID != "" {
		logFields[satelliteIDLogKey] = info.SatelliteID
	}
	return logFields
}

// WithTrackInfo adds TrackInfo to context and logger
func WithTrackInfo(ti *TrackInfo) ContextEnhancer {
	return func(parentCtx context.Context, logFields log.Fields) (context.Context, log.Fields) {
		if ti == nil {
			return parentCtx, logFields
		}
		ctx := context.WithValue(parentCtx, AppInfoKey, *ti)
		logFields = ti.AppendToLogFields(logFields)
		return ctx, logFields
	}
}
