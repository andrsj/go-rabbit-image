package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (*api) Status(ctx *gin.Context) {
	ctx.JSON(
		http.StatusOK,
		gin.H{
			"message": "Hello!",
		},
	)
}

func (*api) LongTimeStatus(ctx *gin.Context) {
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
