package inspireme_test

import (
	"fmt"
	"inspireme"
	"testing"
)

func TestInspirationifier(t *testing.T) {

	quote := "I love Go"
	backgroundURL := "https://images.theweek.com/sites/default/files/styles/tw_image_9_4/public/iStock_95204155_LARGE.jpg"

	imageURL, err := inspireme.Generate(quote, backgroundURL, nil)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}
	fmt.Printf("Quote: %s\n", quote)
	fmt.Printf("Background URL: %s\n", backgroundURL)
	fmt.Printf("Result URL: %s\n", imageURL)

}
