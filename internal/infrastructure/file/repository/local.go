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

// New returns an instance of localFileStorage struct, which implements the FileRepository interface.
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

// getOrCreateDir creates the directory by path, returns errors
func (localFileStorage) getOrCreateDir(path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		if err = os.MkdirAll(path, os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}

// getPathOfFile creates the path name for file
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

	// Get the directory path for the image
	id_path := filepath.Join(l.directoryPath, id)

	// Create the directory if it does not exist already
	err := l.getOrCreateDir(id_path)
	if err != nil {
		l.logger.Error("Error on creating folder", logger.M{
			"error": err,
		})
		return fmt.Errorf("can't create a directory '%s': %s", id_path, err)
	}

	// Get the file path for the image
	path := l.getPathOfFile(data, id, level)

	// Write the image data to the file
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
		// Check if the found item is not a directory
		if !info.IsDir() {
			// Get the name of the found item
			name := info.Name()
			file_name := strings.TrimSuffix(filepath.Base(name), filepath.Ext(name))

			// Check if the found file has the same name as the specified one
			if file_name == fileName {
				result = path // Store the full path of the found file
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

	// Return the full path of the found file
	return result, nil
}

func (l *localFileStorage) GetImage(id string, level string) ([]byte, error) {
	l.logger.Debug("Trying to get image", logger.M{
		"id":    id,
		"level": level,
	})

	// Look for the file with the specified name in the specified directory
	pathImage, err := l.findFileByName(
		filepath.Join(l.directoryPath, id),
		level,
	)

	// If an error occurred during the file search, log it and return it as an error
	if err != nil {
		l.logger.Error("Error finding image file", logger.M{
			"error": err,
			"id":    id,
			"level": level,
		})
		return nil, err
	}

	// Read the file contents into a byte slice
	data, err := ioutil.ReadFile(pathImage)
	if err != nil {
		// If an error occurred during file reading, log it and return it as an error
		l.logger.Error("Error reading image file", logger.M{
			"error": err,
			"id":    id,
			"level": level,
		})
		return nil, fmt.Errorf("can't read the file '%s': %s", id, err)
	}

	// Log that we successfully retrieved an image, along with the ID and quality level
	l.logger.Info("Successfully retrieved image", logger.M{
		"id":    id,
		"level": level,
	})

	// Return the file contents as a byte slice
	return data, nil
}
