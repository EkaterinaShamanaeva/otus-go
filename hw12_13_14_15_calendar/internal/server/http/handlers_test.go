package internalhttp

import (
	"bytes"
	"context"
	"fmt"
	"github.com/EkaterinaShamanaeva/otus-go/hw12_13_14_15_calendar/internal/app"
	"github.com/EkaterinaShamanaeva/otus-go/hw12_13_14_15_calendar/internal/logger"
	"github.com/EkaterinaShamanaeva/otus-go/hw12_13_14_15_calendar/internal/storage/init_storage"
	"github.com/steinfletcher/apitest"
	"github.com/stretchr/testify/assert"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestCreateEvent(t *testing.T) {
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

	jsonFile, _ := os.Open("event.json")

	req := httptest.NewRequest(http.MethodPut, "/create_event", jsonFile)
	w := httptest.NewRecorder()

	server.createEvent(w, req)
	resp := w.Result()
	data, _ := io.ReadAll(resp.Body)
	fmt.Println("res1: ", data)

	req = httptest.NewRequest(http.MethodGet, "/get_events_per_day",
		bytes.NewBuffer([]byte("\"2022-10-02T13:00:00Z\"")))
	w = httptest.NewRecorder()

	server.getEventsPerDay(w, req)
	resp = w.Result()
	data, _ = io.ReadAll(resp.Body)
	fmt.Println("res2: ", string(data))
	fmt.Println(resp.StatusCode)
	fmt.Println(resp.Header.Get("Content-Type"))

	apitest.New().
		HandlerFunc(server.createEvent).
		Put("/create_event").
		JSONFromFile("../../../internal/server/http/event.json").
		Expect(t).
		Assert(func(res *http.Response, req *http.Request) error {
			assert.Equal(t, http.StatusOK, res.StatusCode)
			return nil
		}).
		End()

	apitest.New().
		HandlerFunc(server.getEventsPerDay).
		Get("/get_events_per_day").
		Body("\"2022-08-08T12:00:00Z\"").
		Expect(t).
		Assert(func(res *http.Response, req *http.Request) error {
			assert.Equal(t, http.StatusOK, res.StatusCode)
			return nil
		}).
		End()
}
