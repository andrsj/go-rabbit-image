package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (a *api) Publish(ctx *gin.Context) {
	message := ctx.Param("text")
	err := a.publisherService.Publish(ctx, []byte(message), "NO ID", "plain/text")
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()},
		)
	}

	ctx.JSON(
		http.StatusOK,
		gin.H{"message": fmt.Sprintf("Text '%s' sent", message)},
	)
}
