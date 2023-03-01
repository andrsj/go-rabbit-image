package api

import (
	"net/http"

	"github.com/andrsj/go-rabbit-image/internal/domain/repositories/queue"
	"github.com/andrsj/go-rabbit-image/internal/services/image/storage"
	"github.com/andrsj/go-rabbit-image/pkg/logger"
	"github.com/gin-gonic/gin"
)

type API interface {
	Ping(ctx *gin.Context)
	GetImage(ctx *gin.Context)
	PublishImage(ctx *gin.Context)
}

type api struct {
	imageService     storage.FileStorage
	publisherService queue.Publisher
	logger           logger.Logger
}

var _ API = (*api)(nil)

func New(imageService storage.FileStorage, publisher queue.Publisher, logger logger.Logger) *api {
	return &api{
		imageService:     imageService,
		publisherService: publisher,
		logger:           logger.Named("API"),
	}
}

func (a *api) Ping(ctx *gin.Context) {
	a.logger.Info("Endpoint hit: Ping", nil)
	ctx.String(http.StatusOK, "Ok")
}
