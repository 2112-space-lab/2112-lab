package domain

import (
	"context"

	fx "github.com/org/2112-space-lab/org/app-service/pkg/option"
	xtime "github.com/org/2112-space-lab/org/app-service/pkg/time"
)

// ActionType represents the type of action performed in the audit trail.
type ActionType string

// UserID represents the identifier of the user or system that performed the action.
type UserID string

// AuditTrail represents an audit trail record in the domain layer.
type AuditTrail struct {
	ModelBase
	TableName   string            // Name of the table where the action occurred
	RecordID    string            // ID of the affected record
	Action      ActionType        // Action performed (e.g., INSERT, UPDATE, DELETE)
	ChangesJSON fx.Option[string] // JSON representation of the changes made (optional)
	PerformedBy UserID            // User or system that performed the action
	PerformedAt xtime.UtcTime     // Time the action was performed
}

// AuditTrailRepository defines the interface for audit trail operations.
type AuditTrailRepository interface {
	// Save a new audit trail record
	Save(ctx context.Context, auditTrail AuditTrail) error

	// Find by record ID and table name
	FindByRecordIDAndTable(ctx context.Context, tableName string, recordID string) ([]AuditTrail, error)

	// Retrieve all audit trails with pagination
	FindAllWithPagination(ctx context.Context, page int, pageSize int) ([]AuditTrail, int64, error)
}
