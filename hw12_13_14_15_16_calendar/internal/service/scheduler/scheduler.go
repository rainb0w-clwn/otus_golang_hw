package scheduler

import (
	"context"
	"encoding/json"
	"time"

	"github.com/rainb0w-clwn/otus_golang_hw/hw12_13_14_15_calendar/internal/app"
	"github.com/rainb0w-clwn/otus_golang_hw/hw12_13_14_15_calendar/internal/config"
	"github.com/rainb0w-clwn/otus_golang_hw/hw12_13_14_15_calendar/internal/entity"
	"github.com/rainb0w-clwn/otus_golang_hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/rainb0w-clwn/otus_golang_hw/hw12_13_14_15_calendar/internal/queue"
)

type Scheduler struct {
	app      *app.App
	logger   logger.Logger
	qManager *queue.RabbitManager
}

func New(
	app *app.App,
	logg logger.Logger,
	qManager *queue.RabbitManager,
) *Scheduler {
	return &Scheduler{
		app:      app,
		qManager: qManager,
		logger:   logg,
	}
}

func (s *Scheduler) Run(ctx context.Context) error {
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

	channel, err := qSchedulerAck.Consume()
	if err != nil {
		s.logger.Error("Error registering consumer: %v", err)
		return err
	}

	go func() {
		for {
			select {
			case msg, ok := <-channel:
				if !ok {
					return
				}
				eventMsg := entity.EventMsg{}
				err := json.Unmarshal(msg.Body, &eventMsg)
				if err != nil {
					s.logger.Error("Error reading msg from channel: &v", err)
					continue
				}
				err = s.app.MarkEventAsReminded(eventMsg.ID)
				if err != nil {
					s.logger.Error("Error ack sent message: %v", err)
					continue
				}
				s.logger.Info("Event \"%s\" marked as sent", eventMsg.ID)
			case <-ctx.Done():
				return
			}
		}
	}()
	for {
		s.logger.Info("Looking for events to remind...")
		select {
		case <-time.After(cfg.Scheduler.Period):
			events := s.getEvents()
			if events != nil {
				s.sendEvents(events, qScheduler)
			}
		case <-ctx.Done():
			s.logger.Info("Scheduler stopped.")
			return nil
		}
	}
}

func (s *Scheduler) getEvents() *entity.Events {
	events, err := s.app.GetEventsForRemind()
	if err != nil {
		s.logger.Error("Error getting events for reminder: %v", err)
		return nil
	}

	if len(*events) == 0 {
		return nil
	}

	return events
}

func (s *Scheduler) sendEvents(events *entity.Events, queueSend *queue.RabbitQueueConnection) {
	for _, event := range *events {
		eventMsg := event.ToMsg()
		jsonMsg, err := json.Marshal(eventMsg)
		if err != nil {
			s.logger.Error("Error sending msg to RabbitMQ: %v", err)
			continue
		}

		err = queueSend.Produce(jsonMsg)
		if err != nil {
			s.logger.Error("Error sending msg to RabbitMQ: %v", err)
			continue
		}

		s.logger.Info("Event \"%s\" sent", event.ID)
	}
}
