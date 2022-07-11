package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/log"
	"gitlab.wlink.com.np/nettv-webhook/internal/config"
	"gitlab.wlink.com.np/nettv-webhook/internal/server"
	"gitlab.wlink.com.np/nettv-webhook/messagebroker"
	"golang.org/x/sync/errgroup"
)

var logger log.Logger

func main() {
	C := config.NewConfig()

	rmq := messagebroker.NewRabbitMqClient(C.RabbitMQ)

	os.Exit(start(rmq))

}

func start(msgBrk messagebroker.MessageBrokerService) int {
	host := ""
	port := 8080

	s := server.New(server.Options{
		Port:      port,
		Host:      host,
		MsgBroker: msgBrk,
	})

	var eg errgroup.Group

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	eg.Go(func() error {
		<-ctx.Done()

		if err := s.Stop(); err != nil {
			return err
		}

		return nil
	})

	if err := s.Start(); err != nil {
		return 1
	}

	if err := eg.Wait(); err != nil {
		return 1
	}

	return 0
}
