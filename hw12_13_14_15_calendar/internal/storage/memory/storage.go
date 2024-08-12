package memorystorage

import (
	"context"
	"errors"
	"github.com/EkaterinaShamanaeva/otus-go/hw12_13_14_15_calendar/internal/storage"
	"github.com/gofrs/uuid"
	"sync"
	"time"
)

var (
	ErrBusyTime          = errors.New("this time is already taken by another event")
	ErrAlreadyExist      = errors.New("this event has already exist")
	ErrCanceledByContext = errors.New("context cancel")
	ErrEventNotExist     = errors.New("event does not exist")
)

type Storage struct {
	mu        sync.Mutex
	mapEvents map[uuid.UUID]storage.Event
}

func New() *Storage {
	return &Storage{
		mu:        sync.Mutex{},
		mapEvents: make(map[uuid.UUID]storage.Event),
	}
}

func (s *Storage) CreateEvent(ctx context.Context, event *storage.Event) error {
	select {
	case <-ctx.Done():
		return ErrCanceledByContext
	default:
		s.mu.Lock()
		defer s.mu.Unlock()
		for _, evValue := range s.mapEvents {
			if event.TimeStart.Equal(evValue.TimeStart) && event.ID != evValue.ID {
				return ErrBusyTime
			} else if event.TimeStart.Equal(evValue.TimeStart) && event.ID == evValue.ID {
				return ErrAlreadyExist
			} else if evValue.TimeStart.Sub(event.TimeStart) < evValue.Duration {
				return ErrBusyTime
			}
		}
		s.mapEvents[event.ID] = *event
	}
	return nil
}

func (s *Storage) GetEventID(ctx context.Context, event *storage.Event) (uuid.UUID, error) {
	select {
	case <-ctx.Done():
		return uuid.Nil, ErrCanceledByContext
	default:
		s.mu.Lock()
		defer s.mu.Unlock()
		if _, ok := s.mapEvents[event.ID]; ok {
			return event.ID, nil
		}
	}
	return uuid.Nil, ErrEventNotExist
}

func (s *Storage) DeleteEvent(ctx context.Context, id uuid.UUID) error {
	select {
	case <-ctx.Done():
		return ErrCanceledByContext
	default:
		s.mu.Lock()
		defer s.mu.Unlock()
		if _, ok := s.mapEvents[id]; !ok {
			return ErrEventNotExist
		}
		delete(s.mapEvents, id)
	}
	return nil
}

func (s *Storage) UpdateEvent(ctx context.Context, event *storage.Event) error {
	select {
	case <-ctx.Done():
		return ErrCanceledByContext
	default:
		s.mu.Lock()
		defer s.mu.Unlock()
		if _, ok := s.mapEvents[event.ID]; !ok {
			return ErrEventNotExist
		}
		s.mapEvents[event.ID] = *event
	}
	return nil
}

func (s *Storage) GetEventsPerDay(ctx context.Context, day time.Time) ([]storage.Event, error) {
	eventsPerDay := make([]storage.Event, 0)
	select {
	case <-ctx.Done():
		return nil, ErrCanceledByContext
	default:
		s.mu.Lock()
		defer s.mu.Unlock()
		for _, eventStruct := range s.mapEvents {
			if eventStruct.TimeStart.Year() == day.Year() && eventStruct.TimeStart.Month() == day.Month() &&
				eventStruct.TimeStart.Day() == day.Day() {
				eventsPerDay = append(eventsPerDay, eventStruct)
			}
		}
	}
	return eventsPerDay, nil
}

func (s *Storage) GetEventsPerWeek(ctx context.Context, beginDate time.Time) ([]storage.Event, error) {
	eventsPerWeek := make([]storage.Event, 0)
	select {
	case <-ctx.Done():
		return nil, ErrCanceledByContext
	default:
		s.mu.Lock()
		defer s.mu.Unlock()
		endDay := beginDate.AddDate(0, 0, 7)
		for _, eventStruct := range s.mapEvents {
			if (eventStruct.TimeStart.After(beginDate) || eventStruct.TimeStart.Equal(beginDate)) &&
				(eventStruct.TimeStart.Before(endDay) || eventStruct.TimeStart.Equal(endDay)) {
				eventsPerWeek = append(eventsPerWeek, eventStruct)
			}
		}
	}
	return eventsPerWeek, nil
}

func (s *Storage) GetEventsPerMonth(ctx context.Context, beginDate time.Time) ([]storage.Event, error) {
	eventsPerMonth := make([]storage.Event, 0)
	select {
	case <-ctx.Done():
		return nil, ErrCanceledByContext
	default:
		s.mu.Lock()
		defer s.mu.Unlock()
		endDay := beginDate.AddDate(0, 1, 0)
		for _, eventStruct := range s.mapEvents {
			if (eventStruct.TimeStart.After(beginDate) || eventStruct.TimeStart.Equal(beginDate)) &&
				(eventStruct.TimeStart.Before(endDay) || eventStruct.TimeStart.Equal(endDay)) {
				eventsPerMonth = append(eventsPerMonth, eventStruct)
			}
		}
	}
	return eventsPerMonth, nil
}

func (s *Storage) Close(ctx context.Context) error {
	return nil
}

func (s *Storage) ListForScheduler(ctx context.Context, remindFor time.Duration, period time.Duration) ([]storage.Notification, error) {
	return nil, nil
}

func (s *Storage) ClearEvents(ctx context.Context) error {
	return nil
}
