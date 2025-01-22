package usecase

import (
	"context"
	"scs-session/internal/domain"
	"scs-session/internal/repository"

	"github.com/google/uuid"
)

type AuditTrailUsecase interface {
	HandleAuditTrailMessage(ctx context.Context, message domain.AuditTrail) error
	GetAll(ctx context.Context, filter domain.AuditTrailFilter) ([]domain.AuditTrail, error)
}

type AuditTrailUsecaseImpl struct {
	auditRepository repository.AuditTrailRepository
	nsqRepository   repository.NSQRepository
}

// GetAll implements AuditTrailUsecase.
func (a AuditTrailUsecaseImpl) GetAll(ctx context.Context, filter domain.AuditTrailFilter) ([]domain.AuditTrail, error) {
	data, err := a.auditRepository.GetAll(ctx, filter)
	if err != nil {
		return []domain.AuditTrail{}, err
	}
	return data, nil
}

// HandleAuditTrailMessage implements AuditTrailUsecase.
func (a AuditTrailUsecaseImpl) HandleAuditTrailMessage(ctx context.Context, message domain.AuditTrail) error {
	message.ID = uuid.New().String()
	err := a.auditRepository.Insert(ctx, message)
	if err != nil {
		return err
	}
	return nil
}

func NewAuditTrailUsecase(auditRepository repository.AuditTrailRepository, nsqRepository repository.NSQRepository) AuditTrailUsecase {
	return AuditTrailUsecaseImpl{
		auditRepository: auditRepository,
		nsqRepository:   nsqRepository,
	}
}
