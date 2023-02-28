package api

import (
	"fmt"
	"net/http"
	"regexp"

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

	err := validateGetImageParams(params)
	if err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{"error": fmt.Sprintf("Wrong query parameter: %s", err)},
		)
		return
	}

	img, err := a.imageService.ReadImageFromStorage(params.ID, params.Quality)
	if err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusNotFound,
			gin.H{"error": fmt.Sprintf("Image not found: %s", err)},
		)
		return
	}

	contentType := http.DetectContentType(img)
	ctx.Header("Content-type", contentType)
	ctx.Writer.WriteHeader(http.StatusOK)
	if _, err := ctx.Writer.Write(img); err != nil {
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
