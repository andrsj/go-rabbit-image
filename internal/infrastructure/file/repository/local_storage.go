package repository

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type localFileStorage struct {
	directoryPath string
}

func New(pathToDir string) (*localFileStorage, error) {
	lfs := &localFileStorage{}
	lfs.directoryPath = pathToDir
	err := lfs.getOrCreateDir(pathToDir)
	if err != nil {
		return nil, err
	}
	return lfs, nil
}

func (l *localFileStorage) getOrCreateDir(path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		if err = os.MkdirAll(path, os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}

func (l *localFileStorage) getPathOfFile(data []byte, id string, level string) string {
	var fileExt string
	switch contentType := http.DetectContentType(data); contentType {
	case "image/jpeg":
		fileExt = "jpeg"
	case "image/jpg":
		fileExt = "jpg"
	case "image/png":
		fileExt = "png"
	}
	filename := fmt.Sprintf("%s.%s", level, fileExt)
	path := filepath.Join(l.directoryPath, id, filename)
	return path
}

func (l *localFileStorage) CreateImage(data []byte, id string, level string) error {
	id_path := filepath.Join(l.directoryPath, id)
	err := l.getOrCreateDir(id_path)
	if err != nil {
		return fmt.Errorf("can't create a directory '%s': %s", id_path, err)
	}

	path := l.getPathOfFile(data, id, level)
	err = ioutil.WriteFile(path, data, os.ModePerm)
	if err != nil {
		return fmt.Errorf("can't create an image '%s': %s", id, err)
	}
	return nil
}

func (l *localFileStorage) findFileByName(dirPath, fileName string) (string, error) {
	var result string
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			name := info.Name()
			file_name := strings.TrimSuffix(filepath.Base(name), filepath.Ext(name))

			if file_name == fileName {
				result = path
			}
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	if result == "" {
		return "", fmt.Errorf("file not found with name %s in directory %s", fileName, dirPath)
	}
	return result, nil
}

func (l *localFileStorage) GetImage(id string, level string) ([]byte, error) {
	pathImage, err := l.findFileByName(
		filepath.Join(l.directoryPath, id),
		level,
	)

	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadFile(pathImage)
	if err != nil {
		return nil, fmt.Errorf("can't read the file '%s': %s", id, err)
	}
	return data, nil
}

var _ FileRepositoryInterface = (*localFileStorage)(nil)
