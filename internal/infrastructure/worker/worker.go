package worker

import (
	"context"

	"github.com/andrsj/go-rabbit-image/internal/domain/repositories/file"
	"github.com/andrsj/go-rabbit-image/internal/domain/repositories/queue"
	"github.com/andrsj/go-rabbit-image/internal/infrastructure/worker/compressor"
	"github.com/andrsj/go-rabbit-image/pkg/logger"
)

type WorkerParams struct {
	logger         logger.Logger
	client         queue.Consumer
	fileRepository file.FileRepository
	compressor     compressor.Compressor

	cancelFunc context.CancelFunc
	context    context.Context
}

type WorkerOption func(*WorkerParams)

func WithClient(client queue.Consumer) WorkerOption {
	return func(p *WorkerParams) {
		p.client = client
	}
}

func WithCancel(cancel context.CancelFunc) WorkerOption {
	return func(p *WorkerParams) {
		p.cancelFunc = cancel
	}
}

func WithContext(ctx context.Context) WorkerOption {
	return func(p *WorkerParams) {
		p.context = ctx
	}
}

func WithFileRepository(fileRepository file.FileRepository) WorkerOption {
	return func(p *WorkerParams) {
		p.fileRepository = fileRepository
	}
}

func WithCompressor(compressor compressor.Compressor) WorkerOption {
	return func(p *WorkerParams) {
		p.compressor = compressor
	}
}

func WithLogger(logger logger.Logger) WorkerOption {
	return func(p *WorkerParams) {
		p.logger = logger
	}
}

type Worker interface {
	Start()
	Stop()
}

type worker struct {
	client         queue.Consumer
	compressor     compressor.Compressor
	fileRepository file.FileRepository

	cancelFunc context.CancelFunc
	context    context.Context

	logger logger.Logger
}

func New(options ...WorkerOption) *worker {
	params := &WorkerParams{}
	for _, option := range options {
		option(params)
	}
	return &worker{
		client:         params.client,
		fileRepository: params.fileRepository,
		compressor:     params.compressor,
		cancelFunc:     params.cancelFunc,
		context:        params.context,
		logger:         params.logger.Named("background job"),
	}
}
