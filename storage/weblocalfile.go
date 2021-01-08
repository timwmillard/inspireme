package storage

import (
	"context"
	"fmt"
	"image"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/timwmillard/inspireme"
)

// WebLocalFile used for files stored in a static web server
type WebLocalFile struct {
	LocalDir      string
	ImagesBaseURL string
}

// Store the image to the local file system
func (wf *WebLocalFile) Store(ctx context.Context, img image.Image, format string) (string, error) {

	id := uuid.New()
	fileName := id.String() + "." + format
	filePath := filepath.Join(wf.LocalDir, fileName)

	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("unable to create file: %v", err)
	}

	inspireme.EncodeImage(file, img, format)

	imgURL := filepath.Join(wf.ImagesBaseURL, fileName)

	return imgURL, nil
}
