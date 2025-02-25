package domain

import (
	domainenum "github.com/org/2112-space-lab/org/app-service/internal/domain/domain-enums"
	fx "github.com/org/2112-space-lab/org/app-service/pkg/option"
	xtime "github.com/org/2112-space-lab/org/app-service/pkg/time"
)

// EventType represents the type of an event.
type EventType string

// Event represents a domain model for an event.
type Event struct {
	ModelBase
	EventType   EventType
	EventUID    string
	Payload     fx.Option[string]
	PublishedAt xtime.UtcTime
	Comment     fx.Option[string]
}

// EventHandlerLog represents the execution log of an event handler.
type EventHandler struct {
	ModelBase
	EventID     string
	HandlerName string
	StartedAt   xtime.UtcTime
	CompletedAt fx.Option[xtime.UtcTime]
	Status      domainenum.HandlerState
	Error       fx.Option[string]
}
