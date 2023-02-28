package handler

import (
	"github.com/andrsj/go-rabbit-image/internal/delivery/http/rest/api"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	engine *gin.Engine
}

func New() *Handler {
	r := gin.Default()
	r.HandleMethodNotAllowed = true
	return &Handler{
		engine: r,
	}
}

func (h *Handler) GetGinEngine() *gin.Engine {
	return h.engine
}

func (h *Handler) Register(router api.API) {
	h.engine.GET("/ping", router.Ping)
	h.engine.GET("/img/:id", router.GetImage)
	h.engine.POST("/send-image", router.PublishImage)
}
