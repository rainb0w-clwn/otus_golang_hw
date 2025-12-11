package queue

import (
	"context"
	"errors"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rainb0w-clwn/otus_golang_hw/hw12_13_14_15_calendar/internal/config"
	"github.com/rainb0w-clwn/otus_golang_hw/hw12_13_14_15_calendar/internal/logger"
)

var (
	ErrQueueNotConnected     = errors.New("need to connect first")
	ErrQueueAlreadyConnected = errors.New("already exists")
)

type RabbitQueueConnection struct {
	Name        string
	SendChannel *amqp.Channel
	Queue       *amqp.Queue
}

type RabbitManager struct {
	// sync.Mutex
	connection *amqp.Connection
	queueMap   map[string]*RabbitQueueConnection
	logger     logger.Logger
}

func NewRabbitManager(logger logger.Logger) *RabbitManager {
	return &RabbitManager{
		logger:   logger,
		queueMap: make(map[string]*RabbitQueueConnection),
	}
}

func (q *RabbitManager) Connect(ctx context.Context) error {
	cfg := config.GetFromContext(ctx)
	if cfg == nil {
		return config.ErrNoConfigInContext
	}
	var err error
	q.connection, err = amqp.Dial(fmt.Sprintf(
		"amqp://%s:%s@%s:%s/",
		cfg.RMQ.Login,
		cfg.RMQ.Password,
		cfg.RMQ.Host,
		cfg.RMQ.Port,
	))
	if err != nil {
		return err
	}
	return nil
}

func (q *RabbitManager) CreateQueue(queueName string) (*RabbitQueueConnection, error) {
	if q.connection == nil {
		return nil, ErrQueueNotConnected
	}
	if _, exists := q.queueMap[queueName]; exists {
		return nil, ErrQueueAlreadyConnected
	}
	channel, err := q.connection.Channel()
	if err != nil {
		return nil, err
	}

	ch := channel.NotifyPublish(make(chan amqp.Confirmation, 100))

	queue, err := channel.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		close(ch)
		return nil, err
	}

	rq := &RabbitQueueConnection{
		Name:        queueName,
		SendChannel: channel,
		Queue:       &queue,
	}
	q.queueMap[queueName] = rq
	return rq, nil
}

func (q *RabbitManager) Close() error {
	err := q.connection.Close()
	if err != nil {
		return err
	}
	for queueName, conn := range q.queueMap {
		err = conn.Close()
		if err != nil {
			return err
		}
		delete(q.queueMap, queueName)
	}

	return nil
}

func (q *RabbitQueueConnection) Close() error {
	err := q.SendChannel.Close()
	if err != nil {
		return err
	}
	return nil
}

func (q *RabbitQueueConnection) Produce(jsonMsg []byte) error {
	return q.SendChannel.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        jsonMsg,
		},
	)
}

func (q *RabbitQueueConnection) Consume() (<-chan amqp.Delivery, error) {
	messages, consumeErr := q.SendChannel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if consumeErr != nil {
		return nil, consumeErr
	}

	return messages, nil
}
