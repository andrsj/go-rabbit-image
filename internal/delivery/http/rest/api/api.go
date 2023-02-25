package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type APIInterface interface {
	Status(ctx *gin.Context)
	LongTimeStatus(ctx *gin.Context)
}

type API struct{}

func New() *API {
	return &API{}
}

func (*API) Status(ctx *gin.Context) {
	ctx.JSON(
		http.StatusOK,
		gin.H{
			"message": "Hello!",
		},
	)
}

func (*API) LongTimeStatus(ctx *gin.Context) {
	duration := time.Second * 6
	time.Sleep(duration)
	ctx.JSON(
		http.StatusOK,
		gin.H{
			"message": "Hello!",
			"work":    duration,
		},
	)
}
