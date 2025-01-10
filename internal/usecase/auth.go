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
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
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
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       data.ID,
		"email":    data.Email,
		"fullName": data.FullName,
		"exp":      time.Now().Add(time.Duration(a.conf.TokenExpiry) * time.Minute).Unix(),
	})

	tokenString, err := token.SignedString([]byte(a.conf.TokenSecretKey))
	if err != nil {
		fmt.Printf("error due to: %v", err)
		return dto.LoginResponse{}, errors.New("failed to generate access token")
	}

	// Generate Refresh token
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    data.ID,
		"email": data.Email,
		"exp":   time.Now().Add(time.Duration(a.conf.TokenRefreshExpiry) * time.Minute).Unix(),
	})

	refreshTokenString, err := refreshToken.SignedString([]byte(a.conf.TokenSecretKey))
	if err != nil {
		fmt.Printf("error due to: %v", err)
		return dto.LoginResponse{}, errors.New("failed to generate refresh token")
	}

	a.sessionManager.Put(ctx, "session-id", data.ID)
	blankToken := a.helper.GenerateBlankToken()
	blankTokenExpired := time.Now().Add(time.Duration(a.conf.SessionMaxIdleTime) * time.Minute)
	a.sessionRepository.InsertSession(domain.Session{
		ID:        uuid.New().String(),
		UserID:    data.ID,
		Token:     blankToken,
		ExpiredAt: blankTokenExpired,
		IsValid:   true,
	})
	response := dto.LoginResponse{
		AccessToken:       tokenString,
		RefreshToken:      refreshTokenString,
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
