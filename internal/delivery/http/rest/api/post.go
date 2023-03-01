package api

import (
	"fmt"
	"io"
	"net/http"

	"github.com/andrsj/go-rabbit-image/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (a *api) PublishImage(ctx *gin.Context) {
	file, err := ctx.FormFile("image")
	if err != nil {
		a.logger.Error("Can't get image from form data", logger.M{"error": err})
		return
	}

	src, err := file.Open()
	if err != nil {
		a.logger.Error("Can't open the image file", logger.M{"error": err})
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
		a.logger.Error("Can't read the image", logger.M{"error": err})
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
		a.logger.Error("Can't accept the type of image", logger.M{"content type": contentType})
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{"error": fmt.Sprintf("Can't accept the type '%s': please, use jpg/png type", contentType)},
		)
		return
	}

	imageID := uuid.New().String()

	err = a.publisherService.Publish(ctx, buf, imageID, contentType)
	if err != nil {
		a.logger.Error("Can't publish the image", logger.M{"error": err})
		ctx.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{"error": fmt.Sprintf("Can't publish the image: %s", err)},
		)
		return
	}

	a.logger.Info("Successfully published the image", logger.M{
		"id":           imageID,
		"content type": contentType,
	})
	ctx.JSON(
		http.StatusOK,
		gin.H{
			"message": "Images are being compressed",
			"id":      imageID,
		},
	)
}
