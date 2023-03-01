package storage

import (
	"github.com/andrsj/go-rabbit-image/internal/domain/repositories/file"
	"github.com/andrsj/go-rabbit-image/pkg/logger"
)

// FileStorage interface represents a service to write and read image data to a file storage.
type FileStorage interface {
	WriteImageToStorage(image []byte, id string, level string) error
	ReadImageFromStorage(id string, level string) ([]byte, error)
}

// fileStorageService represents a service that writes and reads image data to/from a file storage.
type fileStorageService struct {
	fileStorage file.FileRepository
	logger      logger.Logger
}

var _ FileStorage = (*fileStorageService)(nil)

// New creates a new instance of fileStorageService.
func New(storage file.FileRepository, logger logger.Logger) *fileStorageService {
	return &fileStorageService{
		fileStorage: storage,
		logger:      logger,
	}
}

// WriteImageToStorage writes image data to a file storage.
func (f *fileStorageService) WriteImageToStorage(image []byte, name string, level string) error {
	if err := f.fileStorage.CreateImage(image, name, level); err != nil {
		f.logger.Error("Error writing image to storage", logger.M{
			"error": err,
			"name":  name,
			"level": level,
		})
		return err
	}

	f.logger.Info("Image written to storage", logger.M{
		"name":  name,
		"level": level,
	})
	return nil
}

// ReadImageFromStorage reads image data from a file storage.
func (f fileStorageService) ReadImageFromStorage(name string, level string) ([]byte, error) {
	data, err := f.fileStorage.GetImage(name, level)
	if err != nil {
		f.logger.Error("Error reading image from storage", logger.M{
			"error": err,
			"name":  name,
			"level": level,
		})
		return nil, err
	}

	f.logger.Info("Image read from storage", logger.M{
		"name":  name,
		"level": level,
	})
	return data, nil
}
