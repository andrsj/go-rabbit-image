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

func (c *worker) Start() {
	messageCh, errorCh := c.client.MustConsumeMessages()
	c.logger.Info("Consumer has started", nil)
	go func() {
		for {
			select {
			case message := <-messageCh:

				img, contentType, err := decodeImage(message.Body)
				if err != nil {
					c.logger.Error("Decoding image", logger.M{"error": err})
					c.logger.Warn("Skipping image", logger.M{"image_id": message.ImageID})
					continue
				}

				go func() {
					err := c.fileRepository.CreateImage(message.Body, message.ImageID, "100")
					if err != nil {
						c.logger.Error("Creating image", logger.M{"error": err})
						c.logger.Warn("Skipping image", logger.M{"image_id": message.ImageID})
						return
					}
				}()
				for _, level := range levels {

					go func(level int) {

						new_img := c.compressor.CompressImage(img, level)
						bufferImage, err := encodeImage(new_img, contentType)
						if err != nil {
							c.logger.Error("Encoding image", logger.M{"error": err})
							c.logger.Warn("Skipping image", logger.M{"image_id": message.ImageID})
							return
						}
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
					c.logger.Fatal("Consuming messages", logger.M{"error": err})
				}
			case <-c.context.Done():
				c.logger.Info("Shutdown job . . .", nil)
				return
			}
		}
	}()

}

func (c *worker) Stop() {
	c.cancelFunc()
}
