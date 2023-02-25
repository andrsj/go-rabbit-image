package api

import (
	"github.com/andrsj/go-rabbit-image/internal/services/image/compress"
	"github.com/andrsj/go-rabbit-image/internal/services/image/storage"
	"github.com/gin-gonic/gin"
)

type APIInterface interface {
	Status(ctx *gin.Context)
	LongTimeStatus(ctx *gin.Context)
	PostImage(ctx *gin.Context)
	GetImage(ctx *gin.Context)
}

type api struct {
	imageService    storage.FileStorageInterface
	compressService compress.CompressorInterface
}

var _ APIInterface = (*api)(nil)

func New(imageService storage.FileStorageInterface, compressService compress.CompressorInterface) *api {
	return &api{
		imageService:    imageService,
		compressService: compressService,
	}
}
