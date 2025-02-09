package services

import (
	"context"
	"encoding/json"

	"github.com/org/2112-space-lab/org/app-service/internal/domain"
	log "github.com/org/2112-space-lab/org/app-service/pkg/log"
	fx "github.com/org/2112-space-lab/org/app-service/pkg/option"
	xtime "github.com/org/2112-space-lab/org/app-service/pkg/time"
	"github.com/org/2112-space-lab/org/app-service/pkg/tracing"
)

// AuditTrailService provides business logic for managing audit trails.
type AuditTrailService struct {
	repo domain.AuditTrailRepository
}

// NewAuditTrailService creates a new instance of AuditTrailService.
func NewAuditTrailService(repo domain.AuditTrailRepository) AuditTrailService {
	return AuditTrailService{repo: repo}
}

// Create creates a new audit trail record.
func (s *AuditTrailService) Create(ctx context.Context, auditTrail domain.AuditTrail) (err error) {
	ctx, span := tracing.NewSpan(ctx, "AuditTrailService.Create")
	defer span.EndWithError(err)

	err = s.repo.Save(ctx, auditTrail)
	if err != nil {
		log.WithFields(log.Fields{
			"func": "AuditTrailService.Create",
		}).WithError(err).Error("failed to create audit trail")
		return err
	}

	return nil
}

// GetByRecordIDAndTable retrieves audit trails by record ID and table name.
func (s *AuditTrailService) GetByRecordIDAndTable(ctx context.Context, tableName, recordID string) (trails []domain.AuditTrail, err error) {
	ctx, span := tracing.NewSpan(ctx, "AuditTrailService.GetByRecordIDAndTable")
	defer span.EndWithError(err)

	trails, err = s.repo.FindByRecordIDAndTable(ctx, tableName, recordID)
	if err != nil {
		log.WithFields(log.Fields{
			"func":      "AuditTrailService.GetByRecordIDAndTable",
			"tableName": tableName,
			"recordID":  recordID,
		}).WithError(err).Error("failed to fetch audit trails")
		return nil, err
	}

	return trails, nil
}

// GetAllWithPagination retrieves all audit trails with pagination support.
func (s *AuditTrailService) GetAllWithPagination(ctx context.Context, page, pageSize int) (trails []domain.AuditTrail, total int64, err error) {
	ctx, span := tracing.NewSpan(ctx, "AuditTrailService.GetAllWithPagination")
	defer span.EndWithError(err)

	trails, total, err = s.repo.FindAllWithPagination(ctx, page, pageSize)
	if err != nil {
		log.WithFields(log.Fields{
			"func":     "AuditTrailService.GetAllWithPagination",
			"page":     page,
			"pageSize": pageSize,
		}).WithError(err).Error("failed to fetch paginated audit trails")
		return nil, 0, err
	}

	return trails, total, nil
}

// LogAudit logs an audit trail entry.
func (s *AuditTrailService) LogAudit(ctx context.Context, tableName, recordID, action, performedBy string, changes map[string]interface{}) (err error) {
	ctx, span := tracing.NewSpan(ctx, "AuditTrailService.LogAudit")
	defer span.EndWithError(err)

	// Serialize changes to JSON (optional field)
	var changesJSON string
	if changes != nil {
		changesJSONBytes, marshalErr := json.Marshal(changes)
		if marshalErr != nil {
			log.WithFields(log.Fields{
				"func": "AuditTrailService.LogAudit",
			}).WithError(marshalErr).Error("failed to marshal changes for audit trail")
			return marshalErr
		}
		changesJSON = string(changesJSONBytes)
	}

	auditTrail := domain.AuditTrail{
		TableName:   tableName,
		RecordID:    recordID,
		Action:      domain.ActionType(action),
		ChangesJSON: fx.NewValueOption(changesJSON),
		PerformedBy: domain.UserID(performedBy),
		PerformedAt: xtime.UtcNow(),
	}

	return s.Create(ctx, auditTrail)
}
