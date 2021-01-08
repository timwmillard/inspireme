# InspireMe

InspireMe is a API to production inspirational images with you faviorite quote overlayed a background image.

1. Find you faviorite Quote
2. Find a cool background image, you just need the URL, so just Google a cool one (just please remeber copyright issues, we take no resposibility).
3. Gererate your image with Inspirationifier and share you inspiration.

## API Server

The API Server will run as a process handling the API requests.

### Usage
As a static image server
```
export BIND_ADDRESS=:8080
export FONTS_DIR=/var/www/resources/fonts
export INSPIREME_STORAGE=local
export IMAGES_STORAGE_PATH=/var/www/images
export IMAGES_BASE_URL=https://<site-domain>/images
./inspireme-api
```

Or store images on Google Cloud Storage
```
export BIND_ADDRESS=:8080
export FONTS_DIR=/var/www/resources/fonts
export INSPIREME_STORAGE=gcloud
export GCLOUD_PROJECT_ID=<gcloud-project-id>
export GCLOUD_BUCKET=<gcloud-bucket-name>
./inspireme-api
```

## API Docs

Please refer to documents, for instruction on the API.
[REST API docs](https://inspireme.docs.apiary.io/)

## CLI
You can use inpireme directly from the command line.
```
export FONTS_DIR=/var/www/resources/fonts
inspireme "My faviorite quote" https://cdn.pixabay.com/photo/2016/11/29/05/45/astronomy-1867616_960_720.jpg
```

## About
InspireMe was created by Tim Millard using Go.  It is a just demo project 😉 .