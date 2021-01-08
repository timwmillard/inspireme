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
	fileStoragePath := os.Getenv("FILE_STORAGE_PATH")
	if fileStoragePath == "" {
		fileStoragePath = "."
	}

	// InspireMe Image Generator
	imgGen := &inspireme.ImageGenerator{
		Storage:  &storage.File{Dir: fileStoragePath},
		Client:   http.DefaultClient,
		FontsDir: fontsDir,
	}

	// InspireMe HTTP Hander
	inspiremeHandler := &inspireme.Handler{
		Log:       logger,
		InspireMe: imgGen,
	}

	// HTTP Server
	server := http.Server{
		Addr:         bindAddress,
		Handler:      inspiremeHandler,
		ErrorLog:     logger,
		ReadTimeout:  5 * time.Second,   // max time to read request from the client
		WriteTimeout: 20 * time.Second,  // max time to write response to the client
		IdleTimeout:  120 * time.Second, // max time for connections using TCP Keep-Alive
	}

	// Run the serve in a go routine
	go func() {
		logger.Printf("Starting server at %s", bindAddress)

		err := server.ListenAndServe()
		if err != nil {
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
