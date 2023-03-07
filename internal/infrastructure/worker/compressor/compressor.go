package compressor

import (
	"image"

	"github.com/andrsj/go-rabbit-image/pkg/logger"
	"github.com/nfnt/resize"
)

const originalSizePercentage = 100

// Compressor is an interface that defines the CompressImage method.
type Compressor interface {
	CompressImage(img image.Image, percentage int) image.Image
}

// compressorService is a struct that holds a logger and implements the Compressor interface.
type compressorService struct {
	logger logger.Logger
}

// New is a constructor of the compressorService struct.
func New(logger logger.Logger) *compressorService {
	return &compressorService{
		logger: logger.Named("Compressor service"),
	}
}

// CompressImage is a method for compressing images by github.com/nfnt/resize package.
func (c *compressorService) CompressImage(img image.Image, percentage int) image.Image {
	c.logger.Info("Compressing image", logger.M{"%": percentage})

	coefficient := float64(percentage) / originalSizePercentage
	newX, newY := float64(img.Bounds().Dx()), float64(img.Bounds().Dy())

	// It resizes the image using the given percentage and the Lanczos3 interpolation method
	newIMG := resize.Resize(
		uint(newX*coefficient),
		uint(newY*coefficient),
		img,
		resize.Lanczos3,
	)

	c.logger.Info("Image compressed successfully", logger.M{"%": percentage})

	return newIMG
}
