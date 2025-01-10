package usecase

import (
	"context"
	"scs-session/internal/dto"
	"scs-session/internal/repository"
)

type UserUsecaseImpl struct {
	ur repository.UserRepository
}

// FindById implements UserUsecase.
func (u *UserUsecaseImpl) FindById(ctx context.Context) (dto.BaseResponse, error) {
	id := ctx.Value("id").(string)
	data, err := u.ur.FindById(id)
	if err != nil {
		return dto.BaseResponse{}, err
	}
	return dto.BaseResponse{Data: data}, nil
}

type UserUsecase interface {
	FindById(ctx context.Context) (dto.BaseResponse, error)
}

func NewUserUsecase(ur repository.UserRepository) UserUsecase {
	return &UserUsecaseImpl{
		ur: ur,
	}
}
