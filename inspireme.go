package inspireme

import (
	"image"
	"net/http"
	"net/url"

	// _ "image/png"
	_ "image/jpeg"
)

// Storage is an interface to store the image
type Storage interface {
	Store(image.Image) (string, error)
	Delete(url string) error
}

// Generate is were all the magic happends
// it fetches the background image from the URL and places
// the quote as an overlay.  Use the style to set style
// the overlay (setting it to null will use defautl style)
// It will return URL of the image
func Generate(quote, backgroundURL string, style *Styles) (image.Image, error) {

	return nil, nil
}

// GenerateAndStore -
func GenerateAndStore(quote, backgroundURL string, style *Styles) (string, error) {
	return "", nil
}

func store(image.Image) (string, error) {
	return "", nil
}

// FetchImage will go to the image URL an return an in memory image
func fetchImage(client *http.Client, url url.URL) (image.Image, error) {

	return nil, nil
}
