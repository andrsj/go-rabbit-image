package main

import (
	"encoding/json"
	"log"

	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	log.Println("Connecting . . .")

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})

	if err != nil {
		log.Fatal("failed to connect database: ", err)
	}
	log.Println("Connected successful")

	// Migrate the schema
	err = db.AutoMigrate(&Image{})
	if err != nil {
		log.Fatal("failed to Auto Migrate Image")
	}
	err = db.AutoMigrate(&ImageSource{})
	if err != nil {
		log.Fatal("failed to Auto Migrate Image")
	}
	log.Println("Auto Migrations were succeed")

	// Create repositories
	imageRepo := NewImageRepository(db)

	log.Println("Creating Image")

	// Add an image
	imageID := uuid.New()
	imageName := NameType(uuid.New().String())
	// imageName := NameType("Wow")
	image := &Image{
		ID:   imageID,
		Name: imageName,
	}
	if err = imageRepo.Save(image); err != nil {
		log.Fatal("failed to create image: ", err)
	}
	log.Println("Creating Image was successful")
	PrintJsonifyImage(image)

	image, err = imageRepo.GetByID(image.ID)
	if err != nil {
		log.Fatal("failed to get image: ", err)
	}
	PrintJsonifyImage(image)

	sources := []*ImageSource{
		{
			ID:        uuid.New(),
			Ready:     true,
			Level:     100,
			ImageName: imageName,
		},
		{
			ID:        uuid.New(),
			Ready:     false,
			Level:     75,
			ImageName: imageName,
		},
		{
			ID:        uuid.New(),
			Ready:     true,
			Level:     50,
			ImageName: imageName,
		},
		{
			ID:        uuid.New(),
			Ready:     true,
			Level:     25,
			ImageName: imageName,
		},
	}

	for _, source := range sources {
		if err = imageRepo.SaveSource(source); err != nil {
			log.Fatal("failed to create source:", err)
		}
	}

	image, _ = imageRepo.GetByID(image.ID)
	PrintJsonifyImage(image)

	log.Println("Getting sources")
	imageSources, _ := imageRepo.GetSourcesByName(image.Name)
	for _, source := range imageSources {
		PrintJsonifyImage(source)
	}

	log.Println("Getting source 75%")
	imageSource, _ := imageRepo.GetSourceByNameAndQuality(image.Name, 75)
	PrintJsonifyImage(imageSource)

}

func PrintJsonifyImage(a any) {
	log.Println("Data:")
	jsonImage, _ := json.MarshalIndent(a, "", "    ")
	log.Println(string(jsonImage))
}
