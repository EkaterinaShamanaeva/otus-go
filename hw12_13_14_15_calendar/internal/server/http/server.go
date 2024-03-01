package internalhttp

import (
	"context"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"time"
)

type Server struct { // TODO: add Application
	server *http.Server
	router *httprouter.Router
	logg   Logger
}

type Logger interface {
	Info(msg string)
	Error(msg string)
	Warn(msg string)
	Debug(msg string)
}

type Application interface { // TODO implement
}

func NewServer(logger Logger) *Server { // app Application
	serv := &Server{logg: logger}
	serv.router = httprouter.New()

	serv.router.GET("/", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		text := "Hello world!"
		fmt.Fprint(writer, text)
	})

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
