package services

import (
	vision "cloud.google.com/go/vision/apiv1"
	"context"
	"fmt"
	"mime/multipart"
)

type Coordinate struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
}

type Object struct {
	Name           string       `json:"name"`
	Score          float32      `json:"score"`
	DetectedByUser bool         `json:"detectedByUser"`
	Coordinates    []Coordinate `json:"coordinates"`
}

// DetectObjects gets objects and bounding boxes from the Vision API for an image at the given file path.
func DetectObjects(file multipart.File) ([]*Object, error) {
	ctx := context.Background()
	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("NewClient: %v", err)

	}

	image, err := vision.NewImageFromReader(file)
	if err != nil {
		return nil, fmt.Errorf("NewImageReader: %v", err)
	}
	annotations, err := client.LocalizeObjects(ctx, image, nil)
	if err != nil {
		return nil, fmt.Errorf("LocalizeObject: %v", err)
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
