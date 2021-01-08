package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path"

	"github.com/timwmillard/inspireme"
	"github.com/timwmillard/inspireme/storage"
)

func main() {

	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <message> <backgroundURL>\n", os.Args[0])
		os.Exit(1)
	}

	message := os.Args[1]
	backgroundURL := os.Args[2]

	ext := path.Ext(backgroundURL)

	localFile := &storage.LocalFile{FileName: "inspireme" + ext}

	// Fonts Directory
	fontsDir := os.Getenv("FONTS_DIR")
	if fontsDir == "" {
		fontsDir = "../../resources/fonts"
	}

	imgGen := inspireme.ImageGenerator{
		Storage:  localFile,
		Client:   http.DefaultClient,
		FontsDir: fontsDir,
	}

	_, err := imgGen.GenerateAndStore(context.Background(), message, backgroundURL, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating image: %v\n", err)
		os.Exit(2)
	}
}
