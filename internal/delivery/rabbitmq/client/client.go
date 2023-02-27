package client

import (
	"context"
	"fmt"

	"github.com/andrsj/go-rabbit-image/internal/domain/repositories/queue"
	"github.com/andrsj/go-rabbit-image/internal/infrastructure/worker"
	amqp "github.com/rabbitmq/amqp091-go"
)

type rabbitMQ struct {
	conn      *amqp.Connection
	channel   *amqp.Channel
	MainQueue string
}

type MessageBroker interface {
	queue.Publisher
	worker.Consumer
}

var _ MessageBroker = (*rabbitMQ)(nil)

func New(url, queue_name string) (*rabbitMQ, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("can't connect to RabbitMQ:%s", err)
	}

	_, err = channel.QueueDeclare(
		queue_name,
		false,
		true,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &rabbitMQ{
		conn:      conn,
		channel:   channel,
		MainQueue: queue_name,
	}, nil
}

func (r *rabbitMQ) Publish(ctx context.Context, message []byte, image_id, contentType string) error {
	return r.channel.PublishWithContext(ctx,
		"",
		r.MainQueue,
		false,
		false,
		amqp.Publishing{
			Headers: map[string]interface{}{
				"id": image_id,
			},
			ContentType: contentType,
			Body:        message,
		},
	)
}

func (r *rabbitMQ) ConsumeMessages() (<-chan []byte, <-chan error) {
	msgs, err := r.channel.Consume(
		r.MainQueue,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic(err)
	}

	messageCh := make(chan []byte)
	errorCh := make(chan error, 1)

	go func() {
		for msg := range msgs {
			messageCh <- msg.Body
		}
		errorCh <- amqp.ErrClosed
	}()

	return messageCh, errorCh
}
