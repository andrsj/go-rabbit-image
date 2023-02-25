package main

import (
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"path/filepath"

	"github.com/nfnt/resize"
)

func main() {
	if len(os.Args) != 3 {
		log.Fatal("Please provide a path as an argument")
	}
	old_path := os.Args[1]
	new_path := os.Args[2]
	log.Printf("Old: %s \tNew: %s", old_path, new_path)

	file, _ := os.Open(old_path)

	fileInfo, _ := file.Stat()

	filename := fileInfo.Name()

	var img image.Image

	switch filepath.Ext(filename) {
	case ".jpg", ".jpeg":
		log.Println("Got .jpg || .jpeg image")
		img, _ = jpeg.Decode(file)
	case ".png":
		log.Println("Got .png")
		img, _ = png.Decode(file)
	}

	coefficient := 0.25
	resized := resize.Resize(
		uint(float64(img.Bounds().Dx())*coefficient),
		uint(float64(img.Bounds().Dy())*coefficient),
		img,
		resize.Lanczos3,
	)

	out, _ := os.Create(new_path)
	defer out.Close()

	if err := jpeg.Encode(out, resized, nil); err != nil {
		panic(err)
	}

}
