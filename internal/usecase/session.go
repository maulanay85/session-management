package usecase

import (
	"context"
	"fmt"
	"scs-session/internal/config"
	"scs-session/internal/domain"
	"scs-session/internal/repository"
	"time"

	"github.com/alexedwards/scs/v2"
)

type SessionUsecaseImpl struct {
	conf config.Config
	sr   repository.SessionRepository
	sm   scs.SessionManager
}

// GetByToken implements SessionUsecase.
func (s *SessionUsecaseImpl) GetByToken(ctx context.Context, token string) (domain.Session, error) {
	data, err := s.sr.GetByToken(ctx, token)
	if err != nil {
		return domain.Session{}, err
	}
	return data, nil
}

type SessionUsecase interface {
	Validate(ctx context.Context, token string) (domain.Session, error)
	GetByToken(ctx context.Context, token string) (domain.Session, error)
}

func NewSessionUsecase(conf config.Config, sr repository.SessionRepository, sm scs.SessionManager) SessionUsecase {
	return &SessionUsecaseImpl{
		conf: conf,
		sr:   sr,
		sm:   sm,
	}
}

// Validate implements SessionUsecase.
func (s *SessionUsecaseImpl) Validate(ctx context.Context, token string) (domain.Session, error) {
	data, err := s.sr.GetByToken(ctx, fmt.Sprintf("token:%s", token))
	if err != nil {
		return domain.Session{}, err
	}
	now := time.Now()
	newExpiry := now.Add(time.Duration(s.conf.SessionMaxIdleTime) * time.Minute)
	data.ExpiredAt = newExpiry
	_, err = s.sr.UpdateSession(ctx, data)
	if err != nil {
		return domain.Session{}, err
	}
	return data, err
}
