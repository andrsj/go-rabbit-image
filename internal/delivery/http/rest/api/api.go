package api

import (
	"net/http"

	"github.com/andrsj/go-rabbit-image/internal/domain/repositories/queue"
	"github.com/andrsj/go-rabbit-image/internal/services/image/storage"
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
}

var _ API = (*api)(nil)

func New(imageService storage.FileStorage, publisher queue.Publisher) *api {
	return &api{
		imageService:     imageService,
		publisherService: publisher,
	}
}

func (*api) Ping(ctx *gin.Context) {
	ctx.String(http.StatusOK, "Ok")
}
