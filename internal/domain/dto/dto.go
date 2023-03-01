package dto

// MessageDTO represents a message received from a message broker queue
type MessageDTO struct {
	// Slice of bytes that contains the message body.
	Body []byte
	// String that represents the ID of the image associated with the message.
	ImageID string
	// String that represents the content type of the image associated with the message.
	ContentType string
}
