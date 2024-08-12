package app

import (
	"context"
	"fmt"
	"github.com/EkaterinaShamanaeva/otus-go/hw12_13_14_15_calendar/internal/storage"
	"github.com/EkaterinaShamanaeva/otus-go/hw12_13_14_15_calendar/internal/storage/init_storage"
	"github.com/gofrs/uuid"
	"time"
)

type App struct {
	logger  Logger
	storage storage.Storage
}

type Logger interface {
	Info(msg string)
	Warn(msg string)
	Error(msg string)
	Debug(msg string)
}

type Event struct {
	ID               uuid.UUID     `json:"id"`
	Title            string        `json:"title"`
	TimeStart        time.Time     `json:"start_date"`
	Duration         time.Duration `json:"duration"`
	Description      string        `json:"description"`
	UserID           uuid.UUID     `json:"user_id"`
	NotifyBeforeDays int           `json:"notify_before"`
}

type Application interface {
	Connect(ctx context.Context, connect string) error
	CreateEvent(ctx context.Context, event *Event) error
	GetEventID(ctx context.Context, event *Event) (uuid.UUID, error)
	UpdateEvent(ctx context.Context, event *Event) error
	DeleteEvent(ctx context.Context, id uuid.UUID) error
	GetEventsPerDay(ctx context.Context, day time.Time) ([]Event, error)
	GetEventsPerWeek(ctx context.Context, beginDate time.Time) ([]Event, error)
	GetEventsPerMonth(ctx context.Context, beginDate time.Time) ([]Event, error)
	Close(ctx context.Context) error
}

func New(logger Logger, storage storage.Storage) *App {
	return &App{
		logger:  logger,
		storage: storage,
	}
}

func (a *App) Connect(ctx context.Context, storageType string, dsn string) error {
	storage, err := init_storage.NewStorage(ctx, storageType, dsn)
	if err != nil {
		return err
	}
	a.storage = storage
	return nil
}

func (a *App) Close(ctx context.Context) error {
	a.logger.Info("storage closing...")
	return a.storage.Close(ctx)
}

func (a *App) CreateEvent(ctx context.Context, event *Event) error {
	if event.Title == "" {
		return fmt.Errorf("empty title of event")
	}
	if event.Duration <= 0 {
		return fmt.Errorf("wrong duration")
	}
	id, err := uuid.NewV4()
	if err != nil {
		return err
	}
	return a.storage.CreateEvent(ctx, &storage.Event{
		ID:               id,
		Title:            event.Title,
		TimeStart:        event.TimeStart,
		Duration:         event.Duration,
		Description:      event.Description,
		UserID:           event.UserID,
		NotifyBeforeDays: event.NotifyBeforeDays,
	})
}

func (a *App) UpdateEvent(ctx context.Context, event *Event) error {
	if event.Title == "" {
		return fmt.Errorf("empty title of event")
	}
	if event.Duration <= 0 {
		return fmt.Errorf("wrong duration")
	}

	return a.storage.UpdateEvent(ctx, &storage.Event{
		ID:               event.ID,
		Title:            event.Title,
		TimeStart:        event.TimeStart,
		Duration:         event.Duration,
		Description:      event.Description,
		UserID:           event.UserID,
		NotifyBeforeDays: event.NotifyBeforeDays,
	})
}

func (a *App) GetEventID(ctx context.Context, event *storage.Event) (uuid.UUID, error) {
	return a.storage.GetEventID(ctx, event)
}

func (a *App) DeleteEvent(ctx context.Context, id uuid.UUID) error {
	return a.storage.DeleteEvent(ctx, id)
}

func (a *App) GetEventsPerDay(ctx context.Context, day time.Time) ([]storage.Event, error) {
	return a.storage.GetEventsPerDay(ctx, day)
}

func (a *App) GetEventsPerWeek(ctx context.Context, beginDate time.Time) ([]storage.Event, error) {
	return a.storage.GetEventsPerWeek(ctx, beginDate)
}

func (a *App) GetEventsPerMonth(ctx context.Context, beginDate time.Time) ([]storage.Event, error) {
	return a.storage.GetEventsPerMonth(ctx, beginDate)
}
