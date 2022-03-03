package inspireme_test

import (
	"bytes"
	"context"
	"crypto/md5"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"
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

func TestImageGeneration(t *testing.T) {
	imgGen := inspireme.ImageGenerator{
		Client: http.DefaultClient,
		// Storage:  &storage.LocalFile{FileName: "testdata/test.png"},
		FontsDir: "resources/fonts",
	}

	// const (
	// 	quote           = "Test"
	// 	bgURL           = "https://img.freepik.com/free-photo/wall-wallpaper-concrete-colored-painted-textured-concept_53876-31799.jpg"
	// 	testImgFileName = "testdata/test.jpeg"
	// )

	const (
		quote           = "Test"
		bgURL           = "https://images.template.net/wp-content/uploads/2015/08/Extraordinary-Paper-Background-for-Free.png"
		testImgFileName = "testdata/test.png"
	)

	ctx := context.Background()

	// Got Image
	gotImg, gotFmt, err := imgGen.Generate(ctx, quote, bgURL, nil)
	if err != nil {
		t.Fatalf("Error generating image: %v", err)
	}

	// Store image in testdata folder (only run this once)
	// imgGen.Storage.Store(ctx, gotImg, gotFmt)

	// Want Image - from testdata
	file, err := os.Open(testImgFileName)
	if err != nil {
		t.Fatalf("Error opening test image: %v", err)
	}
	wantImg, wantFmt, err := image.Decode(file)
	if err != nil {
		t.Fatalf("Error decoding test image: %v", err)
	}

	if wantFmt != gotFmt {
		t.Errorf("wrong image format, wanted %x but got %x", wantFmt, gotFmt)
	}

	gotSum := checksumImage(t, gotImg, wantFmt)
	wantSum := checksumImage(t, wantImg, wantFmt)

	if gotSum != wantSum {
		t.Errorf("image checksum, wanted %x but got %x", wantSum, gotSum)
	}
}

func checksumImage(t *testing.T, img image.Image, fmt string) string {
	var err error

	buf := &bytes.Buffer{}

	switch fmt {
	case "png":
		err = png.Encode(buf, img)
	case "jpeg":
		err = jpeg.Encode(buf, img, nil)
	default:
		t.Fatalf("Invalid test image format: %s", fmt)
	}
	if err != nil {
		t.Fatalf("Unable to endcode image: %v", err)
	}

	hash := md5.New()
	_, err = io.Copy(hash, buf)
	if err != nil {
		t.Fatalf("Error hashing image: %v", err)
	}
	return string(hash.Sum(nil))
}
