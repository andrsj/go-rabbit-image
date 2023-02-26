package broker

import (
	"fmt"

	"github.com/andrsj/go-rabbit-image/internal/delivery/rabbitmq/client"
	"github.com/andrsj/go-rabbit-image/internal/domain/repositories/queue"
)

type consumerService struct {
	consumer queue.Consumer
}

var _ queue.Consumer = (*consumerService)(nil)

func NewConsumer(url, queue_name string) (*consumerService, error) {
	client, err := client.New(url, queue_name)
	if err != nil {
		return nil, fmt.Errorf("can't start message broker: %s", err)
	}
	return &consumerService{
		consumer: client,
	}, nil
}

func (c *consumerService) ConsumeMessages() (<-chan []byte, <-chan error) {
	return c.consumer.ConsumeMessages()
}
