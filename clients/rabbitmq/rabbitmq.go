package rabbitmq

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/Moranilt/http_template/clients/credentials"
	"github.com/Moranilt/http_template/logger"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	reconnectDelay = 5 * time.Second
	reInitDelay    = 2 * time.Second
	resendDelay    = 5 * time.Second
)

var (
	errNotConnected  = errors.New("not connected to a server")
	errAlreadyClosed = errors.New("already closed: not connected to the server")
	errShutdown      = errors.New("client is shutting down")
)

var (
	rabbitMQClient *Client
)

type Client struct {
	logger          *logger.SLogger
	queueName       string
	connection      *amqp.Connection
	channel         *amqp.Channel
	done            chan bool
	notifyConnClose chan *amqp.Error
	notifyChanClose chan *amqp.Error
	notifyConfirm   chan amqp.Confirmation
	isReady         bool
	readyCh         chan bool
	mu              sync.Mutex
}

func (c *Client) setIsReady(ready bool) {
	c.mu.Lock()
	c.isReady = ready
	c.mu.Unlock()
}

func (c *Client) getIsReady() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.isReady
}

func (c *Client) Check(ctx context.Context) error {
	if c.connection.IsClosed() {
		return errors.New("connection closed")
	}
	if c.channel.IsClosed() {
		return errors.New("channel closed")
	}
	return nil
}

func Init(ctx context.Context, queueName string, log *logger.SLogger, creds credentials.SourceStringer) *Client {
	client := Client{
		logger:    log,
		queueName: queueName,
		done:      make(chan bool),
		readyCh:   make(chan bool, 3),
	}
	go client.handleReconnect(ctx, creds.SourceString())

	rabbitMQClient = &client
	return &client
}

func (client *Client) handleReconnect(ctx context.Context, addr string) {
	for {
		client.isReady = false
		client.readyCh <- false
		client.logger.InfoContext(ctx, "Attempting to connect")

		conn, err := client.connect(addr)

		if err != nil {
			client.logger.InfoContext(ctx, "Failed to connect. Retrying...")

			select {
			case <-client.done:
				return
			case <-time.After(reconnectDelay):
			}
			continue
		}

		if done := client.handleReInit(ctx, conn); done {
			break
		}
	}
}

func (client *Client) connect(addr string) (*amqp.Connection, error) {
	conn, err := amqp.Dial(addr)

	if err != nil {
		return nil, err
	}

	client.changeConnection(conn)
	client.logger.Info("Connected!")
	return conn, nil
}

func (client *Client) handleReInit(ctx context.Context, conn *amqp.Connection) bool {
	for {
		client.setIsReady(false)
		client.readyCh <- false

		err := client.init(conn)

		if err != nil {
			client.logger.InfoContext(ctx, "Failed to initialize channel. Retrying...")

			select {
			case <-ctx.Done():
				return true
			case <-client.done:
				return true
			case <-client.notifyConnClose:
				client.logger.InfoContext(ctx, "Connection closed. Reconnecting...")
				return false
			case <-time.After(reInitDelay):
			}
			continue
		}

		select {
		case <-ctx.Done():
			return true
		case <-client.done:
			return true
		case <-client.notifyConnClose:
			client.logger.InfoContext(ctx, "Connection closed. Reconnecting...")
			return false
		case <-client.notifyChanClose:
			client.logger.InfoContext(ctx, "Channel closed. Re-running init...")
		}
	}
}

func (client *Client) init(conn *amqp.Connection) error {
	ch, err := conn.Channel()

	if err != nil {
		return err
	}

	err = ch.Confirm(false)

	if err != nil {
		return err
	}

	_, err = ch.QueueDeclare(
		client.queueName,
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return err
	}

	client.changeChannel(ch)
	client.setIsReady(true)
	client.readyCh <- true
	client.logger.Info("Setup!")

	return nil
}

func (client *Client) ReadMsgs(ctx context.Context, maxAmount int, callback func(d amqp.Delivery) error) {
	stopReceive := make(chan bool, 1)

	for {
		select {
		case <-ctx.Done():
			close(stopReceive)
			return
		case ready := <-client.readyCh:
			if !ready {
				if len(stopReceive) == 0 {
					stopReceive <- true
				} else {
					<-stopReceive
				}
				client.logger.Info("Connection closed. Unable to receive messages. Waiting to reconnect...")
				<-time.After(2 * time.Second)
				continue
			} else {
				if len(stopReceive) != 0 {
					<-stopReceive
				}

				go func(done <-chan bool) {
					client.logger.Info("Run consumer...")
					deliveries, err := client.Consume()
					if err != nil {
						client.logger.Infof("Could not start consuming: %s\n", err)
						return
					}

					counter := 1
					for {
						select {
						case <-done:
							client.logger.Info("Consumer stopped...")
							return
						case d := <-deliveries:
							if counter == maxAmount {
								counter = 1
								<-time.After(15 * time.Second)
							} else {
								counter++
							}
							err := callback(d)
							if err != nil {
								continue
							}
						}
					}
				}(stopReceive)
			}
		}

	}
}

func (client *Client) changeConnection(connection *amqp.Connection) {
	client.connection = connection
	client.notifyConnClose = make(chan *amqp.Error, 1)
	client.connection.NotifyClose(client.notifyConnClose)
}

func (client *Client) changeChannel(channel *amqp.Channel) {
	client.channel = channel
	client.notifyChanClose = make(chan *amqp.Error, 1)
	client.notifyConfirm = make(chan amqp.Confirmation, 1)
	client.channel.NotifyClose(client.notifyChanClose)
	client.channel.NotifyPublish(client.notifyConfirm)
}

func (client *Client) Push(ctx context.Context, data []byte) error {
	if !client.isReady {
		return errors.New("failed to push: not connected")
	}
	for {
		err := client.UnsafePush(ctx, data)
		if err != nil {
			client.logger.InfoContext(ctx, "Push failed. Retrying...")
			select {
			case <-client.done:
				return errShutdown
			case <-time.After(resendDelay):
			}
			continue
		}
		confirm := <-client.notifyConfirm
		if confirm.Ack {
			client.logger.Infof("Push confirmed [%d]!", confirm.DeliveryTag)
			return nil
		}
	}
}

func (client *Client) UnsafePush(ctx context.Context, data []byte) error {
	if !client.isReady {
		return errNotConnected
	}

	msg := amqp.Publishing{
		ContentType: "application/json",
		Body:        data,
	}

	return client.channel.PublishWithContext(
		ctx,
		"",
		client.queueName,
		false,
		false,
		msg,
	)
}

func (client *Client) Consume() (<-chan amqp.Delivery, error) {
loopReady:
	for !client.isReady {
		<-time.After(resendDelay)
		continue loopReady
	}

	if err := client.channel.Qos(
		2,
		0,
		false,
	); err != nil {
		return nil, err
	}

	return client.channel.Consume(
		client.queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
}

func (client *Client) Close() error {
	if !client.getIsReady() {
		return errAlreadyClosed
	}
	close(client.done)
	err := client.channel.Close()
	if err != nil {
		return err
	}
	err = client.connection.Close()
	if err != nil {
		return err
	}

	client.isReady = false
	for len(client.readyCh) > 0 {
		<-client.readyCh
	}
	close(client.readyCh)
	return nil
}
