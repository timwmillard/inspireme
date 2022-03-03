package inspireme

import (
	"context"
	"encoding/json"
	"image"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

const mockStorageImgURL = "https://storage.googleapis.com/inspireme/1234-5678.png"

type MockStorage struct{}

func (t *MockStorage) Store(ctx context.Context, img image.Image, format string) (string, error) {
	return mockStorageImgURL, nil
}

func TestGenerateHandler(t *testing.T) {

	// testCases := []struct {
	// 	inputJSON string
	// 	success   bool
	// }{
	// 	{inputJSON: testReqSuccess, success: true},
	// 	{inputJSON: testReqFailed, success: false},
	// }

	handler := &Handler{
		Log: log.New(os.Stdout, "test", log.LstdFlags),
		InspireMe: &ImageGenerator{
			Client:   http.DefaultClient,
			Storage:  &MockStorage{},
			FontsDir: "../../resources/fonts",
		},
	}
	ts := httptest.NewServer(handler)

	resp, err := http.Post(ts.URL, "application/json", strings.NewReader(testReqSuccess))
	if err != nil {
		t.Errorf("failed generate handler: %v", err)
	}

	var genImg struct {
		ImageURL string `json:"imageUrl"`
	}
	encoder := json.NewDecoder(resp.Body)
	encoder.Decode(&genImg)

	if genImg.ImageURL != mockStorageImgURL {
		t.Errorf("return URL failed wanted %s but got %s", mockStorageImgURL, genImg.ImageURL)
	}
}

const (
	testReqSuccess = `{
	"quote": "Success",
	"backgroundUrl":"https://images.theweek.com/sites/default/files/styles/tw_image_9_4/public/iStock_95204155_LARGE.jpg"
}`

	testReqFailed = `{
	"quote": "Failed",
	"backgroundUrl":"https://github.com/timwmillard/inspireme/fackimage.png"
}`
)
