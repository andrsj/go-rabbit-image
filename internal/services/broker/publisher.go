package broker

import (
	"context"
	"fmt"

	"github.com/andrsj/go-rabbit-image/internal/delivery/rabbitmq/client"
	"github.com/andrsj/go-rabbit-image/internal/domain/repositories/queue"
)

type messagePublisherService struct {
	publisher queue.Publisher
}

var _ queue.Publisher = (*messagePublisherService)(nil)

func NewPublisher(url, queue_name string) (*messagePublisherService, error) {
	client, err := client.New(url, queue_name)
	if err != nil {
		return nil, fmt.Errorf("can't start message broker: %s", err)
	}
	return &messagePublisherService{
		publisher: client,
	}, nil
}

func (m *messagePublisherService) Publish(ctx context.Context, message []byte, image_id string, contentType string) error {
	return m.publisher.Publish(ctx, message, image_id, contentType)
}
