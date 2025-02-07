package models

import (
	models_time "github.com/org/2112-space-lab/org/testing/pkg/testing/x-test-time/models"
)

type EventCallbackActionName string

type EventCallbackInfo struct {
	EventType         string
	Action            EventCallbackActionName
	ActionHandlerArgs string
	ActionDelay       string
}

type ExpectedEvent struct {
	EventType                  string
	Occurence                  int
	FromTime                   models_time.TimeCheckpointExpression
	ToTimeWarn                 models_time.TimeCheckpointExpression
	ToTimeErr                  models_time.TimeCheckpointExpression
	ProduceCheckpointEventTime string
	IsReject                   bool
	AssignRef                  NamedAppEventReference // used to give a known name to matching event for accessing it afterwards
	XPathQuery                 string
	XPathValue                 string
}

type ReferencedEventExpectedInfo struct {
	AssignedRef NamedAppEventReference
	XPathQuery  string
	XPathValue  string
	Description string
}
