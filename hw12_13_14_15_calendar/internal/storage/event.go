package storage

import (
	"github.com/gofrs/uuid"
	"golang.org/x/net/context"
	"time"
)

type Event struct {
	ID               uuid.UUID
	Title            string
	TimeStart        time.Time
	Duration         time.Duration
	Description      string
	UserID           uuid.UUID
	NotifyBeforeDays int
}

type Storage interface {
	CreateEvent(ctx context.Context, event *Event) error
	GetEventID(ctx context.Context, event Event) (uuid.UUID, error)
	UpdateEvent(ctx context.Context, event Event) error
	DeleteEvent(ctx context.Context, id uuid.UUID) error
	GetEventsPerDay(ctx context.Context, day time.Time) ([]Event, error) // TODO int?
	GetEventsPerWeek(ctx context.Context, beginDate time.Time) ([]Event, error)
	GetEventsPerMonth(ctx context.Context, beginDate time.Time) ([]Event, error)
}
