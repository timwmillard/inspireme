package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/timwmillard/inspireme"
)

func main() {

	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <message> <backgroundURL>\n", os.Args[0])
		os.Exit(1)
	}

	message := os.Args[1]
	backgroundURL := os.Args[2]

	_ = message

	file, err := os.Create("inspireme.png")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating file: %v\n", err)
		os.Exit(3)
	}
	defer file.Close()

	err = inspireme.Generate(context.Background(), http.DefaultClient, message, backgroundURL, nil, file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching image: %v\n", err)
		os.Exit(2)
	}

	// 	err = png.Encode(file, img)
	// 	if err != nil {
	// 		fmt.Fprintf(os.Stderr, "Error encoding png: %v\n", err)
	// 		os.Exit(3)
	// 	}
}
