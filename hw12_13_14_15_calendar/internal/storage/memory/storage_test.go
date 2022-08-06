package memorystorage

import (
	"context"
	"fmt"
	"github.com/EkaterinaShamanaeva/otus-go/hw12_13_14_15_calendar/internal/storage"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestStorage(t *testing.T) {
	t.Run("create event test", func(t *testing.T) {
		expectedStorage := New()
		idFirstEvent, err := uuid.NewV4()
		if err != nil {
			fmt.Println(err)
		}
		timeStartFirstEvent, err := time.Parse("2/1/2006 3:04 PM", "8/8/2022 3:00 PM")
		if err != nil {
			fmt.Println(err)
		}
		userIDFirstEvent, err := uuid.NewV4()
		if err != nil {
			fmt.Println(err)
		}
		firstEvent := storage.Event{
			ID:               idFirstEvent,
			Title:            "meeting",
			TimeStart:        timeStartFirstEvent,
			Duration:         time.Hour,
			Description:      "important",
			UserID:           userIDFirstEvent,
			NotifyBeforeDays: 1,
		}
		err = expectedStorage.CreateEvent(context.Background(), &firstEvent)
		if err != nil {
			fmt.Println(err)
		}

		require.NoError(t, err)
		require.Equal(t, 1, len(expectedStorage.mapEvents))

		idSecondEvent, err := uuid.NewV4()
		if err != nil {
			fmt.Println(err)
		}
		userIDSecondEvent, err := uuid.NewV4()
		if err != nil {
			fmt.Println(err)
		}
		secondEvent := storage.Event{
			ID:               idSecondEvent,
			Title:            "meeting 2",
			TimeStart:        timeStartFirstEvent,
			Duration:         time.Hour,
			Description:      "very important",
			UserID:           userIDSecondEvent,
			NotifyBeforeDays: 1,
		}
		err = expectedStorage.CreateEvent(context.Background(), &secondEvent)
		if err != nil {
			fmt.Println(err)
		}
		require.Error(t, err)
		require.ErrorIs(t, err, ErrBusyTime)
		require.Equal(t, 1, len(expectedStorage.mapEvents))

		err = expectedStorage.CreateEvent(context.Background(), &firstEvent)
		if err != nil {
			fmt.Println(err)
		}
		require.Error(t, err)
		require.ErrorIs(t, err, ErrAlreadyExist)
		require.Equal(t, 1, len(expectedStorage.mapEvents))

		idThirdEvent, err := uuid.NewV4()
		if err != nil {
			fmt.Println(err)
		}
		timeStartThirdEvent, err := time.Parse("2/1/2006 3:04 PM", "8/8/2022 3:30 PM")
		if err != nil {
			fmt.Println(err)
		}
		userIDThirdEvent, err := uuid.NewV4()
		if err != nil {
			fmt.Println(err)
		}
		thirdEvent := storage.Event{
			ID:               idThirdEvent,
			Title:            "meeting 2",
			TimeStart:        timeStartThirdEvent,
			Duration:         time.Hour,
			Description:      "very important",
			UserID:           userIDThirdEvent,
			NotifyBeforeDays: 1,
		}
		err = expectedStorage.CreateEvent(context.Background(), &thirdEvent)
		if err != nil {
			fmt.Println(err)
		}
		require.Error(t, err)
		require.ErrorIs(t, err, ErrBusyTime)
		require.Equal(t, 1, len(expectedStorage.mapEvents))
	})
	t.Run("get event ID", func(t *testing.T) {
		testStorage := New()
		idEvent, err := uuid.NewV4()
		if err != nil {
			fmt.Println(err)
		}
		timeStartFirstEvent, err := time.Parse("2/1/2006 3:04 PM", "8/8/2022 3:00 PM")
		if err != nil {
			fmt.Println(err)
		}
		userIDFirstEvent, err := uuid.NewV4()
		if err != nil {
			fmt.Println(err)
		}
		firstEvent := storage.Event{
			ID:               idEvent,
			Title:            "meeting",
			TimeStart:        timeStartFirstEvent,
			Duration:         time.Hour,
			Description:      "important",
			UserID:           userIDFirstEvent,
			NotifyBeforeDays: 1,
		}

		resultID, err := testStorage.GetEventID(context.Background(), &firstEvent)
		if err != nil {
			fmt.Println(err)
		}
		require.Error(t, err)
		require.ErrorIs(t, err, ErrEventNotExist)
		require.Equal(t, uuid.Nil, resultID)

		testStorage.mapEvents[idEvent] = firstEvent
		resultID, err = testStorage.GetEventID(context.Background(), &firstEvent)
		if err != nil {
			fmt.Println(err)
		}
		require.NoError(t, err)
		require.Equal(t, idEvent, resultID)
	})

	t.Run("delete event", func(t *testing.T) {
		testStorage := New()
		idEvent, err := uuid.NewV4()
		if err != nil {
			fmt.Println(err)
		}
		timeStartFirstEvent, err := time.Parse("2/1/2006 3:04 PM", "8/8/2022 3:00 PM")
		if err != nil {
			fmt.Println(err)
		}
		userIDFirstEvent, err := uuid.NewV4()
		if err != nil {
			fmt.Println(err)
		}
		firstEvent := storage.Event{
			ID:               idEvent,
			Title:            "meeting",
			TimeStart:        timeStartFirstEvent,
			Duration:         time.Hour,
			Description:      "important",
			UserID:           userIDFirstEvent,
			NotifyBeforeDays: 1,
		}
		testStorage.mapEvents[idEvent] = firstEvent
		err = testStorage.DeleteEvent(context.Background(), idEvent)
		require.NoError(t, err)
		require.Equal(t, 0, len(testStorage.mapEvents))
	})

	t.Run("update event", func(t *testing.T) {
		testStorage := New()
		idEvent, err := uuid.NewV4()
		if err != nil {
			fmt.Println(err)
		}
		timeStartFirstEvent, err := time.Parse("2/1/2006 3:04 PM", "8/8/2022 3:00 PM")
		if err != nil {
			fmt.Println(err)
		}
		userIDFirstEvent, err := uuid.NewV4()
		if err != nil {
			fmt.Println(err)
		}
		firstEvent := storage.Event{
			ID:               idEvent,
			Title:            "meeting",
			TimeStart:        timeStartFirstEvent,
			Duration:         time.Hour,
			Description:      "important",
			UserID:           userIDFirstEvent,
			NotifyBeforeDays: 1,
		}

		newEvent := storage.Event{
			ID:               idEvent,
			Title:            "meeting",
			TimeStart:        timeStartFirstEvent,
			Duration:         2 * time.Hour,
			Description:      "very important",
			UserID:           userIDFirstEvent,
			NotifyBeforeDays: 2,
		}

		err = testStorage.UpdateEvent(context.Background(), &newEvent)
		require.Error(t, err)
		require.ErrorIs(t, err, ErrEventNotExist)

		testStorage.mapEvents[idEvent] = firstEvent

		err = testStorage.UpdateEvent(context.Background(), &newEvent)
		require.NoError(t, err)
	})

	t.Run("get events per day", func(t *testing.T) {
		testStorage := New()
		idFirstEvent, err := uuid.NewV4()
		if err != nil {
			fmt.Println(err)
		}
		timeStartFirstEvent, err := time.Parse("2/1/2006 3:04 PM", "8/8/2022 3:00 PM")
		if err != nil {
			fmt.Println(err)
		}
		userIDFirstEvent, err := uuid.NewV4()
		if err != nil {
			fmt.Println(err)
		}
		firstEvent := storage.Event{
			ID:               idFirstEvent,
			Title:            "meeting",
			TimeStart:        timeStartFirstEvent,
			Duration:         time.Hour,
			Description:      "important",
			UserID:           userIDFirstEvent,
			NotifyBeforeDays: 1,
		}

		idSecondEvent, err := uuid.NewV4()
		if err != nil {
			fmt.Println(err)
		}
		timeStartSecondEvent, err := time.Parse("2/1/2006 3:04 PM", "8/8/2022 5:00 PM")
		if err != nil {
			fmt.Println(err)
		}
		secondEvent := storage.Event{
			ID:               idSecondEvent,
			Title:            "meeting 2",
			TimeStart:        timeStartSecondEvent,
			Duration:         time.Hour,
			Description:      "very important",
			UserID:           userIDFirstEvent,
			NotifyBeforeDays: 1,
		}

		testStorage.mapEvents[idFirstEvent] = firstEvent
		testStorage.mapEvents[idSecondEvent] = secondEvent

		expectedEventsPerDay := []storage.Event{firstEvent, secondEvent}
		resultEventsPerDay, err := testStorage.GetEventsPerDay(context.Background(),
			time.Date(2022, 8, 8, 0, 0, 0, 0, time.UTC))

		require.Equal(t, expectedEventsPerDay, resultEventsPerDay)
		require.NoError(t, err)

		idThirdEvent, err := uuid.NewV4()
		if err != nil {
			fmt.Println(err)
		}
		timeStartThirdEvent, err := time.Parse("2/1/2006 3:04 PM", "10/8/2022 3:30 PM")
		if err != nil {
			fmt.Println(err)
		}
		userIDThirdEvent, err := uuid.NewV4()
		if err != nil {
			fmt.Println(err)
		}
		thirdEvent := storage.Event{
			ID:               idThirdEvent,
			Title:            "meeting 2",
			TimeStart:        timeStartThirdEvent,
			Duration:         time.Hour,
			Description:      "very important",
			UserID:           userIDThirdEvent,
			NotifyBeforeDays: 1,
		}

		testStorage.mapEvents[idThirdEvent] = thirdEvent

		expectedEventsPerDay = []storage.Event{firstEvent, secondEvent}
		resultEventsPerDay, err = testStorage.GetEventsPerDay(context.Background(),
			time.Date(2022, 8, 8, 0, 0, 0, 0, time.UTC))

		require.Equal(t, expectedEventsPerDay, resultEventsPerDay)
	})

	t.Run("get events per week", func(t *testing.T) {
		testStorage := New()
		idFirstEvent, err := uuid.NewV4()
		if err != nil {
			fmt.Println(err)
		}
		timeStartFirstEvent, err := time.Parse("2/1/2006 3:04 PM", "8/8/2022 3:00 PM")
		if err != nil {
			fmt.Println(err)
		}
		userIDFirstEvent, err := uuid.NewV4()
		if err != nil {
			fmt.Println(err)
		}
		firstEvent := storage.Event{
			ID:               idFirstEvent,
			Title:            "meeting",
			TimeStart:        timeStartFirstEvent,
			Duration:         time.Hour,
			Description:      "important",
			UserID:           userIDFirstEvent,
			NotifyBeforeDays: 1,
		}
		idSecondEvent, err := uuid.NewV4()
		if err != nil {
			fmt.Println(err)
		}
		timeStartSecondEvent, err := time.Parse("2/1/2006 3:04 PM", "8/9/2022 5:00 PM")
		if err != nil {
			fmt.Println(err)
		}
		secondEvent := storage.Event{
			ID:               idSecondEvent,
			Title:            "meeting 2",
			TimeStart:        timeStartSecondEvent,
			Duration:         time.Hour,
			Description:      "very important",
			UserID:           userIDFirstEvent,
			NotifyBeforeDays: 1,
		}
		testStorage.mapEvents[idFirstEvent] = firstEvent
		testStorage.mapEvents[idSecondEvent] = secondEvent

		expectedEventsPerWeek := []storage.Event{firstEvent}
		resultEventsPerWeek, err := testStorage.GetEventsPerWeek(context.Background(),
			timeStartFirstEvent)

		require.NoError(t, err)
		require.Equal(t, expectedEventsPerWeek, resultEventsPerWeek)
	})

	t.Run("get events per month", func(t *testing.T) {
		testStorage := New()
		idFirstEvent, err := uuid.NewV4()
		if err != nil {
			fmt.Println(err)
		}
		timeStartFirstEvent, err := time.Parse("2/1/2006 3:04 PM", "8/8/2022 3:00 PM")
		if err != nil {
			fmt.Println(err)
		}
		userIDFirstEvent, err := uuid.NewV4()
		if err != nil {
			fmt.Println(err)
		}
		firstEvent := storage.Event{
			ID:               idFirstEvent,
			Title:            "meeting",
			TimeStart:        timeStartFirstEvent,
			Duration:         time.Hour,
			Description:      "important",
			UserID:           userIDFirstEvent,
			NotifyBeforeDays: 1,
		}
		idSecondEvent, err := uuid.NewV4()
		if err != nil {
			fmt.Println(err)
		}
		timeStartSecondEvent, err := time.Parse("2/1/2006 3:04 PM", "16/8/2022 5:00 PM")
		if err != nil {
			fmt.Println(err)
		}
		secondEvent := storage.Event{
			ID:               idSecondEvent,
			Title:            "meeting 2",
			TimeStart:        timeStartSecondEvent,
			Duration:         time.Hour,
			Description:      "very important",
			UserID:           userIDFirstEvent,
			NotifyBeforeDays: 1,
		}
		testStorage.mapEvents[idFirstEvent] = firstEvent
		testStorage.mapEvents[idSecondEvent] = secondEvent

		expectedEventsPerMonth := []storage.Event{firstEvent, secondEvent}
		resultEventsPerMonth, err := testStorage.GetEventsPerMonth(context.Background(),
			timeStartFirstEvent)

		require.NoError(t, err)
		require.Equal(t, expectedEventsPerMonth, resultEventsPerMonth)
	})
}
