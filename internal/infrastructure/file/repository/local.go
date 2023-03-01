package repository

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/andrsj/go-rabbit-image/internal/domain/repositories/file"
	"github.com/andrsj/go-rabbit-image/pkg/logger"
)

type localFileStorage struct {
	directoryPath string
	logger        logger.Logger
}

var _ file.FileRepository = (*localFileStorage)(nil)

func New(pathToDir string, log logger.Logger) (*localFileStorage, error) {
	lfs := &localFileStorage{
		directoryPath: pathToDir,
		logger:        log.Named("file repository"),
	}
	lfs.logger.Info("Creating directory", logger.M{
		"path": pathToDir,
	})
	err := lfs.getOrCreateDir(lfs.directoryPath)
	if err != nil {
		lfs.logger.Error("Error on creating folder", logger.M{
			"error": err,
		})
		return nil, err
	}
	return lfs, nil
}

func (localFileStorage) getOrCreateDir(path string) error {
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
	case "image/png":
		fileExt = "png"
	default:
		l.logger.Error("Not accepted content type", logger.M{
			"content type": contentType,
		})
		return ""
	}
	filename := fmt.Sprintf("%s.%s", level, fileExt)
	path := filepath.Join(l.directoryPath, id, filename)
	return path
}

func (l *localFileStorage) CreateImage(data []byte, id string, level string) error {
	l.logger.Info("Creating image", logger.M{
		"id":    id,
		"level": level,
	})

	id_path := filepath.Join(l.directoryPath, id)
	err := l.getOrCreateDir(id_path)
	if err != nil {
		l.logger.Error("Error on creating folder", logger.M{
			"error": err,
		})
		return fmt.Errorf("can't create a directory '%s': %s", id_path, err)
	}

	path := l.getPathOfFile(data, id, level)
	err = ioutil.WriteFile(path, data, os.ModePerm)
	if err != nil {
		l.logger.Error("Error on creating image", logger.M{
			"id":    id,
			"level": level,
			"error": err,
		})
		return fmt.Errorf("can't create an image '%s': %s", id, err)
	}

	l.logger.Info("Image created", logger.M{
		"id":    id,
		"level": level,
	})

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
	l.logger.Debug("Trying to get image", logger.M{
		"id":    id,
		"level": level,
	})

	pathImage, err := l.findFileByName(
		filepath.Join(l.directoryPath, id),
		level,
	)

	if err != nil {
		l.logger.Error("Error finding image file", logger.M{
			"error": err,
			"id":    id,
			"level": level,
		})
		return nil, err
	}

	data, err := ioutil.ReadFile(pathImage)
	if err != nil {
		l.logger.Error("Error reading image file", logger.M{
			"error": err,
			"id":    id,
			"level": level,
		})
		return nil, fmt.Errorf("can't read the file '%s': %s", id, err)
	}

	l.logger.Info("Successfully retrieved image", logger.M{
		"id":    id,
		"level": level,
	})

	return data, nil
}
