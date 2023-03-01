package handler

import (
	"github.com/andrsj/go-rabbit-image/internal/delivery/http/rest/api"
	"github.com/andrsj/go-rabbit-image/pkg/logger"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	engine *gin.Engine
	logger logger.Logger
}

func New(logger logger.Logger) *Handler {
	r := gin.Default()
	r.HandleMethodNotAllowed = true
	return &Handler{
		engine: r,
		logger: logger.Named("Gin engine"),
	}
}

func (h *Handler) GetGinEngine() *gin.Engine {
	return h.engine
}

func (h *Handler) Register(router api.API) {
	h.logger.Info("Registration of controllers", nil)
	h.engine.GET("/ping", router.Ping)
	h.engine.GET("/img/:id", router.GetImage)
	h.engine.POST("/send-image", router.PublishImage)
}
