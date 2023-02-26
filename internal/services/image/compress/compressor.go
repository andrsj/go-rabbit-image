package compress

import (
	"image"

	"github.com/nfnt/resize"
)

type Compressor interface {
	CompressImage(img image.Image, percentage int) image.Image
}

type compressorService struct{}

func New() *compressorService {
	return &compressorService{}
}

func (*compressorService) CompressImage(img image.Image, percentage int) image.Image {
	coefficient := float64(percentage) / 100
	newX, newY := float64(img.Bounds().Dx()), float64(img.Bounds().Dy())
	newIMG := resize.Resize(
		uint(newX*coefficient),
		uint(newY*coefficient),
		img,
		resize.Lanczos3,
	)
	return newIMG
}
