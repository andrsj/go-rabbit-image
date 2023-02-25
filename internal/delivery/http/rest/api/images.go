package api

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const level = "100"

func (a *api) PostImage(ctx *gin.Context) {

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
			gin.H{"error": err.Error()},
		)
	}

	err = a.imageService.WriteImage(buf, uuid.New().String(), level)
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()},
		)
	}
}

func (a *api) GetImage(ctx *gin.Context) {
	id := ctx.Param("id")
	quality := ctx.DefaultQuery("quality", level)
	img, err := a.imageService.ReadImage(id, quality)
	if err != nil {
		ctx.JSON(
			http.StatusNotFound,
			gin.H{"error": fmt.Sprintf("Image not found: %s", err)},
		)
		return
	}

	contentType := http.DetectContentType(img)
	ctx.Header("Content-type", contentType)
	ctx.Writer.WriteHeader(http.StatusOK)
	if _, err := ctx.Writer.Write(img); err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{"error": fmt.Sprintf("Failed to send image '%s'", id)},
		)
		return
	}
}
