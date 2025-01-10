package repository

import (
	"scs-session/internal/domain"

	"gorm.io/gorm"
)

type SessionRepositoryImpl struct {
	db *gorm.DB
}

type SessionRepository interface {
	InsertSession(data domain.Session) (domain.Session, error)
	UpdateSession(data domain.Session) (domain.Session, error)
	GetByToken(token string) (domain.Session, error)
}

func NewSessionRepository(db *gorm.DB) SessionRepository {
	return &SessionRepositoryImpl{
		db: db,
	}
}

// InsertSession implements SessionRepository.
func (s *SessionRepositoryImpl) InsertSession(data domain.Session) (domain.Session, error) {
	result := s.db.Create(&data)
	if result.Error != nil {
		return domain.Session{}, result.Error
	}
	return data, nil
}

// UpdateSession implements SessionRepository.
func (s *SessionRepositoryImpl) UpdateSession(data domain.Session) (domain.Session, error) {
	result := s.db.Save(data)
	if result.Error != nil {
		return domain.Session{}, result.Error
	}
	return data, nil
}

// GetByToken implements SessionRepository.
func (s *SessionRepositoryImpl) GetByToken(token string) (domain.Session, error) {
	var session domain.Session
	if err := s.db.Where("token = ?", token).First(&session).Error; err != nil {
		return domain.Session{}, err
	}
	return session, nil
}
