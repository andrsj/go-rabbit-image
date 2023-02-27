package worker

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"net/http"
	"strconv"

	"github.com/andrsj/go-rabbit-image/internal/domain/dto"
	"github.com/andrsj/go-rabbit-image/internal/domain/repositories/file"
	"github.com/andrsj/go-rabbit-image/internal/infrastructure/worker/compressor"
)

const (
	level75 = 75
	level50 = 50
	level25 = 25
)

var levels = []int{level75, level50, level25}

type Consumer interface {
	ConsumeMessages() (<-chan dto.MessageDTO, <-chan error)
}

type worker struct {
	client         Consumer
	compressor     compressor.Compressor
	fileRepository file.FileRepository
	cancelFunc     context.CancelFunc
}

func New(
	client Consumer,
	cancel context.CancelFunc,
	fileRepository file.FileRepository,
	compressor compressor.Compressor) *worker {
	return &worker{
		client:         client,
		cancelFunc:     cancel,
		fileRepository: fileRepository,
		compressor:     compressor,
	}
}

func (c *worker) Start(ctx context.Context) {
	messageCh, errorCh := c.client.ConsumeMessages()
	go func() {
		for {
			select {
			case message := <-messageCh:
				// process message
				// log.Println("Received message:", string(message))

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
			case <-ctx.Done():
				log.Println("Shutdown job . . .")
				return
			}
		}
	}()

}

func (c *worker) Stop() {
	c.cancelFunc()
}

func decodeImage(buf []byte) (img image.Image, contentType string, err error) {
	contentType = http.DetectContentType(buf)
	switch contentType {
	case "image/jpeg", "image/jpg":
		img, err = jpeg.Decode(bytes.NewReader(buf))
	case "image/png":
		img, err = png.Decode(bytes.NewReader(buf))
	default:
		err = fmt.Errorf("can't decode []byte to image.Image. Unsupported content-type: %s", contentType)
	}
	return
}

func encodeImage(img image.Image, contentType string) ([]byte, error) {
	var err error
	buffer := new(bytes.Buffer)
	switch contentType {
	case "image/jpeg", "image/jpg":
		err = jpeg.Encode(buffer, img, nil)
	case "image/png":
		err = png.Encode(buffer, img)
	}
	return buffer.Bytes(), err
}
