package handler

import (
	"github.com/andrsj/go-rabbit-image/internal/delivery/http/rest/api"
	"github.com/andrsj/go-rabbit-image/pkg/logger"
	"github.com/gin-gonic/gin"
)

// Handler is a struct that holds a gin.Engine instance and a logger instance.
type Handler struct {
	engine *gin.Engine
	logger logger.Logger
}

// New is a constructor function that initializes a gin.Engine instance and returns a Handler instance.
func New(logger logger.Logger) *Handler {
	r := gin.Default()
	r.HandleMethodNotAllowed = true

	return &Handler{
		engine: r,
		logger: logger.Named("Gin engine"),
	}
}

// GetGinEngine is a method that returns the gin.Engine instance from a Handler instance.
func (h *Handler) GetGinEngine() *gin.Engine {
	return h.engine
}

// Register is a method that registers the API routes defined in api.API to the gin.Engine instance in Handler.
func (h *Handler) Register(router api.API) {
	h.logger.Info("Registration of controllers", nil)
	h.engine.GET("/ping", router.Ping)
	h.engine.GET("/img/:id", router.GetImage)
	h.engine.POST("/send-image", router.PublishImage)
}
