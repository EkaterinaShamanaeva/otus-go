package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/EkaterinaShamanaeva/otus-go/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/EkaterinaShamanaeva/otus-go/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
)

var (
	ErrBusyTime      = errors.New("this time is already taken by another event")
	ErrAlreadyExist  = errors.New("this event has already exist")
	ErrEventNotExist = errors.New("event does not exist")
)

type Storage struct {
	Pool *pgxpool.Pool
}

func New() *Storage {
	return &Storage{Pool: nil}
}

func Connect(ctx context.Context, dsn string) (dbpool *pgxpool.Pool, err error) {
	conn, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		err = fmt.Errorf("failed to parse pg config: %w", err)
		return
	}
	// Config
	conn.MaxConns = int32(5)
	conn.MinConns = int32(1)
	conn.HealthCheckPeriod = 1 * time.Minute
	conn.MaxConnLifetime = 24 * time.Hour
	conn.MaxConnIdleTime = 30 * time.Minute
	conn.ConnConfig.ConnectTimeout = 1 * time.Second
	// Create pool
	dbpool, err = pgxpool.ConnectConfig(ctx, conn) // pool
	if err != nil {
		err = fmt.Errorf("failed to connect config: %w", err)
		return
	}
	return
}

func (s *Storage) CreateEvent(ctx context.Context, event *storage.Event) error {
	events, err := s.GetEventsPerDay(ctx, event.TimeStart)
	for _, ev := range events {
		if ev.ID == event.ID {
			return ErrAlreadyExist
		} else if ev.TimeStart == event.TimeStart || ev.TimeStart.Sub(event.TimeStart) < ev.Duration {
			return ErrBusyTime
		}
	}
	query := `INSERT INTO events(id, title, start_date, duration, description, user_id, notify_before)
				VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING id; `
	_, err = s.Pool.Exec(ctx, query, event.ID, event.Title, event.TimeStart, event.Duration,
		event.Description, event.UserID, event.NotifyBeforeDays)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) GetEventID(ctx context.Context, event *storage.Event) (uuid.UUID, error) {
	query := `SELECT * FROM events WHERE id = $1;`

	var events []*storage.Event
	err := pgxscan.Select(ctx, s.Pool, &events, query, event.ID)
	if err != nil {
		return uuid.Nil, err
	}
	return events[0].ID, nil
}

func (s *Storage) UpdateEvent(ctx context.Context, event *storage.Event) error {
	query := `UPDATE events SET id=$1, title=$2, start_date=$3, duration=$4, description=$5,
				user_id=$6, notify_before=$7 WHERE id=$8;`
	_, err := s.Pool.Exec(ctx, query, event.ID, event.Title, event.TimeStart, event.Duration,
		event.Description, event.UserID, event.NotifyBeforeDays, event.ID)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) DeleteEvent(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM events WHERE id=$1`
	_, err := s.Pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) GetEventsPerDay(ctx context.Context, day time.Time) ([]storage.Event, error) {
	query := `SELECT * FROM events WHERE start_date BETWEEN $1 AND $1 + (interval '1d');`
	var eventsByQuery []*storage.Event
	err := pgxscan.Select(ctx, s.Pool, &eventsByQuery, query, day)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrEventNotExist
		}
		return nil, err
	}
	eventsRes := make([]storage.Event, 0, len(eventsByQuery))
	for _, ev := range eventsByQuery {
		eventsRes = append(eventsRes, *ev)
	}
	return eventsRes, nil
}

func (s *Storage) GetEventsPerWeek(ctx context.Context, beginDate time.Time) ([]storage.Event, error) {
	query := `SELECT * FROM events WHERE start_date BETWEEN $1 AND $1 + (interval '7d');`
	var eventsByQuery []*storage.Event
	err := pgxscan.Select(ctx, s.Pool, &eventsByQuery, query, beginDate)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, memorystorage.ErrEventNotExist
		}
		return nil, err
	}
	eventsRes := make([]storage.Event, 0, len(eventsByQuery))
	for _, ev := range eventsByQuery {
		eventsRes = append(eventsRes, *ev)
	}
	return eventsRes, nil
}

func (s *Storage) GetEventsPerMonth(ctx context.Context, beginDate time.Time) ([]storage.Event, error) {
	query := `SELECT * FROM events WHERE start_date BETWEEN $1 AND $1 + (interval '1months');`
	var eventsByQuery []*storage.Event
	err := pgxscan.Select(ctx, s.Pool, &eventsByQuery, query, beginDate)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, memorystorage.ErrEventNotExist
		}
		return nil, err
	}
	eventsRes := make([]storage.Event, 0, len(eventsByQuery))
	for _, ev := range eventsByQuery {
		eventsRes = append(eventsRes, *ev)
	}
	return eventsRes, nil
}

func (s *Storage) Close(ctx context.Context) error {
	s.Pool.Close()
	return nil
}

func (s *Storage) ListForScheduler(ctx context.Context, remindFor time.Duration, period time.Duration) ([]storage.Notification, error) {
	from := time.Now().Add(remindFor)
	fmt.Println(from)
	to := from.Add(period)
	fmt.Println(to)

	query := `SELECT id, title, start_date, user_id FROM events WHERE start_date BETWEEN $1 AND $2;`
	var notices []*storage.Notification
	err := pgxscan.Select(ctx, s.Pool, &notices, query, from.Format("2006-01-02 15:04:00 -0700"),
		to.Format("2006-01-02 15:04:00 -0700"))
	if err != nil {
		return nil, err
	}
	res := make([]storage.Notification, 0, len(notices))
	for _, ev := range notices {
		res = append(res, *ev)
	}
	fmt.Println(res)
	return res, nil
}

func (s *Storage) ClearEvents(ctx context.Context) error {
	hourInDay := 24 * time.Hour
	daysInYear := 365 * hourInDay
	yNow := time.Now().Year()
	if yNow%4 == 0 && (yNow%100 != 0 || yNow%400 == 0) {
		daysInYear = 366 * hourInDay
	}
	yearAgoStr := time.Now().Add(-1 * daysInYear).Format("2006-01-02 15:04:00 -0700")

	_, err := s.Pool.Exec(ctx, `DELETE FROM events WHERE start_date < $1`, yearAgoStr)
	if err != nil {
		return err
	}
	fmt.Println("events cleaned")
	return nil
}
