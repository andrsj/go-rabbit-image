package client

import (
	"context"
	"fmt"

	"github.com/andrsj/go-rabbit-image/internal/domain/dto"
	"github.com/andrsj/go-rabbit-image/internal/domain/repositories/queue"
	"github.com/andrsj/go-rabbit-image/pkg/logger"
	amqp "github.com/rabbitmq/amqp091-go"
)

type rabbitMQ struct {
	conn      *amqp.Connection
	channel   *amqp.Channel
	MainQueue string
	logger    logger.Logger
}

var _ queue.MessageBroker = (*rabbitMQ)(nil)

func New(url, queue_name string, log logger.Logger) (*rabbitMQ, error) {
	log = log.Named("RabbitMQ client")

	log.Info("Establishing connection to RabbitMQ...", logger.M{"URL": url, "queue_name": queue_name})
	conn, err := amqp.Dial(url)
	if err != nil {
		log.Error("Failed to connect to RabbitMQ", logger.M{"error": err})
		return nil, err
	}
	log.Info("Connection established successfully", nil)

	log.Info("Opening a new channel...", nil)
	channel, err := conn.Channel()
	if err != nil {
		log.Error("Failed to open a new channel", logger.M{"error": err})
		return nil, fmt.Errorf("can't connect to RabbitMQ:%s", err)
	}
	log.Info("New channel opened successfully", nil)

	log.Info("Declaring the main queue...", nil)
	_, err = channel.QueueDeclare(
		queue_name,
		false,
		true,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Error("Failed to declare the main queue", logger.M{"error": err})
		return nil, err
	}
	log.Info("Main queue declared successfully", nil)

	return &rabbitMQ{
		conn:      conn,
		channel:   channel,
		MainQueue: queue_name,
		logger:    log,
	}, nil
}

func (r *rabbitMQ) Publish(ctx context.Context, message []byte, image_id, contentType string) error {
	r.logger.Info("Publishing a message to RabbitMQ", logger.M{
		"queue_name":   r.MainQueue,
		"image_id":     image_id,
		"content_type": contentType,
	})
	err := r.channel.PublishWithContext(ctx,
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
	if err != nil {
		r.logger.Error("Failed to publish a message to RabbitMQ", logger.M{
			"error": err,
		})
		return err
	}
	r.logger.Info("Successfully published a message to RabbitMQ", logger.M{
		"queue_name":   r.MainQueue,
		"image_id":     image_id,
		"content_type": contentType,
	})
	return nil
}

func (r *rabbitMQ) MustConsumeMessages() (<-chan dto.MessageDTO, <-chan error) {
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
		r.logger.Error("Error consuming messages: %v", logger.M{"error": err})
		panic(err)
	}

	messageCh := make(chan dto.MessageDTO)
	errorCh := make(chan error, 1)

	go func() {
		for msg := range msgs {
			r.logger.Info("Received message from RabbitMQ", logger.M{"id": msg.Headers["id"].(string)})

			messageCh <- dto.MessageDTO{
				Body:        msg.Body,
				ImageID:     msg.Headers["id"].(string),
				ContentType: msg.ContentType,
			}
		}
		r.logger.Warn("RabbitMQ channel closed", nil)
		errorCh <- amqp.ErrClosed
	}()

	return messageCh, errorCh
}
