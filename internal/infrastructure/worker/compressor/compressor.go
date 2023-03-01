package compressor

import (
	"image"

	"github.com/andrsj/go-rabbit-image/pkg/logger"
	"github.com/nfnt/resize"
)

type Compressor interface {
	CompressImage(img image.Image, percentage int) image.Image
}

type compressorService struct {
	logger logger.Logger
}

func New(logger logger.Logger) *compressorService {
	return &compressorService{
		logger: logger.Named("Comressor service"),
	}
}

func (c *compressorService) CompressImage(img image.Image, percentage int) image.Image {
	c.logger.Info("Compressing image", logger.M{"%": percentage})
	coefficient := float64(percentage) / 100
	newX, newY := float64(img.Bounds().Dx()), float64(img.Bounds().Dy())
	newIMG := resize.Resize(
		uint(newX*coefficient),
		uint(newY*coefficient),
		img,
		resize.Lanczos3,
	)
	c.logger.Info("Image compressed successfully", logger.M{"%": percentage})
	return newIMG
}
