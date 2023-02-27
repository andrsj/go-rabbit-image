package publisher

import (
	"context"

	"github.com/andrsj/go-rabbit-image/internal/domain/repositories/queue"
)

type messagePublisherService struct {
	publisher queue.Publisher
}

var _ queue.Publisher = (*messagePublisherService)(nil)

func New(publisher queue.Publisher) *messagePublisherService {
	return &messagePublisherService{
		publisher: publisher,
	}
}

func (m *messagePublisherService) Publish(ctx context.Context, message []byte, image_id string, contentType string) error {
	return m.publisher.Publish(ctx, message, image_id, contentType)
}
