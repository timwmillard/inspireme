package storage

import (
	"image"
	"os"
)

// File is a local file storage for images
type File struct {
	file os.File
}

// Store the image on the local file system
func (f *File) Store(i image.Image) (string, error) {

	return "nil", nil
}
