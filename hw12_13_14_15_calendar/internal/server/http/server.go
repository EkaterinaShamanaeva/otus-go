package internalhttp

import (
	"context"
	"github.com/EkaterinaShamanaeva/otus-go/hw12_13_14_15_calendar/internal/logger"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"time"
)

type Server struct { // TODO
	server *http.Server
	router *httprouter.Router
	logg   *logger.Logger
}

// type Logger interface { // TODO
// }

// type Application interface { // TODO
// }

func NewServer(logger *logger.Logger) *Server { // app Application
	serv := &Server{logg: logger}
	serv.router = httprouter.New()

	serv.router.GET("/", StartPage)

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
	// <-ctx.Done()
	if err := s.server.ListenAndServe(); err != nil {
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

// TODO
