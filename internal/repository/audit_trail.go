package repository

import (
	"context"
	"scs-session/internal/domain"

	"gorm.io/gorm"
)

type AuditTrailRepositoryImpl struct {
	db *gorm.DB
}

// GetAll implements AuditTrailRepository.
func (a AuditTrailRepositoryImpl) GetAll(ctx context.Context, filter domain.AuditTrailFilter) ([]domain.AuditTrail, error) {
	var data []domain.AuditTrail
	query := a.db.WithContext(ctx)
	if filter.EntityID != "" && filter.ServiceName != "" {
		query = query.Where("entity_id = ? and service_name = ?", filter.EntityID, filter.ServiceName)
	}
	if filter.ID != "" {
		query = query.Where("id = ?", filter.ID)
	}
	query.Order("created_at desc")
	if err := query.Find(&data).Error; err != nil {
		return data, err
	}
	return data, nil
}

// Insert implements AuditTrailRepository.
func (a AuditTrailRepositoryImpl) Insert(ctx context.Context, data domain.AuditTrail) error {
	query := a.db.WithContext(ctx)
	result := query.Save(&data)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

type AuditTrailRepository interface {
	Insert(ctx context.Context, data domain.AuditTrail) error
	GetAll(ctx context.Context, filter domain.AuditTrailFilter) ([]domain.AuditTrail, error)
}

func NewAuditTrailRepository(db *gorm.DB) AuditTrailRepository {
	return AuditTrailRepositoryImpl{
		db: db,
	}
}
