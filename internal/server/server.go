package server

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"affise-server/internal/config"
	"affise-server/internal/handlers"
	"affise-server/internal/middleware"
)

const (
	EndPointRequests = "/requests"

	DefaultReadHeaderTimeout = 5 * time.Second
)

// Server represents server.
type Server struct {
	config *config.Config
	server *http.Server
}

// NewServer creates new instance of server.
func NewServer(conf *config.Config) *Server {
	var s = &Server{
		config: conf,
		server: &http.Server{
			Addr:              conf.Server.Addr,
			ReadHeaderTimeout: DefaultReadHeaderTimeout,
		},
	}

	rateMV := middleware.NewRateLimit(conf.Server.RateInterval, conf.Server.RateLimit)
	requests := handlers.NewRequestsHandler(&conf.Client)

	mux := http.NewServeMux()
	mux.Handle(EndPointRequests, rateMV.Handle(requests))
	s.server.Handler = mux

	return s
}

// Run runs server.
func (s *Server) Run(ctx context.Context) error {
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		<-ctx.Done()
		log.Print("server is shutting down")

		if err := s.server.Shutdown(ctx); err != nil {
			log.Printf("server shutdown err: %v", err)
		}
	}()

	log.Printf("server is listening on %s", s.server.Addr)

	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	wg.Wait()

	return nil
}

// Addr returns the server address.
func (s *Server) Addr() string {
	return s.server.Addr
}
