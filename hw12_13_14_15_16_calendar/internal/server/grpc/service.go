package server

import (
	"context"

	proto "github.com/rainb0w-clwn/otus_golang_hw/hw12_13_14_15_calendar/api"
	"github.com/rainb0w-clwn/otus_golang_hw/hw12_13_14_15_calendar/internal/entity"
	"github.com/rainb0w-clwn/otus_golang_hw/hw12_13_14_15_calendar/internal/logger"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Service struct {
	proto.UnimplementedEventServiceServer
	app    Application
	logger logger.Logger
}

func NewService(app Application, logger logger.Logger) *Service {
	return &Service{
		app:    app,
		logger: logger,
	}
}

func (s Service) CreateEvent(_ context.Context, request *proto.CreateRequest) (*proto.CreateResponse, error) {
	id, err := s.app.CreateEvent(s.proto2entity(request.GetEventData()))
	if err != nil {
		s.logger.Error(err.Error())

		return nil, err
	}

	return &proto.CreateResponse{EventId: &(proto.EventId{Id: id})}, nil
}

func (s Service) UpdateEvent(_ context.Context, request *proto.UpdateRequest) (*proto.UpdateResponse, error) {
	err := s.app.UpdateEvent(
		request.GetEventId().GetId(),
		s.proto2entity(request.GetEventData()),
	)
	if err != nil {
		s.logger.Error(err.Error())

		return nil, err
	}

	return &proto.UpdateResponse{}, nil
}

func (s Service) DeleteEvent(_ context.Context, request *proto.DeleteRequest) (*proto.DeleteResponse, error) {
	err := s.app.DeleteEvent(request.GetEventId().GetId())
	if err != nil {
		s.logger.Error(err.Error())

		return nil, err
	}

	return &proto.DeleteResponse{}, nil
}

func (s Service) GetWeekEvents(_ context.Context, date *proto.StartDate) (*proto.Events, error) {
	events, err := s.app.GetWeekEvents(date.GetStartDate().AsTime())
	if err != nil {
		s.logger.Error(err.Error())

		return nil, err
	}

	return s.entities2Proto(events), nil
}

func (s Service) GetMonthEvents(_ context.Context, date *proto.StartDate) (*proto.Events, error) {
	events, err := s.app.GetMonthEvents(date.GetStartDate().AsTime())
	if err != nil {
		s.logger.Error(err.Error())

		return nil, err
	}

	return s.entities2Proto(events), nil
}

func (s Service) GetDayEvents(_ context.Context, date *proto.StartDate) (*proto.Events, error) {
	events, err := s.app.GetDayEvents(date.GetStartDate().AsTime())
	if err != nil {
		s.logger.Error(err.Error())

		return nil, err
	}

	return s.entities2Proto(events), nil
}

func (s Service) entities2Proto(entityEvents *entity.Events) *proto.Events {
	protoEvents := make([]*proto.Event, 0, len(*entityEvents))

	for _, entityEvent := range *entityEvents {
		protoEvents = append(
			protoEvents,
			s.entity2Proto(entityEvent),
		)
	}

	return &proto.Events{Events: protoEvents}
}

func (s Service) entity2Proto(entityEvent *entity.Event) *proto.Event {
	return &(proto.Event{
		EventId: &proto.EventId{Id: entityEvent.ID},
		EventData: &proto.EventData{
			Title:       entityEvent.Title,
			DateTime:    timestamppb.New(entityEvent.DateTime),
			Description: entityEvent.Description,
			Duration:    entityEvent.Duration,
			RemindTime:  timestamppb.New(entityEvent.RemindTime),
			CreatedAt:   timestamppb.New(entityEvent.CreatedAt),
			UpdatedAt:   timestamppb.New(entityEvent.CreatedAt),
		},
	})
}

func (s Service) proto2entity(protoEvent *proto.EventData) entity.Event {
	return entity.Event{
		Title:       protoEvent.Title,
		Description: protoEvent.Description,
		DateTime:    protoEvent.DateTime.AsTime(),
		Duration:    protoEvent.Duration,
		RemindTime:  protoEvent.RemindTime.AsTime(),
		UserID:      int(protoEvent.UserId),
	}
}
