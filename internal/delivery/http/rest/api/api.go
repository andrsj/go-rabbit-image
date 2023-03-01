package api

import (
	"net/http"

	"github.com/andrsj/go-rabbit-image/internal/domain/repositories/queue"
	"github.com/andrsj/go-rabbit-image/internal/services/image/storage"
	"github.com/andrsj/go-rabbit-image/pkg/logger"
	"github.com/gin-gonic/gin"
)

// API interface representation of controllers for Gin engine
type API interface {
	Ping(ctx *gin.Context)
	GetImage(ctx *gin.Context)
	PublishImage(ctx *gin.Context)
}

// API representation of controllers for Gin engine
type api struct {
	imageService     storage.FileStorage
	publisherService queue.Publisher
	logger           logger.Logger
}

var _ API = (*api)(nil)

// New function is a constructor for the api struct.
func New(imageService storage.FileStorage, publisher queue.Publisher, logger logger.Logger) *api {
	return &api{
		imageService:     imageService,
		publisherService: publisher,
		logger:           logger.Named("API"),
	}
}

// Ping method returns a 200 OK status code [use it e.g. health checking]
func (a *api) Ping(ctx *gin.Context) {
	a.logger.Info("Endpoint hit: Ping", nil)
	ctx.String(http.StatusOK, "Ok")
}
