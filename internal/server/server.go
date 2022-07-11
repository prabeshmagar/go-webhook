package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-kit/log"
	"gitlab.wlink.com.np/nettv-webhook/internal/route"
	"gitlab.wlink.com.np/nettv-webhook/messagebroker"
)

type Server struct {
	address string
	log     log.Logger
	server  *http.Server
}

type Options struct {
	Host      string
	Port      int
	MsgBroker messagebroker.MessageBrokerService
}

func New(opts Options) *Server {
	var logger log.Logger
	logger = log.NewJSONLogger(os.Stderr)
	logger = log.NewSyncLogger(logger)

	address := net.JoinHostPort(opts.Host, strconv.Itoa(opts.Port))
	return &Server{
		address: address,
		server: &http.Server{
			Addr:              address,
			Handler:           route.SetupRoutes(opts.MsgBroker, logger),
			ReadTimeout:       5 * time.Second,
			ReadHeaderTimeout: 5 * time.Second,
			WriteTimeout:      5 * time.Second,
			IdleTimeout:       5 * time.Second,
		},
	}
}

func (s *Server) Start() error {

	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("error starting server: %w", err)

	}
	return nil
}

// Stop the Server gracefully within the timeout.
func (s *Server) Stop() error {
	fmt.Println("Stopping")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {

		return fmt.Errorf("error stopping server: %w", err)
	}

	return nil
}
