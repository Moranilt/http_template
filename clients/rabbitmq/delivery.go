package rabbitmq

import (
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitDelivery interface {
	// Ack acknowledges processing of a Delivery.
	Ack(multiple bool) error
	// Nack negatively acknowledges a Delivery.
	Nack(multiple bool, requeue bool) error
	// Reject rejects a delivery.
	Reject(requeue bool) error

	// Body returns the message body.
	Body() []byte

	// Acknowledger provides acknowledgement information.
	Acknowledger() amqp.Acknowledger
	// Header returns the message header.
	Header() amqp.Table
	// ContentType returns the message content type.
	ContentType() string
	// ContentEncoding returns the message content encoding.
	ContentEncoding() string
	// DeliveryMode returns the delivery mode.
	DeliveryMode() uint8
	// Priority returns the message priority.
	Priority() uint8
	// CorelationId returns the correlation id.
	CorelationId() string
	// ReplyTo returns the reply to value.
	ReplyTo() string
	// Expiration returns the message expiration.
	Expiration() string
	// MessageId returns the message id.
	MessageId() string
	// Timestamp returns the message timestamp.
	Timestamp() time.Time
	// Type returns the message type.
	Type() string
	// UserId returns the creating user id.
	UserId() string
	// AppId returns the creating application id.
	AppId() string
	// ConsumerTag returns the consumer tag.
	ConsumerTag() string
	// MessageCount returns the number of messages pending acknowledgement.
	MessageCount() uint32
	// DeliveryTag returns the delivery tag.
	DeliveryTag() uint64
	// Redelivered returns true if this message is being redelivered.
	Redelivered() bool
	// Exchange returns the exchange this message was published to.
	Exchange() string
	// RoutingKey returns the routing key used when publishing this message.
	RoutingKey() string
}

func NewDelivery(d amqp.Delivery) RabbitDelivery {
	return &rabbitDelivery{d: d}
}

type rabbitDelivery struct {
	d amqp.Delivery
}

func (r *rabbitDelivery) Priority() uint8 {
	return r.d.Priority
}

func (r *rabbitDelivery) Header() amqp.Table {
	return r.d.Headers
}

func (r *rabbitDelivery) DeliveryMode() uint8 {
	return r.d.DeliveryMode
}

func (r *rabbitDelivery) CorelationId() string {
	return r.d.CorrelationId
}

func (r *rabbitDelivery) ContentType() string {
	return r.d.ContentType
}

func (r *rabbitDelivery) ContentEncoding() string {
	return r.d.ContentEncoding
}

func (r *rabbitDelivery) Body() []byte {
	return r.d.Body
}

func (r *rabbitDelivery) Acknowledger() amqp.Acknowledger {
	return r.d.Acknowledger
}

func (r *rabbitDelivery) Ack(multiple bool) error {
	return r.d.Ack(multiple)
}

func (r *rabbitDelivery) Nack(multiple, requeue bool) error {
	return r.d.Nack(multiple, requeue)
}

func (r *rabbitDelivery) Reject(requeue bool) error {
	return r.d.Reject(requeue)
}

func (r *rabbitDelivery) MessageId() string {
	return r.d.MessageId
}

func (r *rabbitDelivery) AppId() string {
	return r.d.AppId
}

func (r *rabbitDelivery) Timestamp() time.Time {
	return r.d.Timestamp
}

func (r *rabbitDelivery) Type() string {
	return r.d.Type
}

func (r *rabbitDelivery) UserId() string {
	return r.d.UserId
}

func (r *rabbitDelivery) ConsumerTag() string {
	return r.d.ConsumerTag
}

func (r *rabbitDelivery) DeliveryTag() uint64 {
	return r.d.DeliveryTag
}

func (r *rabbitDelivery) Redelivered() bool {
	return r.d.Redelivered
}

func (r *rabbitDelivery) Exchange() string {
	return r.d.Exchange
}

func (r *rabbitDelivery) RoutingKey() string {
	return r.d.RoutingKey
}

func (r *rabbitDelivery) CorrelationId() string {
	return r.d.CorrelationId
}

func (r *rabbitDelivery) ReplyTo() string {
	return r.d.ReplyTo
}

func (r *rabbitDelivery) Expiration() string {
	return r.d.Expiration
}

func (r *rabbitDelivery) MessageCount() uint32 {
	return r.d.MessageCount
}
