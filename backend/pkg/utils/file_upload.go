package utils

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/dfanso/parkme-backend/pkg/s3"
)

// UploadFileToS3 handles file upload to S3 and returns the final S3 URL
func UploadFileToS3(file *multipart.FileHeader, s3Client *s3.S3Client, folderName string) (string, error) {
	// Check file type
	contentType := file.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		return "", fmt.Errorf("file must be an image")
	}

	// Open the file
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %v", err)
	}
	defer src.Close()

	// Read file into buffer
	buffer := new(bytes.Buffer)
	if _, err := io.Copy(buffer, src); err != nil {
		return "", fmt.Errorf("failed to read file: %v", err)
	}

	// Generate presigned URL
	presignedURL, _, s3URL, err := s3Client.GeneratePresignedURL(
		file.Filename,
		contentType,
		folderName,
	)
	if err != nil {
		return "", fmt.Errorf("failed to generate upload URL: %v", err)
	}

	// Create HTTP client and request
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPut, presignedURL, bytes.NewReader(buffer.Bytes()))
	if err != nil {
		return "", fmt.Errorf("failed to create upload request: %v", err)
	}

	// Set only Content-Type header
	req.Header.Set("Content-Type", contentType)

	// Upload to S3
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to upload to S3: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Read error response
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to upload to S3: status: %d, body: %s", resp.StatusCode, string(body))
	}

	return s3URL, nil
}
