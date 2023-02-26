package queue

import "context"

type Publisher interface {
	Publish(ctx context.Context, message []byte, image_id, contentType string) error
}

type Consumer interface {
	ConsumeMessages() (<-chan []byte, <-chan error)
}

type MessageBroker interface {
	Publisher
	Consumer
}
