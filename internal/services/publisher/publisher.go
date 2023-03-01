package publisher

import (
	"context"

	"github.com/andrsj/go-rabbit-image/internal/domain/repositories/queue"
	"github.com/andrsj/go-rabbit-image/pkg/logger"
)

// messagePublisherService: defines a struct for message publishing service.
type messagePublisherService struct {
	publisher queue.Publisher
	logger    logger.Logger
}

var _ queue.Publisher = (*messagePublisherService)(nil)

// New is a constructor function that creates and returns a new instance
func New(publisher queue.Publisher, logger logger.Logger) *messagePublisherService {
	return &messagePublisherService{
		publisher: publisher,                         // queue.Publisher which will be used to publish messages.
		logger:    logger.Named("Publisher service"), // a logger instance
	}
}

// Publish is a method that publishes a message to a queue
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
