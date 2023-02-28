package api

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (a *api) PublishImage(ctx *gin.Context) {
	file, err := ctx.FormFile("image")
	if err != nil {
		return
	}

	src, err := file.Open()
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()},
		)
	}
	defer src.Close()

	buf := make([]byte, file.Size)
	_, err = io.ReadFull(src, buf)
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{"error": fmt.Sprintf("Can't read the image: %s", err.Error())},
		)
	}

	imageID := uuid.New().String()

	err = a.publisherService.Publish(ctx, buf, imageID, http.DetectContentType(buf))
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()},
		)
	}

	ctx.JSON(
		http.StatusOK,
		gin.H{
			"message": "Images are being compressed",
			"id":      imageID,
		},
	)
}
