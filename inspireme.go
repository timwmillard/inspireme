package inspireme

import (
	"context"
	"errors"
	"fmt"
	"image"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"image/color"

	"github.com/fogleman/gg"
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
func Generate(ctx context.Context, c *http.Client, quote, backgroundURL string, style *Styles, output io.Writer) error {
	bgImg, imgFormat, err := fetchImage(ctx, c, backgroundURL)
	if err != nil {
		return err
	}
	_ = imgFormat

	// bgBonds := bgImg.Bounds()

	draw := gg.NewContextForImage(bgImg)

	// Draw dark overlay so text is readable
	draw.SetColor(color.RGBA{0, 0, 0, 204})
	draw.DrawRectangle(0, 0, float64(draw.Width()), float64(draw.Height()))
	draw.Fill()

	// Load font
	fontDirPath := os.Getenv("FONTS_DIR") // TODO this should be moved to main func and passed in via a config struct
	fontPath := filepath.Join(fontDirPath, "OpenSans-Bold.ttf")
	draw.LoadFontFace(fontPath, 80) // ignore error, as just just default system font if font not found.

	// Draw text message
	draw.SetColor(color.White)
	const padding = 20
	draw.DrawStringWrapped(quote, float64(draw.Width()/2), float64(draw.Height()/2), 0.5, 0.5, float64(draw.Width()-padding), 1.5, gg.AlignCenter)

	switch imgFormat {
	case "png":
		draw.EncodePNG(output)
	case "jpeg":
		draw.EncodeJPG(output, nil)
	default:
		return errors.New("invalid image format")
	}

	return nil
}

// GenerateAndStore -
func GenerateAndStore(quote, backgroundURL string, style *Styles) (string, error) {
	return "", nil
}

func store(image.Image) (string, error) {
	return "", nil
}

// FetchImage will go to the image URL an return an in memory image
func fetchImage(ctx context.Context, c *http.Client, url string) (image.Image, string, error) {

	resp, err := c.Get(url)
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
