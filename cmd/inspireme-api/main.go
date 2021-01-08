package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/timwmillard/inspireme"
	"github.com/timwmillard/inspireme/storage"
)

func main() {
	logger := log.New(os.Stdout, "inspireme ", log.LstdFlags)

	// Bind Address
	bindAddress := os.Getenv("BIND_ADDRESS")
	if bindAddress == "" {
		bindAddress = ":8080"
	}

	// Fonts Directory
	fontsDir := os.Getenv("FONTS_DIR")
	if fontsDir == "" {
		fontsDir = "../../resources/fonts"
	}

	// File Storage
	imagesStoragePath := os.Getenv("IMAGES_STORAGE_PATH")
	if imagesStoragePath == "" {
		imagesStoragePath = "images"
	}
	// Create directory if not exist
	err := os.MkdirAll(imagesStoragePath, os.ModePerm)
	if err != nil {
		logger.Printf("unable to create image storage path: %s", err)
		os.Exit(1)
	}

	// File Storage
	imagesBaseURL := os.Getenv("IMAGES_BASE_URL")
	if imagesBaseURL == "" {
		imagesBaseURL = "http://localhost" + bindAddress + "/images/"
	}

	webFileStorage := &storage.WebLocalFile{
		LocalDir:      imagesStoragePath,
		ImagesBaseURL: imagesBaseURL,
	}

	// InspireMe Image Generator
	imgGen := &inspireme.ImageGenerator{
		Storage:  webFileStorage,
		Client:   http.DefaultClient,
		FontsDir: fontsDir,
	}

	// Create mux handler
	mux := http.NewServeMux()

	// InspireMe HTTP Hander
	inspiremeHandler := &inspireme.Handler{
		Log:       logger,
		InspireMe: imgGen,
	}
	mux.Handle("/", inspiremeHandler)

	imageServer := http.FileServer(http.Dir(imagesStoragePath))
	mux.Handle("/images/", http.StripPrefix("/images/", imageServer))

	// HTTP Server
	server := http.Server{
		Addr:         bindAddress,
		Handler:      mux,
		ErrorLog:     logger,
		ReadTimeout:  5 * time.Second,   // max time to read request from the client
		WriteTimeout: 20 * time.Second,  // max time to write response to the client
		IdleTimeout:  120 * time.Second, // max time for connections using TCP Keep-Alive
	}

	// Run the serve in a go routine
	go func() {
		logger.Printf("Starting server at %s", bindAddress)

		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			logger.Printf("error starting server: %s", err)
			os.Exit(1)
		}
	}()

	exitCh := make(chan os.Signal, 1)
	signal.Notify(exitCh, os.Interrupt)
	signal.Notify(exitCh, os.Kill)

	sig := <-exitCh
	logger.Printf("performing gracefull shutdown due to: %v", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	cancel()
	server.Shutdown(ctx)

}
