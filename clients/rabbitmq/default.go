package rabbitmq

import (
	"context"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func Push(ctx context.Context, data []byte) error {
	return rabbitMQClient.Push(ctx, data)
}

func UnsafePush(ctx context.Context, data []byte) error {
	return rabbitMQClient.UnsafePush(ctx, data)
}

func ReadMsgs(ctx context.Context, maxAmount int, wait time.Duration, callback ReadMsgCallback) {
	go rabbitMQClient.ReadMsgs(ctx, maxAmount, wait, callback)
}

func Consume() (<-chan amqp.Delivery, error) {
	return rabbitMQClient.Consume()
}

func Close() error {
	return rabbitMQClient.Close()
}
