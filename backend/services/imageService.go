package services

import (
	vision "cloud.google.com/go/vision/apiv1"
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"os"
)

type Coordinate struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
}

type Object struct {
	Name string `json:"name"`
	Score float32 `json:"score"`
	DetectedByUser bool      `json:"DetectedByUser"`
	Coordinates []Coordinate `json:"coordinates"`
}

// DetectObjects gets objects and bounding boxes from the Vision API for an image at the given file path.
func DetectObjects(file multipart.File) ([]*Object, error) {
	ctx := context.Background()
	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		return nil, err
	}

	image, err := vision.NewImageFromReader(file)
	if err != nil {
		return nil, err
	}
	annotations, err := client.LocalizeObjects(ctx, image, nil)
	if err != nil {
		return nil, err
	}

	if len(annotations) == 0 {
		return nil, nil
	}

	objectList := make([]*Object, 0)

	for _, annotation := range annotations {
		coordinates := make([]Coordinate, 0, 4)
		for _, v := range annotation.BoundingPoly.NormalizedVertices {
			coordinates = append(coordinates, Coordinate{v.X, v.Y})
		}
		objectList = append(objectList, &Object{annotation.Name, annotation.Score, false, coordinates})
	}

	return objectList, nil
}

func objectLabels(filename string) error {
	ctx := context.Background()

	// Creates a client.
	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		log.Printf("Failed to create client: %v", err)
		return err
	}
	defer client.Close()

	file, err := os.Open(filename)
	if err != nil {
		log.Printf("Failed to read file: %v", err)
		return err
	}
	defer file.Close()
	image, err := vision.NewImageFromReader(file)
	if err != nil {
		log.Printf("Failed to create image: %v", err)
		return err
	}

	labels, err := client.DetectLabels(ctx, image, nil, 10)
	if err != nil {
		log.Printf("Failed to detect labels: %v", err)
		return err
	}

	fmt.Println("Labels:")
	for _, label := range labels {
		fmt.Printf("%v : %v\n", label.Description, label.Score)
	}
	return nil
}
