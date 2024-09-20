package server

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"log/slog"
	"net/http"
	"url-shortener/config"
)

type Server struct {
	httpServer *http.Server
	log        *slog.Logger
}

func NewServer(cfg *config.ServerConfig, router *mux.Router, log *slog.Logger) (*Server, error) {
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}
	return &Server{
		httpServer: srv,
		log:        log,
	}, nil
}

func (s *Server) Run(ctx context.Context) error {
	errChan := make(chan error)
	go func() {
		s.log.Info(fmt.Sprintf("starting listening: %s", s.httpServer.Addr))

		errChan <- s.httpServer.ListenAndServe()
	}()

	var err error
	select {
	case <-ctx.Done():
		return ctx.Err()
	case err = <-errChan:

	}
	return err
}

func (s *Server) Shutdown(ctx context.Context) error {
	log.Println("Shutting down server...")
	return s.httpServer.Shutdown(ctx)
}
