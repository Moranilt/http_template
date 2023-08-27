package redis_mock

import (
	"github.com/Moranilt/http_template/clients/redis"
	"github.com/go-redis/redismock/v9"
)

func New() (*redis.Client, redismock.ClientMock) {
	client, mockClient := redismock.NewClientMock()

	return &redis.Client{client}, mockClient
}
