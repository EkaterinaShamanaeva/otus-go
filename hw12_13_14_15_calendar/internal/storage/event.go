package storage

import (
	"github.com/gofrs/uuid"
	"golang.org/x/net/context"
	"time"
)

type Event struct {
	ID               uuid.UUID     `db:"id"`
	Title            string        `db:"title"`
	TimeStart        time.Time     `db:"start_date"`
	Duration         time.Duration `db:"duration"`
	Description      string        `db:"description"`
	UserID           uuid.UUID     `db:"user_id"`
	NotifyBeforeDays int           `db:"notify_before"`
}

type Storage interface {
	CreateEvent(ctx context.Context, event *Event) error
	GetEventID(ctx context.Context, event *Event) (uuid.UUID, error)
	UpdateEvent(ctx context.Context, event *Event) error
	DeleteEvent(ctx context.Context, id uuid.UUID) error
	GetEventsPerDay(ctx context.Context, day time.Time) ([]Event, error)
	GetEventsPerWeek(ctx context.Context, beginDate time.Time) ([]Event, error)
	GetEventsPerMonth(ctx context.Context, beginDate time.Time) ([]Event, error)
}
