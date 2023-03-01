package publisher

import (
	"context"

	"github.com/andrsj/go-rabbit-image/internal/domain/repositories/queue"
	"github.com/andrsj/go-rabbit-image/pkg/logger"
)

type messagePublisherService struct {
	publisher queue.Publisher
	logger    logger.Logger
}

var _ queue.Publisher = (*messagePublisherService)(nil)

func New(publisher queue.Publisher, logger logger.Logger) *messagePublisherService {
	return &messagePublisherService{
		publisher: publisher,
		logger:    logger.Named("Publisher service"),
	}
}

func (m *messagePublisherService) Publish(ctx context.Context, message []byte, image_id string, contentType string) error {
	m.logger.Debug("Publishing message", logger.M{
		"image_id":     image_id,
		"content_type": contentType,
	})

	err := m.publisher.Publish(ctx, message, image_id, contentType)
	if err != nil {
		m.logger.Error("Failed to publish message", logger.M{
			"error":        err,
			"image_id":     image_id,
			"content_type": contentType,
		})
		return err
	}

	m.logger.Info("Message published", logger.M{
		"image_id":     image_id,
		"content_type": contentType,
	})
	return nil
}
