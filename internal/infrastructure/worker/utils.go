package worker

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"
)

/*
The decodeImage function decodes a byte slice to an image.Image
and returns the decoded image, the detected
content type of the input byte slice, and an error (if any).

It takes a byte slice as input and first detects
the content type of the image using the http.DetectContentType function.

If the content type is supported, the function decodes
the byte slice to an image using either
the jpeg.Decode or png.Decode function depending on the content type.

If the content type is not supported,
the function returns an error with a message that the content type is unsupported.
*/
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

/*
The encodeImage function encodes an image.Image to a byte slice
and returns the byte slice and an error (if any).

It takes an image and the content type of the image as input.
The function creates a new buffer to hold the encoded image data,
and uses the jpeg.Encode or png.Encode function
depending on the content type to encode the image to the buffer.

If the encoding process succeeds, the function returns
the bytes in the buffer and nil error, otherwise it returns an error.
*/
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
