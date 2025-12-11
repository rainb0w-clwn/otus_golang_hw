package sender

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/rainb0w-clwn/otus_golang_hw/hw12_13_14_15_calendar/internal/config"
	"github.com/rainb0w-clwn/otus_golang_hw/hw12_13_14_15_calendar/internal/entity"
	"github.com/rainb0w-clwn/otus_golang_hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/rainb0w-clwn/otus_golang_hw/hw12_13_14_15_calendar/internal/queue"
)

type Sender struct {
	logger   logger.Logger
	qManager *queue.RabbitManager
}

func New(logger logger.Logger, qManager *queue.RabbitManager) *Sender {
	return &Sender{
		logger:   logger,
		qManager: qManager,
	}
}

func (s *Sender) Run(ctx context.Context) error {
	cfg := config.GetFromContext(ctx)
	if cfg == nil {
		return config.ErrNoConfigInContext
	}

	qScheduler, err := s.qManager.CreateQueue(cfg.Scheduler.Queue)
	if err != nil {
		s.logger.Error("Error declaring scheduler queue: %v", err)
		return err
	}

	qSchedulerAck, err := s.qManager.CreateQueue(cfg.Scheduler.Queue + "_ACK")
	if err != nil {
		s.logger.Error("Error declaring scheduler_ack queue: %v", err)
		return err
	}

	channel, err := qScheduler.Consume()
	if err != nil {
		s.logger.Error("Error registering consumer: %v", err)
		return err
	}

	for {
		s.logger.Info("Waiting for events...")

		select {
		case <-ctx.Done():
			s.logger.Info("Sender stopped.")
			return nil
		case msg := <-channel:
			eventMsg := entity.EventMsg{}

			err := json.Unmarshal(msg.Body, &eventMsg)
			if err != nil {
				s.logger.Error("Error reading msg from channel: " + err.Error())

				continue
			}

			s.logger.Info(fmt.Sprintf(
				"Sending reminder about \"%s\" event to #%d user. Event time: %s.",
				eventMsg.Title,
				eventMsg.UserID,
				eventMsg.DateTime.Format(time.RFC822),
			))

			err = qSchedulerAck.Produce(msg.Body)
			if err != nil {
				s.logger.Error("Error sending msg to RabbitMQ: " + err.Error())
				continue
			}
		}
	}
}
