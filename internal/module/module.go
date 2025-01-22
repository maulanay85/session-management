package module

import (
	"fmt"
	"scs-session/internal/config"
	"scs-session/internal/consumer"
	"scs-session/internal/controller"
	"scs-session/internal/helper"
	"scs-session/internal/middleware"
	"scs-session/internal/repository"
	"scs-session/internal/usecase"
	"time"

	scsredis "github.com/alexedwards/scs/redisstore"
	"github.com/alexedwards/scs/v2"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
)

func Init(conf config.Config) *gin.Engine {
	r := gin.New()
	r.Use(cors.Default())
	r.Use(gin.Recovery())

	scsRedisPool := &redis.Pool{
		MaxIdle: 10,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", fmt.Sprintf("%s:%s", conf.RedisHost, conf.RedisPort))
		},
	}

	sessionManager := scs.New()
	sessionManager.Store = scsredis.New(scsRedisPool)
	sessionManager.Lifetime = time.Duration(conf.TokenExpiry) * time.Minute
	sessionManager.Cookie.Name = "UserSession"
	sessionManager.Cookie.HttpOnly = true
	sessionManager.Cookie.Secure = false
	sessionManager.Cookie.Path = "/"
	sessionManager.Cookie.Domain = "localhost"

	// helper
	util := helper.NewUtil()

	// database
	db := config.InitializeDatabase(conf)

	// message broker
	nsq, err := config.InitializeNSQ(conf)
	if err != nil {
		panic(err)
	}

	// redis
	redisClient := config.InitializeRedis(&conf)

	// repository
	auditTrailRepository := repository.NewAuditTrailRepository(db)
	userRepository := repository.NewUserRepository(db)
	sessionRepository := repository.NewSessionRepository(db, redisClient, conf)
	nsqRepository := repository.NewNSQRepository(&nsq)

	// usecase
	auditTrailUsecase := usecase.NewAuditTrailUsecase(auditTrailRepository, nsqRepository)
	authUsecase := usecase.NewAuthUseCase(conf, userRepository, sessionRepository, *sessionManager, redisClient, util)
	sessionUsecase := usecase.NewSessionUsecase(conf, sessionRepository, *sessionManager)
	userUsecase := usecase.NewUserUsecase(userRepository, nsqRepository)

	// controller
	authController := controller.NewAuthController(authUsecase)
	userController := controller.NewUserController(userUsecase)
	auditTrailController := controller.NewAuditTrailController(auditTrailUsecase)

	go func() {
		ansq := consumer.NewNSQConsumer(auditTrailUsecase)
		consumer.StartNSQConsumer(conf.NsqUrl, "audit_trail", ansq)
	}()

	r.Use(middleware.LoadAndSave(sessionManager))
	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, "pong!")
	})
	authGroup := r.Group("/auth")
	{
		authGroup.POST("", authController.Login)
	}
	v1 := r.Group("/v1")
	v1.Use(middleware.SessionMiddleware(conf, sessionManager, sessionUsecase))
	{
		v1.GET("/profile", userController.GetProfile)
		v1.PUT("/user/:id", userController.Update)
		v1.GET("/audit-trail", auditTrailController.GetAll)
	}
	return r
}
