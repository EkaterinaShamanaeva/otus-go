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
		// new storage
		testStorage := New()

		// first event creation
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
		// create new event
		err = testStorage.CreateEvent(context.Background(), &firstEvent)

		require.NoError(t, err)
		require.Equal(t, 1, len(testStorage.mapEvents))

		// second event creation
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
			TimeStart:        timeStartFirstEvent, // busy time
			Duration:         time.Hour,
			Description:      "very important",
			UserID:           userIDSecondEvent,
			NotifyBeforeDays: 1,
		}

		// create second event (time is already taken by the first event)
		err = testStorage.CreateEvent(context.Background(), &secondEvent)

		require.Error(t, err)
		require.ErrorIs(t, err, ErrBusyTime)
		require.Equal(t, 1, len(testStorage.mapEvents))

		// event already exist
		err = testStorage.CreateEvent(context.Background(), &firstEvent)

		require.Error(t, err)
		require.ErrorIs(t, err, ErrAlreadyExist)
		require.Equal(t, 1, len(testStorage.mapEvents))

		// third event creation
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
			TimeStart:        timeStartThirdEvent, // busy time by 1st event at 3 pm (duration 1 hour)
			Duration:         time.Hour,
			Description:      "very important",
			UserID:           userIDThirdEvent,
			NotifyBeforeDays: 1,
		}
		// create event
		err = testStorage.CreateEvent(context.Background(), &thirdEvent)

		require.Error(t, err)
		require.ErrorIs(t, err, ErrBusyTime)
		require.Equal(t, 1, len(testStorage.mapEvents))
	})
	t.Run("get event ID", func(t *testing.T) {
		// new storage
		testStorage := New()

		// first event creation
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

		// get ID of a non-existent event
		resultID, err := testStorage.GetEventID(context.Background(), &firstEvent)

		require.Error(t, err)
		require.ErrorIs(t, err, ErrEventNotExist)
		require.Equal(t, uuid.Nil, resultID)

		// create event
		testStorage.mapEvents[idEvent] = firstEvent

		// get ID
		resultID, err = testStorage.GetEventID(context.Background(), &firstEvent)

		require.NoError(t, err)
		require.Equal(t, idEvent, resultID)
	})

	t.Run("delete event", func(t *testing.T) {
		// new storage
		testStorage := New()

		// first event creation
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

		// create event
		testStorage.mapEvents[idEvent] = firstEvent

		// delete event
		err = testStorage.DeleteEvent(context.Background(), idEvent)

		require.NoError(t, err)
		require.Equal(t, 0, len(testStorage.mapEvents))
	})

	t.Run("update event", func(t *testing.T) {
		// new storage
		testStorage := New()

		// event creation
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

		// event before update
		firstEvent := storage.Event{
			ID:               idEvent,
			Title:            "meeting",
			TimeStart:        timeStartFirstEvent,
			Duration:         time.Hour,
			Description:      "important",
			UserID:           userIDFirstEvent,
			NotifyBeforeDays: 1,
		}

		// event after update
		newEvent := storage.Event{
			ID:               idEvent,
			Title:            "meeting",
			TimeStart:        timeStartFirstEvent,
			Duration:         2 * time.Hour,
			Description:      "very important",
			UserID:           userIDFirstEvent,
			NotifyBeforeDays: 2,
		}

		// update non-existent event
		err = testStorage.UpdateEvent(context.Background(), &newEvent)

		require.Error(t, err)
		require.ErrorIs(t, err, ErrEventNotExist)

		// create event
		testStorage.mapEvents[idEvent] = firstEvent

		// update event
		err = testStorage.UpdateEvent(context.Background(), &newEvent)
		require.NoError(t, err)
	})

	t.Run("get events per day", func(t *testing.T) {
		// new storage
		testStorage := New()

		// first event creation
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

		// second event creation
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

		// create 1st and 2d events
		testStorage.mapEvents[idFirstEvent] = firstEvent
		testStorage.mapEvents[idSecondEvent] = secondEvent

		expectedEventsPerDay := []storage.Event{firstEvent, secondEvent}

		// get all events 8/8/2022
		resultEventsPerDay, err := testStorage.GetEventsPerDay(context.Background(),
			time.Date(2022, 8, 8, 0, 0, 0, 0, time.UTC))

		require.Equal(t, expectedEventsPerDay, resultEventsPerDay)
		require.NoError(t, err)

		// third event creation
		idThirdEvent, err := uuid.NewV4()
		if err != nil {
			fmt.Println(err)
		}

		// new date - 10/8/2022
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

		// create 3d event
		testStorage.mapEvents[idThirdEvent] = thirdEvent

		expectedEventsPerDay = []storage.Event{firstEvent, secondEvent}

		// get all events (8/8/2022)
		resultEventsPerDay, err = testStorage.GetEventsPerDay(context.Background(),
			time.Date(2022, 8, 8, 0, 0, 0, 0, time.UTC))

		require.Equal(t, expectedEventsPerDay, resultEventsPerDay)
	})

	t.Run("get events per week", func(t *testing.T) {
		// new storage
		testStorage := New()

		// first event creation
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

		// second event creation
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

		// create 1st and 2d events
		testStorage.mapEvents[idFirstEvent] = firstEvent
		testStorage.mapEvents[idSecondEvent] = secondEvent

		expectedEventsPerWeek := []storage.Event{firstEvent}

		// get events (8/8/2022 - 14/8/2022)
		resultEventsPerWeek, err := testStorage.GetEventsPerWeek(context.Background(),
			timeStartFirstEvent)

		require.NoError(t, err)
		require.Equal(t, expectedEventsPerWeek, resultEventsPerWeek)
	})

	t.Run("get events per month", func(t *testing.T) {
		// new storage
		testStorage := New()

		// first event creation
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

		// second event creation
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

		// create 1st and 2d events
		testStorage.mapEvents[idFirstEvent] = firstEvent
		testStorage.mapEvents[idSecondEvent] = secondEvent

		expectedEventsPerMonth := []storage.Event{firstEvent, secondEvent}

		// get events (1/8/2022-31/8/2022)
		resultEventsPerMonth, err := testStorage.GetEventsPerMonth(context.Background(),
			timeStartFirstEvent)

		require.NoError(t, err)
		require.Equal(t, expectedEventsPerMonth, resultEventsPerMonth)
	})
}
