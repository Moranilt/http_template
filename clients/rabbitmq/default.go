package rabbitmq

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

func Push(ctx context.Context, data []byte) error {
	return rabbitMQClient.Push(ctx, data)
}

func UnsafePush(ctx context.Context, data []byte) error {
	return rabbitMQClient.UnsafePush(ctx, data)
}

func ReadMsgs(ctx context.Context, maxAmount int, callback func(d amqp.Delivery) error) {
	go rabbitMQClient.ReadMsgs(ctx, maxAmount, callback)
}

func Consume() (<-chan amqp.Delivery, error) {
	return rabbitMQClient.Consume()
}

func Close() error {
	return rabbitMQClient.Close()
}
