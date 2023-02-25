package storage

import "github.com/andrsj/go-rabbit-image/internal/infrastructure/file/repository"

type FileStorageInterface interface {
	WriteImage(image []byte, id string, level string) error
	ReadImage(id string, level string) ([]byte, error)
}

type fileStorageService struct {
	fileStorage repository.FileRepositoryInterface
}

var _ FileStorageInterface = (*fileStorageService)(nil)

func New(storage repository.FileRepositoryInterface) *fileStorageService {
	return &fileStorageService{
		fileStorage: storage,
	}
}

func (f *fileStorageService) WriteImage(image []byte, name string, level string) error {
	if err := f.fileStorage.CreateImage(image, name, level); err != nil {
		return err
	}
	return nil
}

func (f fileStorageService) ReadImage(name string, level string) ([]byte, error) {
	data, err := f.fileStorage.GetImage(name, level)
	if err != nil {
		return nil, err
	}
	return data, nil
}
