package models

import "encoding/json"

type ServiceName string

type GlobalPropKeyValueMap map[string]string

type ServiceAppName string

const (
	AppName           ServiceAppName = "app"
	PropagatorAppName ServiceAppName = "propagator"
)

type NamedAppEventReference string
type AppEventRawJSON json.RawMessage
