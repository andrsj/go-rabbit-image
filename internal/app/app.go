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

// New creates a new App object and returns a pointer to it.
func New(log logger.Logger) (*App, error) {
	// It takes a logger.Logger object as input and returns an error
	// if it fails to create any of the necessary components.
	log = log.Named("app")

	// It initializes a rabbitMQ client with the given URL
	// and queue name and logs errors if any occur.
	rabbitClient, err := client.New(rabbitURL, queueName, log)
	if err != nil {
		err = fmt.Errorf("error connected with RabbitMQ: %s", err)
		log.Error("Connection error to Message Broker", logger.M{
			"error": err,
		})
		return nil, err
	}
	publisher := publisher.New(rabbitClient, log)

	// It creates a file storage repository and
	// an associated file service, logging any errors that occur.
	pathToServerFiles := filepath.Join(path, image_folder)
	fileStorage, err := repository.New(pathToServerFiles, log)
	if err != nil {
		log.Error("Can't create file storage", logger.M{
			"error": err,
		})
		return nil, fmt.Errorf("can't create file storage: %s", err)
	}
	fileService := storage.New(fileStorage, log)

	// It creates an API router and handler with the file service
	// and publisher, and registers the router to the handler.
	api_router := api.New(fileService, publisher, log)
	api_handler := handler.New(log)
	api_handler.Register(api_router)

	// It creates a compressor with the logger.
	compressor := compressor.New(log)

	// It creates a job, job's context, cancel function for the worker using the logger.
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

// Start method is responsible for starting the server and background job.
func (a *App) Start() {
	// Start the background job.
	a.log.Info("Starting background job", nil)
	a.job.Start()

	// Launch the server in a separate goroutine.
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
	// Stop background job
	a.log.Info("Stopping background job", nil)
	a.job.Stop()

	// Close keep-alive connections
	a.log.Info("Closing keep-alive connections", nil)
	a.srv.SetKeepAlivesEnabled(false)

	// Shutdown server with a timeout
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

/*
WaitForShutdown function sets up a channel to listen for signals
indicating that the server should be shut down.

It logs when it begins waiting for shutdown
and when it receives a signal to terminate the server.
*/
func (a *App) WaitForShutdown() {
	a.log.Info("Waiting for shutdown", nil)

	/*
		Specifically, it creates a channel to listen for SIGINT and SIGTERM signals,
		which are sent to the process when the user requests to terminate the program
		(e.g. by pressing Ctrl+C in the terminal).

		Once the signal is received, the function logs the signal
		that was received and the server will begin the shutdown process.
	*/
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	a.log.Info("Got signal for terminating server", logger.M{
		"signal": sig,
	})
}
