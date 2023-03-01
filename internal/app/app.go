package app

import (
	"context"
	"fmt"
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
	"github.com/andrsj/go-rabbit-image/pkg/logger"
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
	log logger.Logger
}

func New(log logger.Logger) (*App, error) {
	log = log.Named("app")

	rabbitClient, err := client.New(rabbitURL, queueName, log)
	if err != nil {
		err = fmt.Errorf("error connected with RabbitMQ: %s", err)
		log.Error("Connection error to Message Broker", logger.M{
			"error": err,
		})
		return nil, err
	}
	publisher := publisher.New(rabbitClient, log)

	pathToServerFiles := filepath.Join(path, image_folder)
	fileStorage, err := repository.New(pathToServerFiles, log)
	if err != nil {
		log.Error("Can't create file storage", logger.M{
			"error": err,
		})
		return nil, fmt.Errorf("can't create file storage: %s", err)
	}
	fileService := storage.New(fileStorage, log)

	api_router := api.New(fileService, publisher, log)
	api_handler := handler.New(log)
	api_handler.Register(api_router)

	compressor := compressor.New(log)
	jobContext, jobCancelFunc := context.WithCancel(context.Background())
	job := worker.New(
		worker.WithClient(rabbitClient),
		worker.WithFileRepository(fileStorage),
		worker.WithCompressor(compressor),
		worker.WithCancel(jobCancelFunc),
		worker.WithContext(jobContext),
		worker.WithLogger(log),
	)

	server := server.New(api_handler)

	return &App{
		srv: server,
		job: job,
		log: log,
	}, nil
}

func (a *App) Start() {

	a.log.Info("Starting background job", nil)
	a.job.Start()

	a.log.Info("Starting server", logger.M{
		"address": a.srv.Addr,
	})
	go func() {
		if err := a.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.log.Fatal("Server starting error.", logger.M{
				"error": err,
			})
		}
	}()
}

func (a *App) Stop() error {
	a.log.Info("Stopping background job", nil)
	a.job.Stop()

	a.log.Info("Closing keep-alive connections", nil)
	a.srv.SetKeepAlivesEnabled(false)

	a.log.Info("Shutdown server . . . Timeout", logger.M{
		"timeout": TimeoutDuration,
	})
	ctx, cancel := context.WithTimeout(context.Background(), TimeoutDuration)
	defer cancel()

	if err := a.srv.Shutdown(ctx); err != nil {
		a.log.Error("Error shutdown", logger.M{
			"error": err,
		})
		return err
	}

	a.log.Info("Server exiting", nil)

	return nil
}

func (a *App) WaitForShutdown() {
	a.log.Info("Waiting for shutdown", nil)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	a.log.Info("Got signal for terminating server", logger.M{
		"signal": sig,
	})
}
