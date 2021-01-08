package storage

import (
	"context"
	"fmt"
	"image"

	"cloud.google.com/go/storage"
	"github.com/google/uuid"
	"github.com/timwmillard/inspireme"
)

const gCloudBaseURL = "https://storage.googleapis.com/"

// GCloud storage for storing image to Google Cloud Storage
type GCloud struct {
	Credntials string
	ProjectID  string
	Bucket     string
}

// Store the image on gcloud
func (gc *GCloud) Store(ctx context.Context, img image.Image, format string) (string, error) {

	id := uuid.New()
	fileName := id.String() + "." + format

	client, err := storage.NewClient(ctx)
	if err != nil {
		return "", fmt.Errorf("unable to create gcloud client: %v", err)
	}

	bucket := client.Bucket(gc.Bucket)

	// err = bucket.Create(ctx, gc.ProjectID, nil)
	// if err != nil {
	// 	return "", fmt.Errorf("unable to create gcloud bucket: %v", err)
	// }

	obj := bucket.Object(fileName)
	writer := obj.NewWriter(ctx)

	inspireme.EncodeImage(writer, img, format)

	return gCloudBaseURL + obj.BucketName() + "/" + obj.ObjectName(), nil
}
