package app

import (
	"context"
	"fmt"
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
	"github.com/andrsj/go-rabbit-image/internal/infrastructure/worker"
	"github.com/andrsj/go-rabbit-image/internal/infrastructure/worker/compressor"
	"github.com/andrsj/go-rabbit-image/internal/services/image/storage"
	"github.com/andrsj/go-rabbit-image/internal/services/publisher"
)

const (
	path            = "C:/Users/ADerkach/Desktop/go-rabbit-image/"
	image_folder    = "/server_images"
	rabbitURL       = "amqp://guest:guest@localhost:5672/"
	queueName       = "Queue"
	TimeoutDuration = time.Second * 5
)

type App struct {
	srv *http.Server
	job worker.Worker
}

func New() (*App, error) {
	pathToServerFiles := filepath.Join(path, image_folder)
	fileStorage, err := repository.New(pathToServerFiles)
	if err != nil {
		return nil, fmt.Errorf("can't create file storage: %s", err)
	}
	fileService := storage.New(fileStorage)

	rabbitClient, err := client.New(rabbitURL, queueName)
	if err != nil {
		return nil, fmt.Errorf("error connected with RabbitMQ: %s", err)
	}
	publisher := publisher.New(rabbitClient)

	api_router := api.New(fileService, publisher)
	api_handler := handler.New()
	api_handler.Register(api_router)

	compressor := compressor.New()
	jobContext, jobCancelFunc := context.WithCancel(context.Background())
	job := worker.New(
		worker.WithClient(rabbitClient),
		worker.WithFileRepository(fileStorage),
		worker.WithCompressor(compressor),
		worker.WithCancel(jobCancelFunc),
		worker.WithContext(jobContext),
	)

	server := server.New(api_handler)

	return &App{
		srv: server,
		job: job,
	}, nil
}

func (a *App) Start() {
	a.job.Start()

	go func() {
		if err := a.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
}

func (a *App) Stop() error {
	a.job.Stop()
	a.srv.SetKeepAlivesEnabled(false)

	log.Printf("Shutdown server . . . Timeout: %v", TimeoutDuration)
	ctx, cancel := context.WithTimeout(context.Background(), TimeoutDuration)
	defer cancel()

	if err := a.srv.Shutdown(ctx); err != nil {
		return err
	}

	log.Println("Server exiting")

	return nil
}

func (a *App) WaitForShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}
