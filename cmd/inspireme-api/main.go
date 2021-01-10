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

var (
	mux              = http.NewServeMux()
	logger           = log.New(os.Stdout, "inspireme ", log.LstdFlags)
	bindAddress      = os.Getenv("BIND_ADDRESS")
	port             = os.Getenv("PORT")
	inspireMeStorage = os.Getenv("INSPIREME_STORAGE")
)

func main() {

	// Bind Address
	if bindAddress == "" {
		if port == "" {
			port = "8080"
		}
		bindAddress = ":" + port
	}

	// Fonts Directory
	fontsDir := os.Getenv("FONTS_DIR")
	if fontsDir == "" {
		fontsDir = "../../resources/fonts"
	}

	// InspireMe Image Generator
	imgGen := &inspireme.ImageGenerator{
		Storage:  imageStorage(),
		Client:   http.DefaultClient,
		FontsDir: fontsDir,
	}

	// InspireMe HTTP Hander
	inspiremeHandler := &inspireme.Handler{
		Log:       logger,
		InspireMe: imgGen,
	}
	mux.Handle("/", corsHandler(inspiremeHandler))

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

func imageStorage() inspireme.Storage {

	switch inspireMeStorage {
	case "local":
		return webLocalStorage()
	case "gcloud":
		return gCloudStorage()
	// case "s3":
	// 	return s3Storage()
	default:
		return webLocalStorage()
	}

}

func webLocalStorage() *storage.WebLocalFile {

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

	// Setup Local image server
	imageServer := http.FileServer(http.Dir(imagesStoragePath))
	mux.Handle("/images/", http.StripPrefix("/images/", imageServer))

	return webFileStorage
}

func gCloudStorage() *storage.GCloud {

	return &storage.GCloud{
		ProjectID: os.Getenv("GCLOUD_PROJECT_ID"),
		Bucket:    os.Getenv("GCLOUD_BUCKET"),
	}
}

// func s3Storage() *storage.S3 {
// 	return nil
// }

// CorsHandler handles cross origin request
func corsHandler(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Origin, Accept")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		h.ServeHTTP(w, r)
	}
}
