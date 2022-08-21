package internalhttp

import (
	"context"
	"github.com/EkaterinaShamanaeva/otus-go/hw12_13_14_15_calendar/internal/app"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

type Server struct {
	server *http.Server
	router *mux.Router //*httprouter.Router
	logg   Logger
	app    *app.App
}

type Logger interface {
	Info(msg string)
	Error(msg string)
	Warn(msg string)
	Debug(msg string)
}

type Application interface { // TODO implement
}

func NewServer(logger Logger, app *app.App) *Server { // app Application
	serv := &Server{logg: logger, app: app}
	// serv.router = httprouter.New() //TODO
	serv.router = mux.NewRouter()

	serv.router.HandleFunc("/create_event", serv.createEvent).Methods("PUT")
	serv.router.HandleFunc("/update_event", serv.updateEvent).Methods("PUT")
	serv.router.HandleFunc("/delete_event", serv.deleteEvent).Methods("DELETE")
	serv.router.HandleFunc("/get_events_per_day", serv.getEventsPerDay).Methods("GET")
	serv.router.HandleFunc("/get_events_per_week", serv.getEventsPerWeek).Methods("GET")
	serv.router.HandleFunc("/get_events_per_month", serv.getEventsPerMonth).Methods("GET")

	//serv.router.PUT("/create_event", serv.createEvent)
	//serv.router.PUT("/update_event", serv.updateEvent)
	//serv.router.DELETE("/delete_event", serv.deleteEvent)
	//serv.router.GET("/get_events_per_day", serv.getEventsPerDay)
	//serv.router.GET("/get_events_per_week", serv.getEventsPerWeek)
	//serv.router.GET("/get_events_per_month", serv.getEventsPerMonth)

	return serv
}

func (s *Server) Start(ctx context.Context, addr string) error {
	s.logg.Info("HTTP server " + addr + " starting...")
	s.server = &http.Server{
		Addr:         addr,
		Handler:      s.router,
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}

	errChan := make(chan error)

	go func() {
		if err := s.server.ListenAndServe(); err != nil {
			errChan <- err
		}
	}()

	select {
	case <-ctx.Done():
	case err := <-errChan:
		return err
	}
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.logg.Info("HTTP server was stopped...")
	if err := s.server.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}
