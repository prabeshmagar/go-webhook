package messagebroker

import (
	"context"
	"fmt"

	"github.com/streadway/amqp"
	"gitlab.wlink.com.np/nettv-webhook/internal/config"
)

type RabbitMqClient struct {
	Connection *amqp.Connection
	Exchange   string
	Queue      amqp.Queue
	Channel    *amqp.Channel
}

type MessageBrokerService interface {
	Publish(msg []byte) error
	PublishOnQueue(msg []byte, queueName string) error
	Close()
}

func NewRabbitMqClient(C config.RabbitMQ) MessageBrokerService {
	var ctx context.Context
	var r RabbitMqClient
	var err error

	connectionstring := fmt.Sprintf("amqp://%s:%s@%s:%s", C.Username, C.Password, C.HostName, C.Port)

	r.Connection, err = amqp.Dial(connectionstring)

	failOnError(ctx, err, "Failed to connect to RabbitMQ")
	defer r.Connection.Close()

	r.Channel, err = r.Connection.Channel()
	failOnError(ctx, err, "Failed to open a channel")

	err = r.Channel.ExchangeDeclare(
		C.ExchangeName,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(ctx, err, "Failed to declare an exchange")

	r.Queue, err = r.Channel.QueueDeclare(
		C.QueueName, // name
		false,       // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	failOnError(ctx, err, "Failed to declare an queue")

	err = r.Channel.QueueBind(
		r.Queue.Name,   // name of the queue
		"",             // bindingKey
		C.ExchangeName, // sourceExchange
		false,          // noWait
		nil,            // arguments
	)

	failOnError(ctx, err, "Failed to bind queue")

	return &r
}

func (r *RabbitMqClient) Publish(msg []byte) error {
	fmt.Printf("Published data %v", string(msg))
	err := r.Channel.Publish(r.Exchange, r.Queue.Name, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        msg,
	})

	return err
}

func (r *RabbitMqClient) PublishOnQueue(msg []byte, queueName string) error {
	fmt.Println(msg)
	return nil

}

func (r *RabbitMqClient) Close() {
	fmt.Println(r)
}

func failOnError(ctx context.Context, err error, msg string) {
	if err != nil {
		// level.Error(log).Log("msg", msg, "err", err)
		fmt.Println(err)
	}

}
