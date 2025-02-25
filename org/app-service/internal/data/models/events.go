package models

import (
	"time"

	"github.com/org/2112-space-lab/org/app-service/internal/domain"
	domainenum "github.com/org/2112-space-lab/org/app-service/internal/domain/domain-enums"
	fx "github.com/org/2112-space-lab/org/app-service/pkg/option"
	xtime "github.com/org/2112-space-lab/org/app-service/pkg/time"
)

// Event represents the database model for events.
type Event struct {
	ModelBase
	EventType   string    `gorm:"size:255;not null;index"` // Event type (e.g., "SATELLITE_TLE_PROPAGATED")
	EventUID    string    `gorm:"size:255;not null;unique"`
	Payload     string    `gorm:"type:json;not null"` // Event payload in JSON format
	PublishedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
	Comment     *string   `gorm:"type:text"` // Optional event comment
}

// EventHandlerLog tracks when an event handler starts, ends, and the event that triggered it.
type EventHandler struct {
	ModelBase
	EventID     string     `gorm:"size:255;not null;index"` // Reference to the event that triggered the handler
	HandlerName string     `gorm:"size:255;not null"`       // Handler name (e.g., "RehydrateGameContextHandler")
	StartedAt   time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP"`
	CompletedAt *time.Time `gorm:"null"`             // Completion timestamp (null if still running)
	Status      string     `gorm:"size:50;not null"` // Status: "STARTED", "COMPLETED", "FAILED"
	Error       *string    `gorm:"type:text"`        // Error message if failed

	Event Event `gorm:"constraint:OnDelete:CASCADE;foreignKey:EventID;references:ID"`
}

// MapToEventDomain converts an Event database model to a domain Event model.
func MapToEventDomain(e Event) domain.Event {
	return domain.Event{
		ModelBase: domain.ModelBase{
			ID:          e.ID,
			CreatedAt:   e.CreatedAt,
			UpdatedAt:   &e.UpdatedAt,
			DeleteAt:    e.DeleteAt,
			ProcessedAt: e.ProcessedAt,
			IsActive:    e.IsActive,
			IsFavourite: e.IsFavourite,
			DisplayName: e.DisplayName,
		},
		EventType:   domain.EventType(e.EventType),
		EventUID:    e.EventUID,
		Payload:     fx.AsOption(&e.Payload),
		PublishedAt: xtime.NewUtcTimeIgnoreZone(e.PublishedAt),
		Comment:     fx.AsOption(e.Comment),
	}
}

// MapToEventModel converts a domain Event model to an Event database model.
func MapToEventModel(e domain.Event) Event {
	return Event{
		ModelBase: ModelBase{
			ID:          e.ModelBase.ID,
			CreatedAt:   e.ModelBase.CreatedAt,
			UpdatedAt:   *e.ModelBase.UpdatedAt,
			DeleteAt:    e.ModelBase.DeleteAt,
			ProcessedAt: e.ModelBase.ProcessedAt,
			IsActive:    e.ModelBase.IsActive,
			IsFavourite: e.ModelBase.IsFavourite,
			DisplayName: e.ModelBase.DisplayName,
		},
		EventType:   string(e.EventType),
		EventUID:    e.EventUID,
		Payload:     fx.GetOrDefault(e.Payload, ""),
		PublishedAt: e.PublishedAt.Inner(),
		Comment:     fx.ConvertToStrPtr(e.Comment),
	}
}

// MapToEventHandlerDomain converts an EventHandlerLog database model to a domain EventHandlerLog model.
func MapToEventHandlerDomain(h EventHandler) (domain.EventHandler, error) {

	state, err := domainenum.PotentialHandlerState(h.Status).Validate()
	if err != nil {
		return domain.EventHandler{}, err
	}

	return domain.EventHandler{
		ModelBase: domain.ModelBase{
			ID:          h.ID,
			CreatedAt:   h.CreatedAt,
			UpdatedAt:   &h.UpdatedAt,
			DeleteAt:    h.DeleteAt,
			ProcessedAt: h.ProcessedAt,
			IsActive:    h.IsActive,
			IsFavourite: h.IsFavourite,
			DisplayName: h.DisplayName,
		},
		EventID:     h.EventID,
		HandlerName: h.HandlerName,
		StartedAt:   xtime.NewUtcTimeIgnoreZone(h.StartedAt),
		CompletedAt: xtime.ConvertToUtcTime(h.CompletedAt),
		Status:      state,
		Error:       fx.AsOption(h.Error),
	}, nil
}

// MapToEventHandlerModel converts a domain EventHandlerLog model to an EventHandler database model.
func MapToEventHandlerModel(h domain.EventHandler) EventHandler {
	return EventHandler{
		ModelBase: ModelBase{
			ID:          h.ModelBase.ID,
			CreatedAt:   h.ModelBase.CreatedAt,
			UpdatedAt:   *h.ModelBase.UpdatedAt,
			DeleteAt:    h.ModelBase.DeleteAt,
			ProcessedAt: h.ModelBase.ProcessedAt,
			IsActive:    h.ModelBase.IsActive,
			IsFavourite: h.ModelBase.IsFavourite,
			DisplayName: h.ModelBase.DisplayName,
		},
		EventID:     h.EventID,
		HandlerName: h.HandlerName,
		StartedAt:   h.StartedAt.Inner(),
		CompletedAt: xtime.ConvertToTimePtr(h.CompletedAt),
		Status:      h.Status.String(),
		Error:       fx.ConvertToStrPtr(h.Error),
	}
}
