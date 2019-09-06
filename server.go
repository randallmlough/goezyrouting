package goezyrouting

import (
	"context"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type Server struct {
	Logger
	*http.Server
}

func NewServer(opts ...ServerOptions) *Server {
	s := new(Server)
	for _, opt := range opts {
		opt(s)
	}
	if s.Server == nil {
		s.Server = defaultServer
	}
	return s
}

type ServerOptions func(s *Server)

var defaultServer = &http.Server{
	Addr:         ":3000",
	ReadTimeout:  5 * time.Second,
	WriteTimeout: 10 * time.Second,
	IdleTimeout:  15 * time.Second,
}

func DefaultServer() ServerOptions {
	return func(s *Server) {
		s.Server = defaultServer
	}
}
func WithHandler(h http.Handler) ServerOptions {
	return func(s *Server) {
		s.Server.Handler = h
	}
}
func WithPort(p string) ServerOptions {
	return func(s *Server) {
		s.Server.Addr = ":" + p
	}
}
func WithErrorLog(l *log.Logger) ServerOptions {
	return func(s *Server) {
		s.Server.ErrorLog = l
	}
}
func WithLogger(l Logger) ServerOptions {
	return func(s *Server) {
		s.Logger = l
	}
}

// Start starts the http server
func (s *Server) Start() {
	s.Log("Starting server...")

	go func() {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.Fatal("failed to start server", err)
		}
	}()
	s.Log("server is ready to handle requests on:", s.Addr)
	s.GracefullShutdown()
}

// StartTLS starts the https server
func (s *Server) StartTLS(cf, kf string) error {
	go func() {
		s.Log("starting server on port %v", s.Addr)
		if err := s.ListenAndServeTLS(cf, kf); err != nil {
			s.Fatal("failed to start server", err)
		}
	}()
	// Setting up signal capturing
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	// Waiting for SIGINT (pkill -2)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		return errors.Wrap(err, "error stopping server")
	}

	s.Log("gracefully stopped server")
	return nil
}

func (s *Server) GracefullShutdown() {
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt)
	sig := <-quit
	s.Log("server is shutting down", "reason:", sig.String())

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	s.SetKeepAlivesEnabled(false)
	if err := s.Shutdown(ctx); err != nil {
		s.Fatal("Could not gracefully shutdown the server", "error", err)
	}
	s.Log("server stopped")
}
