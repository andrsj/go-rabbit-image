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
	// TODO remove unnecessary urls
	h.engine.GET("/", router.Status)
	h.engine.GET("/l", router.LongTimeStatus)
	h.engine.POST("/img", router.PostImage)
	h.engine.GET("/img/:id", router.GetImage)
	h.engine.GET("/send/:text", router.Publish)
}
