package queue

import (
	"context"

	"github.com/andrsj/go-rabbit-image/internal/domain/dto"
)

type Publisher interface {
	Publish(ctx context.Context, message []byte, imageID, contentType string) error
}

type Consumer interface {
	MustConsumeMessages() (<-chan dto.MessageDTO, <-chan error)
}

type MessageBroker interface {
	Publisher
	Consumer
}
