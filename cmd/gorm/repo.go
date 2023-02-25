package main

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type NameType string

type Image struct {
	ID           uuid.UUID     `gorm:"primaryKey"`
	Name         NameType      `gorm:"not null"`
	ImageSources []ImageSource `gorm:"foreignKey:ImageName;references:Name"`
}

type ImageSource struct {
	ID        uuid.UUID `gorm:"primaryKey"`
	Ready     bool      `gorm:"not null"`
	Level     int       `gorm:"not null"`
	ImageName NameType
}

func (i *ImageSource) GetBasePath() string {
	return fmt.Sprintf("%s\\%v", i.ImageName, i.Level)
}

type ImageRepository interface {
	GetByID(id uuid.UUID) (*Image, error)
	GetByName(name NameType) (*Image, error)
	GetSourceByNameAndQuality(name NameType, quality int) (*ImageSource, error)
	GetSourcesByName(name NameType) ([]ImageSource, error)
	SaveSource(imageSource *ImageSource) error
	Save(image *Image) error
}

type imageRepository struct {
	db *gorm.DB
}

func NewImageRepository(db *gorm.DB) ImageRepository {
	return &imageRepository{db: db}
}

func (ir *imageRepository) GetByID(id uuid.UUID) (*Image, error) {
	var image Image
	result := ir.db.Preload("ImageSources").First(&image, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &image, nil
}

func (ir *imageRepository) GetByName(name NameType) (*Image, error) {
	var image Image
	result := ir.db.Preload("ImageSources").First(&image, "name = ?", name)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &image, nil
}

func (ir *imageRepository) GetSourceByNameAndQuality(name NameType, quality int) (*ImageSource, error) {
	var imageSource ImageSource
	result := ir.db.Where("image_name = ? AND level = ?", name, quality).First(&imageSource)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &imageSource, nil
}

func (ir *imageRepository) GetSourcesByName(name NameType) ([]ImageSource, error) {
	var imageSources []ImageSource
	result := ir.db.Where("image_name = ?", name).Find(&imageSources)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return imageSources, nil
}

func (ir *imageRepository) SaveSource(imageSource *ImageSource) error {
	return ir.db.Create(imageSource).Error
}

func (ir *imageRepository) Save(image *Image) error {
	return ir.db.Create(image).Error
}
