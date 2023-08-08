package redis

import (
	"context"
	"time"

	"github.com/Moranilt/http_template/clients/credentials"
	"github.com/redis/go-redis/v9"
)

type Client struct {
	*redis.Client
}

func (r *Client) Check(ctx context.Context) error {
	return r.Ping(ctx).Err()
}

func New(ctx context.Context, creds *credentials.Redis) (*Client, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:         creds.Host,
		Password:     creds.Password,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	})

	if ping := redisClient.Ping(ctx); ping.Err() != nil {
		return nil, ping.Err()
	}

	return &Client{
		redisClient,
	}, nil
}
