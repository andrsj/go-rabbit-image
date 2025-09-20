package repository

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/andrsj/go-rabbit-image/pkg/logger"
)

type noopLogger struct{}

func (l noopLogger) Named(string) logger.Logger { return l }
func (noopLogger) Debug(string, logger.M)       {}
func (noopLogger) Info(string, logger.M)        {}
func (noopLogger) Warn(string, logger.M)        {}
func (noopLogger) Error(string, logger.M)       {}
func (noopLogger) Fatal(string, logger.M)       {}

var png1x1 = []byte{
	0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a,
	0x00, 0x00, 0x00, 0x0d, 0x49, 0x48, 0x44, 0x52,
	0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
	0x08, 0x02, 0x00, 0x00, 0x00, 0x90, 0x77, 0x53,
	0xde, 0x00, 0x00, 0x00, 0x0b, 0x49, 0x44, 0x41,
	0x54, 0x08, 0xd7, 0x63, 0xf8, 0xcf, 0xc0, 0x00,
	0x00, 0x04, 0x00, 0x01, 0xe2, 0x26, 0x05, 0x9b,
	0x00, 0x00, 0x00, 0x00, 0x49, 0x45, 0x4e, 0x44,
	0xae, 0x42, 0x60, 0x82,
}

func TestNewCreatesDirectory(t *testing.T) {
	t.Parallel()

	baseDir := filepath.Join(t.TempDir(), "server_images")
	if _, err := New(baseDir, noopLogger{}); err != nil {
		t.Fatalf("expected no error creating repository: %v", err)
	}

	if info, err := os.Stat(baseDir); err != nil {
		t.Fatalf("expected directory to be created: %v", err)
	} else if !info.IsDir() {
		t.Fatalf("expected %s to be a directory", baseDir)
	}
}

func TestCreateImageCreatesFile(t *testing.T) {
	t.Parallel()

	baseDir := filepath.Join(t.TempDir(), "server_images")
	repo, err := New(baseDir, noopLogger{})
	if err != nil {
		t.Fatalf("expected no error creating repository: %v", err)
	}

	const (
		imageID = "test-image"
		level   = "75"
	)

	if err := repo.CreateImage(png1x1, imageID, level); err != nil {
		t.Fatalf("expected no error creating image: %v", err)
	}

	expectedPath := filepath.Join(baseDir, imageID, level+".png")
	if _, err := os.Stat(expectedPath); err != nil {
		t.Fatalf("expected image file to exist: %v", err)
	}
}
