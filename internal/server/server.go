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

	"affise-server/internal/handlers"
	"affise-server/internal/middleware"
)

const (
	EndPointRequests = "/requests"

	DefaultAddr         = ":8080"
	DefaultRateInterval = 1 * time.Second
	DefaultRateLimit    = 100
)

// Option is a function that configures the server.
type Option func(*Server)

// WithServer is an option to set the http.Server.
func WithServer(s *http.Server) Option {
	return func(srv *Server) {
		srv.server = s
	}
}

// WithAddr is an option to set the server address.
func WithAddr(addr string) Option {
	return func(srv *Server) {
		srv.addr = addr
	}
}

type Server struct {
	addr   string
	server *http.Server
}

func NewServer(opts ...Option) *Server {
	var s = &Server{}
	for _, opt := range opts {
		opt(s)
	}

	if s.addr == "" {
		s.addr = DefaultAddr
	}

	if s.server == nil {
		s.server = &http.Server{
			Addr: s.addr,
		}
	}

	limiter := middleware.NewRateLimit(DefaultRateInterval, DefaultRateLimit)

	mux := http.NewServeMux()
	mux.Handle(EndPointRequests, limiter.Handle(handlers.NewRequestsHandler()))
	s.server.Handler = mux

	return s
}

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
