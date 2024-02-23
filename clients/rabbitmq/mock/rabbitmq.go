package rabbitmq_mock

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/Moranilt/http-utils/mock"
	"github.com/Moranilt/http_template/clients/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	ERR_Not_Valid_Data = "expected %v data to be %s, got %s"
)

type mockRabbit struct {
	history      *mock.MockHistory[[]byte]
	msgsCh       chan amqp.Delivery
	mockedMsgsCh chan RabbitDeliveryMocker
	msgStorage   []RabbitDeliveryMocker
	t            *testing.T
}
type RabbitMQMocker interface {
	// RabbitMQClient is the interface for the real RabbitMQ client
	rabbitmq.RabbitMQClient

	// ExpectReadMsg records an expectation that ReadMsg will be called
	ExpectReadMsg()
	// ExpectPush records an expectation that Push will be called with the given data and error
	ExpectPush(data []byte, err error)
	// ExpectUnsafePush records an expectation that UnsafePush will be called with the given data and error
	ExpectUnsafePush(data []byte, err error)
	// ExpectConsume records an expectation that Consume will be called with the given error
	ExpectConsume(err error)
	// ExpectClose records an expectation that Close will be called with the given error
	ExpectClose(err error)
	// ExpectCheck records an expectation that Check will be called with the given error
	ExpectCheck(err error)
	// Clear clears all recorded expectations
	Clear()
	// AllExpectationsDone checks if all recorded expectations have been satisfied
	AllExpectationsDone() error
}

func NewRabbitMQ(t *testing.T) RabbitMQMocker {
	// mockedMsgsCh is a channel for mocked RabbitMQ deliveries
	mockedMsgsCh := make(chan RabbitDeliveryMocker, 16)

	// msgStorage stores mocked RabbitMQ deliveries
	msgStorage := make([]RabbitDeliveryMocker, 0)

	// msgsCh is a channel for RabbitMQ deliveries
	msgsCh := make(chan amqp.Delivery, 1)

	return &mockRabbit{
		history:      mock.NewMockHistory[[]byte](),
		mockedMsgsCh: mockedMsgsCh,
		msgStorage:   msgStorage,
		msgsCh:       msgsCh,
		t:            t,
	}
}

func (m *mockRabbit) ExpectReadMsg() {
	m.history.Push("ReadMsg", nil, nil)
}

func (m *mockRabbit) ExpectPush(data []byte, err error) {
	m.history.Push("Push", data, err)
}

func (m *mockRabbit) ExpectUnsafePush(data []byte, err error) {
	m.history.Push("UnsafePush", data, err)
}

func (m *mockRabbit) ExpectConsume(err error) {
	m.history.Push("Consume", nil, err)
}

func (m *mockRabbit) ExpectClose(err error) {
	m.history.Push("Close", nil, err)
}

func (m *mockRabbit) ExpectCheck(err error) {
	m.history.Push("Check", nil, err)
}

func (m *mockRabbit) Clear() {
	m.history.Clear()
	m.msgsCh = make(chan amqp.Delivery, 1)
	m.mockedMsgsCh = make(chan RabbitDeliveryMocker, 16)
	m.msgStorage = make([]RabbitDeliveryMocker, 0)
}

// AllExpectationsDone checks if all expected calls were done.
// Returns error if some expectations were not met.
func (m *mockRabbit) AllExpectationsDone() error {
	return m.history.AllExpectationsDone()
}

func (m *mockRabbit) Check(ctx context.Context) error {
	item, err := m.history.Get("Check")
	if err != nil {
		return err
	}
	return item.Err
}

func (m *mockRabbit) ReadMsgs(ctx context.Context, maxAmount int, wait time.Duration, callback rabbitmq.ReadMsgCallback) {
	m.t.Helper()
	_, err := m.history.Get("ReadMsg")
	if err != nil {
		m.t.Error(err)
	}
}

func (m *mockRabbit) testPush(name string, data []byte) error {
	item, err := m.history.Get(name)
	if err != nil {
		return err
	}
	if !reflect.DeepEqual(item.Data, data) {
		return fmt.Errorf(ERR_Not_Valid_Data, name, item.Data, data)
	}
	m.mockedMsgsCh <- NewDelivery(m.t, MockRabbitDeliveryFields{
		Body: data,
	})
	return item.Err
}

func (m *mockRabbit) Push(ctx context.Context, data []byte) error {
	m.t.Helper()
	return m.testPush("Push", data)
}

func (m *mockRabbit) UnsafePush(ctx context.Context, data []byte) error {
	m.t.Helper()
	return m.testPush("UnsafePush", data)
}

func (m *mockRabbit) Consume() (<-chan amqp.Delivery, error) {
	item, err := m.history.Get("Consume")
	if err != nil {
		return nil, err
	}

	go func(mockedMsgs chan RabbitDeliveryMocker, target chan<- amqp.Delivery) {
		var counter uint64
		for val := range mockedMsgs {
			target <- amqp.Delivery{
				DeliveryTag: counter,
				Body:        val.Body(),
			}
			counter++
		}

	}(m.mockedMsgsCh, m.msgsCh)
	return m.msgsCh, item.Err
}

func (m *mockRabbit) Close() error {
	item, err := m.history.Get("Close")
	if err != nil {
		return err
	}

	return item.Err
}

type RabbitDeliveryMocker interface {
	rabbitmq.RabbitDelivery

	ExpectNack(multiple bool, requeue bool, err error)
	ExpectAck(multiple bool, err error)
	ExpectReject(requeue bool, err error)
	AllExpectationsDone() error
}

type MockRabbitDeliveryFields struct {
	Priority        uint8
	Headers         amqp.Table
	DeliveryMode    uint8
	CorelationId    string
	ContentType     string
	ContentEncoding string
	Body            []byte
	Acknowledger    amqp.Acknowledger
	MessageId       string
	AppId           string
	Timestamp       time.Time
	MsgType         string
	UserId          string
	ConsumerTag     string
	DeliveryTag     uint64
	Redelivered     bool
	Exchange        string
	RoutingKey      string
	ReplyTo         string
	Expiration      string
	MessageCount    uint32
}
type mockRabbitDelivery struct {
	history *mock.MockHistory[mockDeliveryData]

	t *testing.T

	d MockRabbitDeliveryFields
}

func NewDelivery(t *testing.T, d MockRabbitDeliveryFields) RabbitDeliveryMocker {
	return &mockRabbitDelivery{
		history: mock.NewMockHistory[mockDeliveryData](),
		t:       t,
		d:       d,
	}
}

type mockDeliveryData struct {
	multiple bool
	requeue  bool
}

func (r *mockRabbitDelivery) Priority() uint8 {
	return r.d.Priority
}

func (r *mockRabbitDelivery) Header() amqp.Table {
	return r.d.Headers
}

func (r *mockRabbitDelivery) DeliveryMode() uint8 {
	return r.d.DeliveryMode
}

func (r *mockRabbitDelivery) CorelationId() string {
	return r.d.CorelationId
}

func (r *mockRabbitDelivery) ContentType() string {
	return r.d.ContentType
}

func (r *mockRabbitDelivery) ContentEncoding() string {
	return r.d.ContentEncoding
}

func (r *mockRabbitDelivery) Body() []byte {
	return r.d.Body
}

func (r *mockRabbitDelivery) Acknowledger() amqp.Acknowledger {
	return r.d.Acknowledger
}

func (r *mockRabbitDelivery) MessageId() string {
	return r.d.MessageId
}

func (r *mockRabbitDelivery) AppId() string {
	return r.d.AppId
}

func (r *mockRabbitDelivery) Timestamp() time.Time {
	return r.d.Timestamp
}

func (r *mockRabbitDelivery) Type() string {
	return r.d.MsgType
}

func (r *mockRabbitDelivery) UserId() string {
	return r.d.UserId
}

func (r *mockRabbitDelivery) ConsumerTag() string {
	return r.d.ConsumerTag
}

func (r *mockRabbitDelivery) DeliveryTag() uint64 {
	return r.d.DeliveryTag
}

func (r *mockRabbitDelivery) Redelivered() bool {
	return r.d.Redelivered
}

func (r *mockRabbitDelivery) Exchange() string {
	return r.d.Exchange
}

func (r *mockRabbitDelivery) RoutingKey() string {
	return r.d.RoutingKey
}

func (r *mockRabbitDelivery) ReplyTo() string {
	return r.d.ReplyTo
}

func (r *mockRabbitDelivery) Expiration() string {
	return r.d.Expiration
}

func (r *mockRabbitDelivery) MessageCount() uint32 {
	return r.d.MessageCount
}

func (m *mockRabbitDelivery) ExpectAck(multiple bool, err error) {
	m.history.Push("Ack", mockDeliveryData{
		multiple: multiple,
	}, err)
}

func (m *mockRabbitDelivery) ExpectNack(multiple bool, requeue bool, err error) {
	m.history.Push("Nack", mockDeliveryData{
		multiple: multiple,
		requeue:  requeue,
	}, err)
}

func (m *mockRabbitDelivery) ExpectReject(requeue bool, err error) {
	m.history.Push("Reject", mockDeliveryData{
		requeue: requeue,
	}, err)
}

func (m *mockRabbitDelivery) Ack(multiple bool) error {
	m.t.Helper()
	item, err := m.history.Get("Ack")
	if err != nil {
		m.t.Error(err)
		return err
	}

	if item.Data.multiple != multiple {
		err := fmt.Errorf("expected multiple ack to be %t, got %t", item.Data.multiple, multiple)
		m.t.Error(err)
		return err
	}

	return item.Err
}

func (m *mockRabbitDelivery) Nack(multiple bool, requeue bool) error {
	m.t.Helper()
	item, err := m.history.Get("Nack")
	if err != nil {
		m.t.Error(err)
		return err
	}

	if item.Data.multiple != multiple {
		return fmt.Errorf("expected multiple Nack to be %t, got %t", item.Data.multiple, multiple)
	}

	if item.Data.requeue != requeue {
		return fmt.Errorf("expected requeue Nack to be %t, got %t", item.Data.requeue, requeue)
	}

	return item.Err
}

func (m *mockRabbitDelivery) Reject(requeue bool) error {
	m.t.Helper()
	item, err := m.history.Get("Reject")
	if err != nil {
		return err
	}

	if item.Data.requeue != requeue {
		return fmt.Errorf("expected requeue Reject to be %t, got %t", item.Data.requeue, requeue)
	}

	return item.Err
}

func (m *mockRabbitDelivery) AllExpectationsDone() error {
	return m.history.AllExpectationsDone()
}
