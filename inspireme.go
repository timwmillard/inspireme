package inspireme

import (
	"context"
	"errors"
	"fmt"
	"image"
	"io"
	"net/http"
	"path/filepath"

	"image/color"
	"image/jpeg"
	"image/png"

	"github.com/fogleman/gg"
)

// ImageGenerator is used the generate the images.
// call the Genrate method for all the magic to happen.
type ImageGenerator struct {
	Storage  Storage
	Client   *http.Client
	FontsDir string
}

// Storage is an interface to store the image
type Storage interface {
	Store(img image.Image, format string) (string, error)

	// TOOD add delete method
	// Delete(url string) error

	// Was trying to make the interface return a WriteCloser
	// WriteCloser(format string) (io.WriteCloser, error)
}

// Generate is were all the magic happends
// it fetches the background image from the URL and places
// the quote as an overlay.  Use the style to set style
// the overlay (setting it to null will use defautl style)
// It will return URL of the image
func (im *ImageGenerator) Generate(ctx context.Context, quote, backgroundURL string, style *Styles) (image.Image, string, error) {
	// fetch the backgroudn image
	bgImg, imgFormat, err := im.FetchImage(ctx, backgroundURL)
	if err != nil {
		return nil, "", err
	}

	// generate the image with the text
	img, err := im.GenerateImage(quote, bgImg, style)
	if err != nil {
		return nil, "", err
	}

	return img, imgFormat, nil
}

// GenerateImage will add the message to the background image
func (im *ImageGenerator) GenerateImage(message string, background image.Image, style *Styles) (image.Image, error) {

	// create a drawing context
	draw := gg.NewContextForImage(background)

	// Draw dark overlay so text is readable
	draw.SetColor(color.RGBA{0, 0, 0, 204})
	draw.DrawRectangle(0, 0, float64(draw.Width()), float64(draw.Height()))
	draw.Fill()

	// Load font
	fontPath := filepath.Join(im.FontsDir, "OpenSans-Bold.ttf")
	draw.LoadFontFace(fontPath, 80) // ignore error, as just just default system font if font not found.

	// Draw text message
	draw.SetColor(color.White)
	const padding = 20
	draw.DrawStringWrapped(message, float64(draw.Width()/2), float64(draw.Height()/2), 0.5, 0.5, float64(draw.Width()-padding), 1.5, gg.AlignCenter)

	// switch imgFormat {
	// case "png":
	// 	draw.EncodePNG(output)
	// case "jpeg":
	// 	draw.EncodeJPG(output, nil)
	// default:
	// 	return errors.New("invalid image format")
	// }

	return draw.Image(), nil
}

// GenerateAndStore this is the main
func (im *ImageGenerator) GenerateAndStore(ctx context.Context, quote, backgroundURL string, style *Styles) (string, error) {

	img, imgFormat, err := im.Generate(ctx, quote, backgroundURL, style)
	if err != nil {
		return "", err
	}

	url, err := im.Storage.Store(img, imgFormat)
	if err != nil {
		return "", err
	}

	return url, nil
}

// Store the image
// func (im *ImageGenerator) Store(img image.Image, format string) (string, error) {
// 	writer := im.storage.WriteCloser(format)

// 	defer writer.Close()

// 	switch format {
// 	case "png":
// 		png.Encode(writer, img)
// 	case "jpeg":
// 		jpeg.Encode(writer, img, nil)
// 	default:
// 		return "", errors.New("invalid image formate")
// 	}

// 	return im.storage.Location(), nil
// }

// FetchImage will go to the image URL an return an in memory image
func (im *ImageGenerator) FetchImage(ctx context.Context, url string) (image.Image, string, error) {

	resp, err := im.Client.Get(url)
	if err != nil {
		return nil, "", fmt.Errorf("unable to fetch image URL: %v", err)
	}
	defer resp.Body.Close()

	img, imgFormat, err := image.Decode(resp.Body)
	if err != nil {
		return nil, "", fmt.Errorf("unable decode image from body: %v", err)
	}

	return img, imgFormat, nil
}

// EncodeImage utily to encode image in the given format
func EncodeImage(w io.WriteCloser, img image.Image, format string) error {
	defer w.Close()

	switch format {
	case "png":
		png.Encode(w, img)
	case "jpeg":
		jpeg.Encode(w, img, nil)
	default:
		return errors.New("invalid image formate")
	}

	return nil
}
