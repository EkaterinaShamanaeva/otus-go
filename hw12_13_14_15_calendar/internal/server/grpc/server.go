package grpc

import (
	"context"
	"fmt"
	"github.com/EkaterinaShamanaeva/otus-go/hw12_13_14_15_calendar/internal/app"
	"github.com/EkaterinaShamanaeva/otus-go/hw12_13_14_15_calendar/internal/storage"
	"github.com/gofrs/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/EkaterinaShamanaeva/otus-go/hw12_13_14_15_calendar/internal/logger"
	"google.golang.org/grpc"
	"net"
	"time"
)

//go:generate protoc -I ../../../api EventService.proto --go_out=. --go-grpc_out=.
type server struct {
	srv  *grpc.Server
	app  *app.App
	logg *logger.Logger
	UnimplementedCalendarServer
}

func New(logg *logger.Logger, app *app.App) *server {
	return &server{
		app:  app,
		logg: logg,
	}
}

func (s *server) Start(ctx context.Context, addr string) error {
	s.logg.Info("gRPC server starting...")
	lsn, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	s.srv = grpc.NewServer(grpc.UnaryInterceptor(loggingServerInterceptor(*s.logg)))
	RegisterCalendarServer(s.srv, s)
	if err = s.srv.Serve(lsn); err != nil {
		return err
	}
	return nil
}

func (s *server) Stop(ctx context.Context) error {
	s.logg.Info("gRPC server stopping...")
	s.srv.GracefulStop()
	return nil
}

func loggingServerInterceptor(logger app.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ interface{}, err error) {
		logger.Info(fmt.Sprintf("method: %s, duration: %s, request: %+v", info.FullMethod, time.Since(time.Now()), req))
		h, err := handler(ctx, req)
		return h, err
	}
}

func (s *server) CreateEvent(ctx context.Context, request *CreateEventRequest) (*CreateEventResponse, error) {
	id, err := uuid.FromString(request.Event.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("error to create event: %v", err))
	}
	userId, err := uuid.FromString(request.Event.UserId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("error to create event: %v", err))
	}

	ev := app.Event{
		ID:               id,
		Title:            request.Event.Title,
		TimeStart:        request.Event.TimeStart.AsTime(),
		Duration:         request.Event.Duration.AsDuration(),
		Description:      request.Event.Description,
		UserID:           userId,
		NotifyBeforeDays: int(request.Event.NotifyBefore),
	}

	err = s.app.CreateEvent(ctx, &ev)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("error to create event: %v", err))
	}

	return &CreateEventResponse{Id: ev.ID.String()}, nil
}

func (s *server) UpdateEvent(ctx context.Context, request *UpdateEventRequest) (*emptypb.Empty, error) {
	id, err := uuid.FromString(request.Event.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("error to create event: %v", err))
	}
	userId, err := uuid.FromString(request.Event.UserId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("error to create event: %v", err))
	}

	ev := app.Event{
		ID:               id,
		Title:            request.Event.Title,
		TimeStart:        request.Event.TimeStart.AsTime(),
		Duration:         request.Event.Duration.AsDuration(),
		Description:      request.Event.Description,
		UserID:           userId,
		NotifyBeforeDays: int(request.Event.NotifyBefore),
	}

	err = s.app.UpdateEvent(ctx, &ev)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("error to create event: %v", err))
	}

	return &emptypb.Empty{}, nil
}

func (s *server) DeleteEvent(ctx context.Context, request *DeleteEventRequest) (*emptypb.Empty, error) {
	id, err := uuid.FromString(request.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("error to delete event: %v", err))
	}

	err = s.app.DeleteEvent(ctx, id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("error to delete event: %v", err))
	}

	return &emptypb.Empty{}, nil
}

func (s *server) GetEventsPerDay(ctx context.Context, request *GetEventsPerDayRequest) (*GetEventsPerDayResponse, error) {
	listEvents, err := s.app.GetEventsPerDay(ctx, request.TimeStart.AsTime())
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("error to get evens list: %v", err))
	}

	return &GetEventsPerDayResponse{Events: convertStorageEvToGrpcEv(listEvents)}, nil
}

func (s *server) GetEventsPerWeek(ctx context.Context, request *GetEventsPerWeekRequest) (*GetEventsPerWeekResponse, error) {
	listEvents, err := s.app.GetEventsPerWeek(ctx, request.Day.AsTime())
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("error to get evens list: %v", err))
	}

	return &GetEventsPerWeekResponse{Events: convertStorageEvToGrpcEv(listEvents)}, nil
}

func (s *server) GetEventsPerMonth(ctx context.Context, request *GetEventsPerMonthRequest) (*GetEventsPerMonthResponse, error) {
	listEvents, err := s.app.GetEventsPerMonth(ctx, request.BeginDate.AsTime())
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("error to get evens list: %v", err))
	}

	return &GetEventsPerMonthResponse{Events: convertStorageEvToGrpcEv(listEvents)}, nil
}

func convertStorageEvToGrpcEv(events []storage.Event) []*Event {
	resultEvents := make([]*Event, 0, len(events))
	for _, event := range events {
		resultEvent := &Event{
			Id:           event.ID.String(),
			Title:        event.Title,
			TimeStart:    timestamppb.New(event.TimeStart),
			Duration:     durationpb.New(event.Duration),
			Description:  event.Description,
			UserId:       event.UserID.String(),
			NotifyBefore: int32(event.NotifyBeforeDays),
		}
		resultEvents = append(resultEvents, resultEvent)
	}
	return resultEvents
}
