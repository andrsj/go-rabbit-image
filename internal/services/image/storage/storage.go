package storage

import "github.com/andrsj/go-rabbit-image/internal/domain/repositories/file"

type FileStorage interface {
	WriteImage(image []byte, id string, level string) error
	ReadImage(id string, level string) ([]byte, error)
}

type fileStorageService struct {
	fileStorage file.FileRepository
}

var _ FileStorage = (*fileStorageService)(nil)

func New(storage file.FileRepository) *fileStorageService {
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
