package api

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/andrsj/go-rabbit-image/pkg/logger"
	"github.com/gin-gonic/gin"
)

const (
	DefaultLevel = "100"
	Level1       = "75"
	Level2       = "50"
	Level3       = "25"
)

type imageParams struct {
	ID            string
	Quality       string
	AllowedValues []string
}

func (a *api) GetImage(ctx *gin.Context) {
	params := imageParams{
		ID:      ctx.Param("id"),
		Quality: ctx.DefaultQuery("quality", DefaultLevel),
	}

	a.logger.Debug("GetImage: Validating image params", logger.M{"params": params})
	err := validateGetImageParams(params)
	if err != nil {
		a.logger.Error("GetImage: Invalid image params", logger.M{
			"error":  err,
			"params": params,
		})
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{"error": fmt.Sprintf("Wrong query parameter: %s", err)},
		)
		return
	}

	a.logger.Debug("GetImage: Reading image from storage", logger.M{
		"image_id": params.ID,
		"quality":  params.Quality},
	)
	img, err := a.imageService.ReadImageFromStorage(params.ID, params.Quality)
	if err != nil {
		a.logger.Error("GetImage: Failed to read image from storage", logger.M{
			"error":    err,
			"image_id": params.ID,
			"quality":  params.Quality,
		})
		ctx.AbortWithStatusJSON(
			http.StatusNotFound,
			gin.H{"error": fmt.Sprintf("Image not found: %s", err)},
		)
		return
	}

	contentType := http.DetectContentType(img)
	ctx.Header("Content-type", contentType)
	ctx.Writer.WriteHeader(http.StatusOK)

	a.logger.Info("GetImage: Sending image to client", logger.M{
		"image_id": params.ID,
		"quality":  params.Quality,
	})
	if _, err := ctx.Writer.Write(img); err != nil {
		a.logger.Error("GetImage: Failed to send image to client", logger.M{
			"error":    err,
			"image_id": params.ID,
		})
		ctx.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{"error": fmt.Sprintf("Failed to send image '%s'", params.ID)},
		)
		return
	}
}

func validateGetImageParams(params imageParams) error {
	if !isValidUUID(params.ID) {
		return fmt.Errorf("invalid format of ID")
	}

	switch params.Quality {
	case DefaultLevel, Level1, Level2, Level3:
		return nil
	default:
		return fmt.Errorf("invalid quality parameter: use 100, 75, 50 or 25")
	}
}

func isValidUUID(s string) bool {
	regex := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)
	return regex.MatchString(s)
}
