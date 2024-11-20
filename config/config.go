package config

import (
	"fmt"
	"os"

	"github.com/Moranilt/http-utils/clients/database"
	"github.com/Moranilt/http-utils/clients/rabbitmq"
	"github.com/Moranilt/http-utils/clients/redis"
	"github.com/spf13/viper"
)

const (
	ENV_PRODUCTION = "PRODUCTION"
	ENV_PORT       = "PORT"

	ENV_DB_NAME     = "DB_NAME"
	ENV_DB_HOST     = "DB_HOST"
	ENV_DB_USER     = "DB_USER"
	ENV_DB_PASSWORD = "DB_PASSWORD"
	ENV_DB_SSL_MODE = "DB_SSL_MODE"

	ENV_RABBITMQ_HOST     = "RABBITMQ_HOST"
	ENV_RABBITMQ_USERNAME = "RABBITMQ_USERNAME"
	ENV_RABBITMQ_PASSWORD = "RABBITMQ_PASSWORD"

	ENV_REDIS_HOST     = "REDIS_HOST"
	ENV_REDIS_PASSWORD = "REDIS_PASSWORD"

	ENV_TRACER_URL  = "TRACER_URL"
	ENV_TRACER_NAME = "TRACER_NAME"
)

var envVariables []string = []string{
	ENV_PORT,
	ENV_DB_NAME,
	ENV_DB_HOST,
	ENV_DB_USER,
	ENV_DB_PASSWORD,
	ENV_DB_SSL_MODE,
	ENV_RABBITMQ_HOST,
	ENV_RABBITMQ_USERNAME,
	ENV_RABBITMQ_PASSWORD,
	ENV_REDIS_HOST,
	ENV_TRACER_URL,
	ENV_TRACER_NAME,
}

type TracerConfig struct {
	URL  string `yaml:"url"`
	Name string `yaml:"name"`
}

type Config struct {
	DB         *database.Credentials
	RabbitMQ   *rabbitmq.Credentials
	Redis      *redis.Credentials
	Tracer     *TracerConfig
	Port       string
	Production bool
}

func Read() (*Config, error) {
	var envCfg Config
	viper.AutomaticEnv()
	isProduction := viper.GetBool(ENV_PRODUCTION)

	result := make(map[string]string, len(envVariables))
	for _, name := range envVariables {
		value := os.Getenv(name)
		if value == "" {
			return nil, fmt.Errorf("env %q is empty", name)
		}
		result[name] = value
	}

	dbCreds := &database.Credentials{
		Username: result[ENV_DB_USER],
		Password: result[ENV_DB_PASSWORD],
		DBName:   result[ENV_DB_NAME],
		Host:     result[ENV_DB_HOST],
	}

	if result[ENV_DB_SSL_MODE] != "" {
		sslMode := result[ENV_DB_SSL_MODE]
		dbCreds.SSLMode = &sslMode
	}

	rabbitMQCreds := &rabbitmq.Credentials{
		Host:     result[ENV_RABBITMQ_HOST],
		Username: result[ENV_RABBITMQ_USERNAME],
		Password: result[ENV_RABBITMQ_PASSWORD],
	}

	redisCreds := &redis.Credentials{
		Host:     result[ENV_REDIS_HOST],
		Password: result[ENV_REDIS_PASSWORD],
	}

	envCfg = Config{
		DB:       dbCreds,
		RabbitMQ: rabbitMQCreds,
		Redis:    redisCreds,
		Tracer: &TracerConfig{
			URL:  result[ENV_TRACER_URL],
			Name: result[ENV_TRACER_NAME],
		},
		Port:       result[ENV_PORT],
		Production: isProduction,
	}

	return &envCfg, nil
}
