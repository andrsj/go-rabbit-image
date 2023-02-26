package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/andrsj/go-rabbit-image/internal/delivery/http/handler"
	"github.com/andrsj/go-rabbit-image/internal/delivery/http/rest/api"
	"github.com/andrsj/go-rabbit-image/internal/delivery/http/server"
	"github.com/andrsj/go-rabbit-image/internal/delivery/rabbitmq/client"
	"github.com/andrsj/go-rabbit-image/internal/infrastructure/file/repository"
	"github.com/andrsj/go-rabbit-image/internal/services/image/compress"
	"github.com/andrsj/go-rabbit-image/internal/services/image/storage"
)

const (
	path      = "C:/Users/ADerkach/Desktop/Image"
	rabbitURL = "amqp://guest:guest@localhost:5672/"
)

func main() {
	pathToServerFiles := filepath.Join(path, "/server_images")
	fileStorage, err := repository.New(pathToServerFiles)
	if err != nil {
		log.Fatalf("Can't create file storage: %s", err)
	}
	fileService := storage.New(fileStorage)
	compressService := compress.New()

	publisher, err := client.New(rabbitURL, "QUEUE")
	if err != nil {
		panic(err)
	}

	api_router := api.New(fileService, compressService, publisher)
	api_handler := handler.New()
	api_handler.Register(api_router)

	server := server.New(api_handler)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	server.SetKeepAlivesEnabled(false)
	duration := time.Second * 5
	log.Printf("Shutdown server . . . Timeout: %v", duration)
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown: %s", err)
	}

	log.Println("Server exiting")
}
