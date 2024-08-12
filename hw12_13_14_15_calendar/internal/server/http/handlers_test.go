package internalhttp

import (
	"context"
	"fmt"
	"github.com/EkaterinaShamanaeva/otus-go/hw12_13_14_15_calendar/internal/app"
	"github.com/EkaterinaShamanaeva/otus-go/hw12_13_14_15_calendar/internal/logger"
	"github.com/EkaterinaShamanaeva/otus-go/hw12_13_14_15_calendar/internal/storage/init_storage"
	"github.com/steinfletcher/apitest"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"testing"
)

func TestAPICalendar(t *testing.T) {
	logg, err := logger.New("DEBUG", "logfile.log")
	if err != nil {
		log.Fatalf("Config error: %v", err)
	}

	ctx := context.Background()

	db, err := init_storage.NewStorage(ctx, "memory",
		"postgres://postgres:password@localhost:5432/calendar?sslmode=disable")
	if err != nil {
		logg.Error("failed to connect DB: " + err.Error())
	}
	defer db.Close(ctx)

	calendar := app.New(logg, db)

	server := NewServer(logg, calendar)

	t.Run("create event", func(t *testing.T) {
		apitest.New().
			HandlerFunc(server.createEvent).
			Put("/create_event").
			JSONFromFile("../../../internal/server/http/tests/event.json").
			Expect(t).
			Status(http.StatusOK).
			Assert(func(res *http.Response, req *http.Request) error {
				fmt.Println("request: ", req.Body)
				assert.Equal(t, http.StatusOK, res.StatusCode)
				return nil
			}).
			End()

		apitest.New().
			HandlerFunc(server.getEventsPerDay).
			Get("/get_events_per_day").
			Body("\"2022-10-02T12:00:00Z\"").
			Expect(t).
			Assert(func(res *http.Response, req *http.Request) error {
				fmt.Println("request: ", req.Body)
				fmt.Println("result: ", res.Body)
				assert.Equal(t, http.StatusOK, res.StatusCode)
				return nil
			}).
			End()

		apitest.New().
			HandlerFunc(server.getEventsPerWeek).
			Delete("/get_events_per_week").
			Body("\"2022-10-02T12:00:00Z\"").
			Expect(t).
			Assert(func(res *http.Response, req *http.Request) error {
				fmt.Println("request: ", req.Body)
				fmt.Println("result: ", res.Body)
				assert.Equal(t, http.StatusOK, res.StatusCode)
				return nil
			}).
			End()
	})
}
