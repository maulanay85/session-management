package usecase

import (
	"context"
	"encoding/json"
	"log"
	"scs-session/internal/domain"
	"scs-session/internal/dto"
	"scs-session/internal/repository"
)

type UserUsecaseImpl struct {
	ur        repository.UserRepository
	nsqCLient repository.NSQRepository
}

// Update implements UserUsecase.
func (u *UserUsecaseImpl) Update(ctx context.Context, id string, data dto.UserUpdateRequest) error {
	exist, err := u.ur.FindById(id)
	if err != nil {
		return err
	}
	exist.Address = data.Address
	exist.Email = data.Email
	exist.FullName = data.FullName
	_, err = u.ur.Update(exist)
	if err != nil {
		return err
	}

	rawMessage, _ := json.Marshal(exist)
	userId := ctx.Value("id").(string)
	go func(id string, rawMessage []byte) {
		err := u.nsqCLient.PublishMessage(domain.AuditTrail{
			ServiceName: "scs-session",
			TableName:   "user",
			EntityID:    id,
			CreatedBy:   userId,
			Message:     "data telah diperbaharui",
			Data:        string(rawMessage),
		}, "audit_trail")

		if err != nil {
			log.Printf("error when publish message : %+v", err)
		}
	}(id, rawMessage)
	return nil
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
	Update(ctx context.Context, id string, data dto.UserUpdateRequest) error
}

func NewUserUsecase(ur repository.UserRepository, nsqClient repository.NSQRepository) UserUsecase {
	return &UserUsecaseImpl{
		ur:        ur,
		nsqCLient: nsqClient,
	}
}
