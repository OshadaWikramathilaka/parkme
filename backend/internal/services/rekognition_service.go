package services

import (
	"context"
	"fmt"
	"io"
	"regexp"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/rekognition"
	"github.com/aws/aws-sdk-go-v2/service/rekognition/types"
	localconfig "github.com/dfanso/parkme-backend/config"
)

type RekognitionService struct {
	config *localconfig.Config
	client *rekognition.Client
}

type RekognitionResult struct {
	Text       string  `json:"text"`
	Confidence float64 `json:"confidence"`
}

func NewRekognitionService(cfg *localconfig.Config) (*RekognitionService, error) {
	ctx := context.Background()
	awsCfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(cfg.AWS.Region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.AWS.AccessKeyID,
			cfg.AWS.SecretAccessKey,
			"",
		)),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to load AWS config: %w", err)
	}

	return &RekognitionService{
		config: cfg,
		client: rekognition.NewFromConfig(awsCfg),
	}, nil
}

func (s *RekognitionService) DetectText(ctx context.Context, data io.Reader) (*RekognitionResult, error) {
	// Read image data
	imageBytes, err := io.ReadAll(data)
	if err != nil {
		return nil, fmt.Errorf("failed to read image: %w", err)
	}

	input := &rekognition.DetectTextInput{
		Image: &types.Image{
			Bytes: imageBytes,
		},
	}

	// Detect text
	output, err := s.client.DetectText(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to detect text: %w", err)
	}

	// Debug print all detections
	fmt.Println("All detected text:")
	for _, detection := range output.TextDetections {
		fmt.Printf("Type: %s, Text: %s, Confidence: %f\n",
			detection.Type, *detection.DetectedText, *detection.Confidence)
	}

	var plateText string
	var bestConfidence float64

	// Find the code (excluding province codes like WP, SP, etc.)
	for _, detection := range output.TextDetections {
		if detection.Type == "WORD" {
			text := *detection.DetectedText
			// Skip known province codes
			if text != "WP" && text != "SP" && text != "CP" && text != "EP" && text != "NP" && text != "UP" && text != "NW" {
				// Check if it's exactly 2 or 3 uppercase letters
				if (len(text) == 2 || len(text) == 3) && regexp.MustCompile(`^[A-Z]{2,3}$`).MatchString(text) {
					plateText = text
					bestConfidence = float64(*detection.Confidence)
					break
				}
			}
		}
	}

	// Then find the 4-digit number
	if plateText != "" {
		for _, detection := range output.TextDetections {
			text := *detection.DetectedText
			if detection.Type == "LINE" || detection.Type == "WORD" {
				cleanText := regexp.MustCompile(`\s+`).ReplaceAllString(text, "")
				if regexp.MustCompile(`^\d{4}$`).MatchString(cleanText) {
					plateText += " " + cleanText
					bestConfidence = (bestConfidence + float64(*detection.Confidence)) / 2
					break
				}
			}
		}
	}

	return &RekognitionResult{
		Text:       plateText,
		Confidence: bestConfidence,
	}, nil
}
