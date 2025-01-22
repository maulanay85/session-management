package domain

import "time"

// deprecated
type Session struct {
	ID        string    `gorm:"primaryKey;column:id"`
	UserID    string    `gorm:"column:user_id"`
	Token     string    `gorm:"column:token"`
	IsValid   bool      `gorm:"column:is_valid"`
	ExpiredAt time.Time `gorm:"expired_at"`
}
