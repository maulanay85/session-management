package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	ServiceEnv         string
	ServiceName        string
	ServicePort        string
	RedisHost          string
	RedisPort          string
	RedisUsername      string
	RedisPassword      string
	TokenExpiry        int
	TokenRefreshExpiry int
	TokenSecretKey     string
	DatabaseHost       string
	DatabasePort       int
	DatabaseName       string
	DatabaseUser       string
	DatabasePassword   string
	SessionMaxIdleTime int
	LoginFinaltyTime   int64
	LoginMaxTry        int
	NsqTopic           string
	NsqChannel         string
	NsqUrl             string
}

func InitializeConfig() (*Config, error) {
	viper.AutomaticEnv()
	viper.SetConfigFile(".env")

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("cant read config due to: %v", err)
		return nil, err
	}

	config := &Config{
		ServiceEnv:         viper.GetString("service.env"),
		ServiceName:        viper.GetString("service.name"),
		ServicePort:        viper.GetString("service.port"),
		RedisHost:          viper.GetString("redis.host"),
		RedisPort:          viper.GetString("redis.port"),
		RedisUsername:      viper.GetString("redis.username"),
		RedisPassword:      viper.GetString("redis.password"),
		TokenExpiry:        viper.GetInt("token.expiry"),
		TokenRefreshExpiry: viper.GetInt("token.refresh.expiry"),
		TokenSecretKey:     viper.GetString("token.secret.key"),
		DatabaseHost:       viper.GetString("database.host"),
		DatabasePort:       viper.GetInt("database.port"),
		DatabaseName:       viper.GetString("database.name"),
		DatabaseUser:       viper.GetString("database.user"),
		DatabasePassword:   viper.GetString("database.password"),
		SessionMaxIdleTime: viper.GetInt("session.max.idle.time"),
		LoginMaxTry:        viper.GetInt("login.max.try"),
		LoginFinaltyTime:   viper.GetInt64("login.finalty.time"),
		NsqTopic:           viper.GetString("nsq.topic"),
		NsqChannel:         viper.GetString("nsq.channel"),
		NsqUrl:             viper.GetString("nsq.url"),
	}
	return config, nil
}
