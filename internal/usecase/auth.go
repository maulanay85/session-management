package usecase

import (
	"context"
	"errors"
	"fmt"
	"scs-session/internal/config"
	"scs-session/internal/domain"
	"scs-session/internal/dto"
	"scs-session/internal/helper"
	"scs-session/internal/repository"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/redis/go-redis/v9"
)

type AuthUsecaseImpl struct {
	conf              config.Config
	userRepository    repository.UserRepository
	sessionRepository repository.SessionRepository
	sessionManager    scs.SessionManager
	redisClient       *redis.Client
	helper            helper.UtilInterface
}

type AuthUsecase interface {
	Login(ctx context.Context, req dto.LoginRequest) (dto.LoginResponse, error)
	IncrementFailedAttemps(ctx context.Context, userId string) (int64, error)
	LockUser(ctx context.Context, userId string) error
	Logout(ctx context.Context) error
}

func NewAuthUseCase(conf config.Config, userRepository repository.UserRepository, sessionRepository repository.SessionRepository, sessionManager scs.SessionManager, redisClient *redis.Client, helper helper.UtilInterface) AuthUsecase {
	return &AuthUsecaseImpl{
		conf:              conf,
		userRepository:    userRepository,
		sessionRepository: sessionRepository,

		sessionManager: sessionManager,
		redisClient:    redisClient,

		helper: helper,
	}
}

func (a *AuthUsecaseImpl) Login(ctx context.Context, req dto.LoginRequest) (dto.LoginResponse, error) {
	data, err := a.userRepository.GetByEmail(req.Email)
	if err != nil {
		fmt.Printf("error due to: %v", err)
		return dto.LoginResponse{}, errors.New("username or password is wrong")
	}
	val, err := a.redisClient.Exists(ctx, "lockout:"+data.ID).Result()
	if err != nil {
		return dto.LoginResponse{}, err
	}
	if val > 0 {
		return dto.LoginResponse{}, errors.New("your id has been locked")
	}
	if data.Password != req.Password {

		failTime, err := a.IncrementFailedAttemps(ctx, data.ID)
		if err != nil {
			return dto.LoginResponse{}, err
		}

		if failTime >= int64(a.conf.LoginMaxTry) {
			a.redisClient.Del(ctx, "failed_attempts:"+data.ID)
			a.LockUser(ctx, data.ID)
		}

		return dto.LoginResponse{}, errors.New("username or password is wrong")
	}

	a.sessionManager.Put(ctx, "session-id", data.ID)
	blankToken := a.helper.GenerateBlankToken()
	blankTokenExpired := time.Now().Add(time.Duration(a.conf.SessionMaxIdleTime) * time.Minute)
	err = a.sessionRepository.InsertSession(ctx, domain.Session{
		UserID:    data.ID,
		Token:     blankToken,
		ExpiredAt: blankTokenExpired,
	})
	if err != nil {
		return dto.LoginResponse{}, err
	}
	response := dto.LoginResponse{
		ID:                data.ID,
		FullName:          data.FullName,
		Email:             data.Email,
		BlankToken:        blankToken,
		BlankTokenExpired: blankTokenExpired,
	}
	return response, nil
}

func (a *AuthUsecaseImpl) IncrementFailedAttemps(ctx context.Context, userId string) (int64, error) {
	key := "failed_attempts:" + userId
	failedAttemps, err := a.redisClient.Incr(ctx, key).Result()
	return failedAttemps, err
}

func (a *AuthUsecaseImpl) LockUser(ctx context.Context, userId string) error {
	key := "lockout:" + userId
	expiration := time.Duration(a.conf.LoginFinaltyTime) * time.Minute
	return a.redisClient.Set(ctx, key, "true", expiration).Err()
}

func (a *AuthUsecaseImpl) Logout(ctx context.Context) error {
	err := a.sessionManager.Destroy(ctx)
	if err != nil {
		return err
	}
	return nil
}
