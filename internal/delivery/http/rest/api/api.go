package api

import (
	"github.com/andrsj/go-rabbit-image/internal/domain/repositories/queue"
	"github.com/andrsj/go-rabbit-image/internal/services/image/compress"
	"github.com/andrsj/go-rabbit-image/internal/services/image/storage"
	"github.com/gin-gonic/gin"
)

type API interface {
	// TODO remove unnecessary URLs

	Status(ctx *gin.Context)
	LongTimeStatus(ctx *gin.Context)
	PostImage(ctx *gin.Context)
	GetImage(ctx *gin.Context)

	// Temporary
	Publish(ctx *gin.Context)
	PublishImage(ctx *gin.Context)
}

type api struct {
	imageService     storage.FileStorage
	compressService  compress.Compressor
	publisherService queue.Publisher
}

var _ API = (*api)(nil)

func New(imageService storage.FileStorage, compressService compress.Compressor, publisher queue.Publisher) *api {
	return &api{
		imageService:     imageService,
		compressService:  compressService,
		publisherService: publisher,
	}
}
