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
		ctx.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{"error": fmt.Sprintf("Can't open the file: %s", err)},
		)
		return
	}
	defer src.Close()

	buf := make([]byte, file.Size)
	_, err = io.ReadFull(src, buf)
	if err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{"error": fmt.Sprintf("Can't read the image: %s", err)},
		)
		return
	}

	contentType := http.DetectContentType(buf)
	switch contentType {
	case "image/jpeg", "image/png":
	default:
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{"error": fmt.Sprintf("Can't accept the type '%s': please, use jpg/png type", contentType)},
		)
		return
	}

	imageID := uuid.New().String()

	err = a.publisherService.Publish(ctx, buf, imageID, contentType)
	if err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{"error": fmt.Sprintf("Can't publish the image: %s", err)},
		)
		return
	}

	ctx.JSON(
		http.StatusOK,
		gin.H{
			"message": "Images are being compressed",
			"id":      imageID,
		},
	)
}
