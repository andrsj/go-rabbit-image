package worker

import (
	"context"
	"log"
)

type Consumer interface {
	ConsumeMessages() (<-chan []byte, <-chan error)
}

type worker struct {
	client     Consumer
	cancelFunc context.CancelFunc
}

func New(client Consumer, cancel context.CancelFunc) *worker {
	return &worker{
		client:     client,
		cancelFunc: cancel,
	}
}

func (c *worker) Start(ctx context.Context) {
	messageCh, errorCh := c.client.ConsumeMessages()
	go func() {
		for {
			select {
			case message := <-messageCh:
				// process message
				log.Println("Received message:", string(message))
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
