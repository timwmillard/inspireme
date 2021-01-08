package storage

import (
	"fmt"
	"image"
	"os"

	"github.com/timwmillard/inspireme"
)

// LocalFile stores the image in a local file
type LocalFile struct {
	FileName string
}

// Store the image to the local file system
func (f *LocalFile) Store(img image.Image, format string) (string, error) {

	file, err := os.Create(f.FileName)
	if err != nil {
		return "", fmt.Errorf("unable to create file: %v", err)
	}

	inspireme.EncodeImage(file, img, format)

	return f.FileName, nil
}
