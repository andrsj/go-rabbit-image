package worker

import (
	"context"

	"github.com/andrsj/go-rabbit-image/internal/domain/repositories/file"
	"github.com/andrsj/go-rabbit-image/internal/domain/repositories/queue"
	"github.com/andrsj/go-rabbit-image/internal/infrastructure/worker/compressor"
	"github.com/andrsj/go-rabbit-image/pkg/logger"
)

type Params struct {
	logger         logger.Logger
	client         queue.Consumer
	fileRepository file.Repository
	compressor     compressor.Compressor

	cancelFunc context.CancelFunc
	context    context.Context
}

/*
Actually the FunctionParameters pattern is no needed here, but...
I just want to show that I know it and don't want to delete the implementation of this pattern.
*/
type Option func(*Params)

func WithClient(client queue.Consumer) Option {
	return func(p *Params) {
		p.client = client
	}
}

func WithCancel(cancel context.CancelFunc) Option {
	return func(p *Params) {
		p.cancelFunc = cancel
	}
}

func WithContext(ctx context.Context) Option {
	return func(p *Params) {
		p.context = ctx
	}
}

func WithFileRepository(fileRepository file.Repository) Option {
	return func(p *Params) {
		p.fileRepository = fileRepository
	}
}

func WithCompressor(compressor compressor.Compressor) Option {
	return func(p *Params) {
		p.compressor = compressor
	}
}

func WithLogger(logger logger.Logger) Option {
	return func(p *Params) {
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
	fileRepository file.Repository

	cancelFunc context.CancelFunc
	context    context.Context

	logger logger.Logger
}

func New(options ...Option) *worker {
	params := &Params{nil, nil, nil, nil, nil, nil}

	// There is a problem that I DON'T CHECK
	// if some REQUIRED parameter is not provided
	//
	// I can validate the <nil> value in functions,
	// but how to check if all With<Parameter> functions
	// was called?
	for _, option := range options {
		option(params)
	}

	return &worker{
		client:         params.client,
		fileRepository: params.fileRepository,
		compressor:     params.compressor,
		cancelFunc:     params.cancelFunc,
		context:        params.context,
		// Question: is it good to pass the name here?
		// Because it's a constructor of the worker instance
		// ..., but is it idiomatic way of GO?
		logger: params.logger.Named("background job"),
	}
}
