package repository

import (
	"context"

	"github.com/org/2112-space-lab/org/app-service/internal/data"
	"github.com/org/2112-space-lab/org/app-service/internal/data/models"
	"github.com/org/2112-space-lab/org/app-service/internal/domain"
)

// AuditTrailRepository manages audit trail data access.
type AuditTrailRepository struct {
	db *data.Database
}

// NewAuditTrailRepository creates a new AuditTrailRepository instance.
func NewAuditTrailRepository(db *data.Database) AuditTrailRepository {
	return AuditTrailRepository{db: db}
}

// Save creates a new audit trail record.
func (r *AuditTrailRepository) Save(ctx context.Context, auditTrail domain.AuditTrail) error {
	model := models.MapToAuditTrailModel(auditTrail)
	return r.db.DbHandler.Create(&model).Error
}

// FindByRecordIDAndTable retrieves audit trail records by record ID and table name.
func (r *AuditTrailRepository) FindByRecordIDAndTable(ctx context.Context, tableName, recordID string) ([]domain.AuditTrail, error) {
	var results []models.AuditTrail
	err := r.db.DbHandler.Where("table_name = ? AND record_id = ?", tableName, recordID).Find(&results).Error
	if err != nil {
		return nil, err
	}

	var auditTrails []domain.AuditTrail
	for _, result := range results {
		auditTrails = append(auditTrails, models.MapToAuditTrailDomain(result))
	}
	return auditTrails, nil
}

// FindAllWithPagination retrieves all audit trails with pagination.
func (r *AuditTrailRepository) FindAllWithPagination(ctx context.Context, page, pageSize int) ([]domain.AuditTrail, int64, error) {
	var results []models.AuditTrail
	var totalRecords int64

	// Calculate offset for pagination
	offset := (page - 1) * pageSize

	query := r.db.DbHandler.Table("audit_trails")

	// Count total records before applying limit and offset
	if err := query.Count(&totalRecords).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination and retrieve results
	if err := query.Limit(pageSize).Offset(offset).Find(&results).Error; err != nil {
		return nil, 0, err
	}

	// Map results to domain.AuditTrail
	var auditTrails []domain.AuditTrail
	for _, result := range results {
		auditTrails = append(auditTrails, models.MapToAuditTrailDomain(result))
	}

	return auditTrails, totalRecords, nil
}
