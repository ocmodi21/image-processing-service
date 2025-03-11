package service

import (
	"errors"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"math/rand"
	"net/http"
	"time"

	"github.com/ocmodi21/image-processing-service/internal/models"
)

// ImageProcessor handles the processing of images
type ImageProcessor struct{}

func NewImageProcessor() *ImageProcessor {
	return &ImageProcessor{}
}

// ProcessImages downloads and processes a list of image URLs
func (p *ImageProcessor) ProcessImages(urls []string) ([]models.ImageResult, error) {
	results := make([]models.ImageResult, 0, len(urls))

	for _, url := range urls {
		result, err := p.processImage(url)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}

	return results, nil
}

// processImage downloads and processes a single image
func (p *ImageProcessor) processImage(url string) (models.ImageResult, error) {
	// Download the image
	resp, err := http.Get(url)
	if err != nil {
		return models.ImageResult{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return models.ImageResult{}, errors.New("failed to download image")
	}

	// Decode the image
	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return models.ImageResult{}, err
	}

	// Calculate the perimeter
	bounds := img.Bounds()
	width := bounds.Max.X - bounds.Min.X
	height := bounds.Max.Y - bounds.Min.Y
	perimeter := 2 * float64(width+height)

	// Random sleep to simulate GPU processing
	sleepTime := 100 + rand.Intn(301) // 100-400 milliseconds
	time.Sleep(time.Duration(sleepTime) * time.Millisecond)

	return models.ImageResult{
		URL:         url,
		Perimeter:   perimeter,
		ProcessedAt: time.Now(),
	}, nil
}
