package models

import (
	"encoding/json"

	xtime "github.com/org/2112-space-lab/org/testing/pkg/x-time"
)

type ServiceName string

type GlobalPropKeyValueMap map[string]string

type ServiceAppName string

const (
	AppName           ServiceAppName = "app"
	PropagatorAppName ServiceAppName = "propagator"
)

type NamedAppEventReference string
type AppEventRawJSON json.RawMessage

type EventRoot struct {
	EventTimeUtc string          `json:"event_time_utc"`    // ISO-8601 string format (matches Python)
	EventUid     string          `json:"event_uid"`         // UUID string
	EventType    string          `json:"event_type"`        // Event type as string
	Comment      string          `json:"comment,omitempty"` // Optional field
	Payload      json.RawMessage `json:"payload,omitempty"` // Raw JSON for flexible structure
}

func (e *EventRoot) GetEventTimeUtc() xtime.UtcTime {
	v, err := xtime.FromString(xtime.DateTimeFormat(e.EventTimeUtc))
	if err != nil {
		panic("error while parsing date from event")
	}
	return v
}
