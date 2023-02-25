package handler

import (
	"github.com/andrsj/go-rabbit-image/internal/delivery/http/rest/api"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	Engine *gin.Engine
}

func New() *Handler {
	r := gin.Default()
	return &Handler{
		Engine: r,
	}
}

func (s *Handler) Register(router api.APIInterface) {
	s.Engine.GET("/", router.Status)
	s.Engine.GET("/l", router.LongTimeStatus)
}
