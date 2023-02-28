package worker

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"
)

func decodeImage(buf []byte) (img image.Image, contentType string, err error) {
	contentType = http.DetectContentType(buf)
	switch contentType {
	case "image/jpeg", "image/jpg":
		img, err = jpeg.Decode(bytes.NewReader(buf))
	case "image/png":
		img, err = png.Decode(bytes.NewReader(buf))
	default:
		err = fmt.Errorf("can't decode []byte to image.Image. Unsupported content-type: %s", contentType)
	}
	return
}

func encodeImage(img image.Image, contentType string) ([]byte, error) {
	var err error
	buffer := new(bytes.Buffer)
	switch contentType {
	case "image/jpeg", "image/jpg":
		err = jpeg.Encode(buffer, img, nil)
	case "image/png":
		err = png.Encode(buffer, img)
	}
	return buffer.Bytes(), err
}
