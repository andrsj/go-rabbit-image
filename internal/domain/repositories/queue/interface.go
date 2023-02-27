package queue

import "context"

type Publisher interface {
	Publish(ctx context.Context, message []byte, image_id, contentType string) error
}
