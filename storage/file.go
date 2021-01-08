package storage

import (
	"fmt"
	"image"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/timwmillard/inspireme"
)

// File is a local file storage for images
// dir contains the directory the files will be saved
type File struct {
	Dir string
}

// Store the image to the local file system
func (f *File) Store(img image.Image, format string) (string, error) {

	id := uuid.New()
	fileName := id.String() + "." + format
	filePath := filepath.Join(f.Dir, fileName)

	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("unable to create file: %v", err)
	}

	inspireme.EncodeImage(file, img, format)

	return filePath, nil
}
