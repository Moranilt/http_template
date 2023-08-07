package redis

import (
	"context"
	"time"

	"github.com/Moranilt/http_template/clients/credentials"
	"github.com/redis/go-redis/v9"
)

func New(ctx context.Context, creds *credentials.Redis) (*redis.Client, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:         creds.Host,
		Password:     creds.Password,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	})

	if ping := redisClient.Ping(ctx); ping.Err() != nil {
		return nil, ping.Err()
	}

	return redisClient, nil
}
