package worker

import (
	"strconv"

	"github.com/andrsj/go-rabbit-image/pkg/logger"
)

const (
	level75 = 75
	level50 = 50
	level25 = 25
)

var levels = []int{level75, level50, level25}

// Start() method of the worker struct
func (c *worker) Start() {
	// Start consuming messages from the Consumer
	messageCh, errorCh := c.client.MustConsumeMessages()
	c.logger.Info("Consumer has started", nil)
	go func() {
		for {
			select {
			case message := <-messageCh:

				// Decode the image from the message body
				img, contentType, err := decodeImage(message.Body)
				if err != nil {
					c.logger.Error("Decoding image", logger.M{"error": err})
					c.logger.Warn("Skipping image", logger.M{"image_id": message.ImageID})
					continue
				}

				// Create image with 100% quality
				go func() {
					err := c.fileRepository.CreateImage(message.Body, message.ImageID, "100")
					if err != nil {
						c.logger.Error("Creating image", logger.M{"error": err})
						c.logger.Warn("Skipping image", logger.M{"image_id": message.ImageID})
						return
					}
				}()

				// Compress the image and create images with different levels of quality
				for _, level := range levels {

					go func(level int) {
						// Compress the image to a specific quality level
						new_img := c.compressor.CompressImage(img, level)

						// Encode the compressed image
						bufferImage, err := encodeImage(new_img, contentType)
						if err != nil {
							c.logger.Error("Encoding image", logger.M{"error": err})
							c.logger.Warn("Skipping image", logger.M{"image_id": message.ImageID})
							return
						}

						// Create image with the given quality level
						err = c.fileRepository.CreateImage(bufferImage, message.ImageID, strconv.Itoa(level))
						if err != nil {
							c.logger.Error("Creating image", logger.M{"error": err})
							c.logger.Warn("Skipping image", logger.M{"image_id": message.ImageID})
							return
						}

					}(level)

				}

			case err := <-errorCh:
				if err != nil {
					// Log the error and exit the application
					c.logger.Fatal("Consuming messages", logger.M{"error": err})
				}

			// Graceful shutdown?
			case <-c.context.Done():
				// Log the shutdown and return from the method
				c.logger.Info("Shutdown job . . .", nil)
				return
			}
		}
	}()

}

func (c *worker) Stop() {
	c.cancelFunc()
}
