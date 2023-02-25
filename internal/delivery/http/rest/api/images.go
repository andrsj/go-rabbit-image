package api

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	level100 = "100"
	level75  = "75"
	level50  = "50"
	level25  = "25"
)

var levels = [4]string{level100, level75, level50, level25}

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
			gin.H{"error": fmt.Sprintf("Can't read the image: %s", err.Error())},
		)
	}

	imageID := uuid.New().String()

	var img image.Image
	contentType := http.DetectContentType(buf)
	switch contentType {
	case "image/jpeg", "image/jpg":
		img, err = jpeg.Decode(bytes.NewReader(buf))
	case "image/png":
		img, err = png.Decode(bytes.NewReader(buf))
	}
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{"error": fmt.Sprintf("Can't proceed the image type '%s': %s", contentType, err.Error())},
		)
	}

	for _, level := range levels {
		go func(level string) {
			intLevel, _ := strconv.Atoi(level)

			new_img := a.compressService.CompressImage(img, intLevel)

			newBuffer := new(bytes.Buffer)
			switch contentType {
			case "image/jpeg", "image/jpg":
				err = jpeg.Encode(newBuffer, new_img, nil)
			case "image/png":
				err = png.Encode(newBuffer, new_img)
			}

			err = a.imageService.WriteImage(newBuffer.Bytes(), imageID, level)
			if err != nil {
				log.Printf("Error: '%s' can't be compressed by %s\n", imageID, level)
			}
		}(level)
	}

	ctx.JSON(
		http.StatusOK,
		gin.H{
			"message": "Images are being compressed",
			"id":      imageID,
		},
	)

}

func (a *api) GetImage(ctx *gin.Context) {
	id := ctx.Param("id")
	quality := ctx.DefaultQuery("quality", level100)
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
