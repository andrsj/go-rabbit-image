package worker

import (
	"log"
	"strconv"
)

const (
	level75 = 75
	level50 = 50
	level25 = 25
)

var levels = []int{level75, level50, level25}

func (c *worker) Start() {
	messageCh, errorCh := c.client.MustConsumeMessages()
	go func() {
		for {
			select {
			case message := <-messageCh:

				img, contentType, err := decodeImage(message.Body)
				if err != nil {
					// TODO change logger
					log.Printf("Error: %s\n", err)
					log.Printf("Warning: skipping id:'%s'", message.ImageID)
					continue
				}

				go func() {
					err := c.fileRepository.CreateImage(message.Body, message.ImageID, "100")
					if err != nil {
						// TODO report errors in DB
						log.Printf("Error: %s\n", err)
						log.Printf("Warning: skipping id:'%s'", message.ImageID)
						return
					}
				}()
				for _, level := range levels {

					go func(level int) {

						new_img := c.compressor.CompressImage(img, level)
						bufferImage, err := encodeImage(new_img, contentType)
						if err != nil {
							// TODO report errors in DB
							log.Printf("Error: %s\n", err)
							log.Printf("Warning: skipping id:'%s'", message.ImageID)
							return
						}
						err = c.fileRepository.CreateImage(bufferImage, message.ImageID, strconv.Itoa(level))
						if err != nil {
							// TODO report errors in DB
							log.Printf("Error: %s\n", err)
							log.Printf("Warning: skipping id:'%s'", message.ImageID)
							return
						}

					}(level)

				}

			case err := <-errorCh:
				if err != nil {
					log.Fatalln("Error while consuming messages:", err)
				}
			case <-c.context.Done():
				log.Println("Shutdown job . . .")
				return
			}
		}
	}()

}

func (c *worker) Stop() {
	c.cancelFunc()
}
