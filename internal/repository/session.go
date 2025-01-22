package repository

import (
	"context"
	"scs-session/internal/config"
	"scs-session/internal/domain"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type SessionRepositoryImpl struct {
	db   *gorm.DB
	r    *redis.Client
	conf config.Config
}

type SessionRepository interface {
	InsertSession(ctx context.Context, data domain.Session) error
	UpdateSession(ctx context.Context, data domain.Session) (domain.Session, error)
	GetByToken(ctx context.Context, token string) (domain.Session, error)
}

func NewSessionRepository(db *gorm.DB, r *redis.Client, conf config.Config) SessionRepository {
	return &SessionRepositoryImpl{
		db:   db,
		r:    r,
		conf: conf,
	}
}

// InsertSession implements SessionRepository.
func (s *SessionRepositoryImpl) InsertSession(ctx context.Context, data domain.Session) error {
	if err := s.r.Set(ctx, "token:"+data.Token, data.UserID, time.Until(data.ExpiredAt)).Err(); err != nil {
		return err
	}
	return nil
}

// UpdateSession implements SessionRepository.
func (s *SessionRepositoryImpl) UpdateSession(ctx context.Context, data domain.Session) (domain.Session, error) {
	if err := s.r.Expire(ctx, data.Token, time.Until(data.ExpiredAt)).Err(); err != nil {
		return domain.Session{}, err
	}
	return data, nil
}

// GetByToken implements SessionRepository.
func (s *SessionRepositoryImpl) GetByToken(ctx context.Context, token string) (domain.Session, error) {
	result, err := s.r.Get(ctx, token).Result()
	if err != nil {
		return domain.Session{}, err
	}
	return domain.Session{
		Token:  token,
		UserID: result,
	}, nil

}
