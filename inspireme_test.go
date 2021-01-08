package inspireme_test

import (
	"context"
	"fmt"
	"image"
	"net/http"
	"testing"

	"github.com/timwmillard/inspireme"
)

const fackImageURL = "https://github.com/timwmillard/inspireme/fackimage.png"

type MockStorage struct{}

func (t *MockStorage) Store(ctx context.Context, img image.Image, format string) (string, error) {
	return fackImageURL, nil
}

func TestInspireMe(t *testing.T) {

	quote := "I love Go"
	backgroundURL := "https://images.theweek.com/sites/default/files/styles/tw_image_9_4/public/iStock_95204155_LARGE.jpg"

	imgGen := inspireme.ImageGenerator{
		Client:   http.DefaultClient,
		Storage:  &MockStorage{},
		FontsDir: "../../resources/fonts",
	}

	imageURL, err := imgGen.GenerateAndStore(context.Background(), quote, backgroundURL, nil)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}
	if imageURL != fackImageURL {
		t.Fatalf("invalid return image URL wanted %v but got %v", fackImageURL, imageURL)
	}

	fmt.Printf("Quote: %s\n", quote)
	fmt.Printf("Background URL: %s\n", backgroundURL)
	fmt.Printf("Result URL: %s\n", imageURL)
}
