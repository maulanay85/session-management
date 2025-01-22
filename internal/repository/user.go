package repository

import (
	"fmt"
	"scs-session/internal/domain"

	"gorm.io/gorm"
)

type UserRepositoryImpl struct {
	db *gorm.DB
}

// Update implements UserRepository.
func (u *UserRepositoryImpl) Update(data domain.User) (domain.User, error) {
	result := u.db.Save(&data)
	if result.Error != nil {
		return domain.User{}, result.Error
	}
	return data, nil
}

// FindById implements UserRepository.
func (u *UserRepositoryImpl) FindById(id string) (domain.User, error) {
	var user domain.User
	err := u.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

// GetByEmail implements UserRepository.
func (u *UserRepositoryImpl) GetByEmail(email string) (domain.User, error) {
	var user domain.User
	err := u.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return domain.User{}, err
	}
	if user.Email == "" {
		return domain.User{}, fmt.Errorf("data not found")
	}
	return user, nil

}

type UserRepository interface {
	GetByEmail(email string) (domain.User, error)
	FindById(id string) (domain.User, error)
	Update(domain.User) (domain.User, error)
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &UserRepositoryImpl{
		db: db,
	}
}
