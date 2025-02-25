package domain

import (
	"time"

	"github.com/google/uuid"
	xtime "github.com/org/2112-space-lab/org/app-service/pkg/time"
)

type ModelBase struct {
	ID          string
	DisplayName string
	CreatedAt   time.Time
	UpdatedAt   *time.Time
	DeleteAt    *time.Time
	ProcessedAt *time.Time
	IsActive    bool
	IsFavourite bool
}

// NewModelBase creates a new instance.
func NewModelBase(name string, spaceID string, satType SatelliteType, isFavourite bool, isActive bool, createdAt time.Time) (ModelBase, error) {
	if err := satType.IsValid(); err != nil {
		return ModelBase{}, err
	}
	return ModelBase{
		ID:          uuid.NewString(),
		CreatedAt:   createdAt,
		UpdatedAt:   &createdAt,
		DisplayName: name,
		IsActive:    isActive,
		ProcessedAt: &createdAt,
		IsFavourite: isFavourite,
	}, nil
}

func NewModelBaseDefault() ModelBase {

	nowUtc := xtime.UtcNow().Inner()
	return ModelBase{
		ID:          uuid.NewString(),
		CreatedAt:   nowUtc,
		UpdatedAt:   &nowUtc,
		DisplayName: "default",
		IsActive:    true,
		ProcessedAt: &nowUtc,
		IsFavourite: false,
	}
}
