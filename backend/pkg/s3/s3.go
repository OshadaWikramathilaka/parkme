package s3

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

type S3Client struct {
	client     *s3.Client
	bucketName string
}

func NewS3Client(bucketName string) (*S3Client, error) {
	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = "ap-south-1" // fallback default
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(
			os.Getenv("AWS_ACCESS_KEY_ID"),
			os.Getenv("AWS_SECRET_ACCESS_KEY"),
			"",
		))),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config: %v", err)
	}

	client := s3.NewFromConfig(cfg)
	return &S3Client{
		client:     client,
		bucketName: bucketName,
	}, nil
}

func (s *S3Client) GeneratePresignedURL(filename, fileType, folderName string) (string, string, string, error) {
	// Split filename and extension
	ext := filepath.Ext(filename)
	nameWithoutExt := strings.TrimSuffix(filename, ext)

	// Generate UUID and create new filename
	uniqueID := uuid.New().String()
	uniqueFilename := fmt.Sprintf("%s-%s%s", nameWithoutExt, uniqueID, ext)

	key := fmt.Sprintf("%s/%s", folderName, uniqueFilename)

	presignClient := s3.NewPresignClient(s.client)

	// Create the PutObject input
	putInput := &s3.PutObjectInput{
		Bucket:      aws.String(s.bucketName),
		Key:         aws.String(key),
		ContentType: aws.String(fileType),
	}

	request, err := presignClient.PresignPutObject(context.TODO(), putInput,
		s3.WithPresignExpires(time.Minute*15),
	)

	if err != nil {
		return "", "", "", fmt.Errorf("failed to generate presigned URL: %v", err)
	}

	// Generate the final S3 URL using environment region
	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = "ap-south-1" // fallback default
	}

	// Use direct S3 URL format
	s3URL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s.bucketName, region, key)

	return request.URL, uniqueFilename, s3URL, nil
}

func (s *S3Client) GetBucketName() string {
	return s.bucketName
}
