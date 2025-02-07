package models

import (
	"encoding/json"
	"time"
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
	EventTimeUtc time.Time
	EventUid     string
	EventType    string
	Comment      string
	Payload      AppEventRawJSON
}
