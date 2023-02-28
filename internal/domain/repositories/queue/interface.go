package queue

import (
	"context"

	"github.com/andrsj/go-rabbit-image/internal/domain/dto"
)

type Publisher interface {
	Publish(ctx context.Context, message []byte, image_id, contentType string) error
}

type Consumer interface {
	ConsumeMessages() (<-chan dto.MessageDTO, <-chan error)
}

type MessageBroker interface {
	Publisher
	Consumer
}
