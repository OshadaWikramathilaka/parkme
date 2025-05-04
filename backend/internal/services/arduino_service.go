package services

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	_ "image/jpeg" // Add support for JPEG
	_ "image/png"  // Add support for PNG

	localconfig "github.com/dfanso/parkme-backend/config"
	"github.com/dfanso/parkme-backend/internal/models"
	"github.com/dfanso/parkme-backend/internal/repositories"
	"go.mongodb.org/mongo-driver/bson"
)

var (
	ErrModuleNotFound    = errors.New("module not found")
	ErrInvalidSensorData = errors.New("invalid sensor data")
	ErrModuleUnavailable = errors.New("module is currently unavailable")
)

type SensorData struct {
	ModuleID  string                 `json:"moduleId"`
	Timestamp time.Time              `json:"timestamp"`
	Type      string                 `json:"sensorType"`
	Readings  map[string]interface{} `json:"readings"`
}

type AlertData struct {
	ModuleID  string                 `json:"moduleId"`
	Type      string                 `json:"alertType"`
	Severity  string                 `json:"severity"`
	Timestamp time.Time              `json:"timestamp"`
	Details   map[string]interface{} `json:"details"`
}

type ModuleStatus struct {
	ModuleID     string    `json:"moduleId"`
	Status       string    `json:"status"`
	LastSeen     time.Time `json:"lastSeen"`
	BatteryLevel int       `json:"batteryLevel"`
	Error        string    `json:"error,omitempty"`
}

type ModuleConfig struct {
	SampleRate    int     `json:"sampleRate"`
	Threshold     float64 `json:"threshold"`
	AlertEnabled  bool    `json:"alertEnabled"`
	LedBrightness int     `json:"ledBrightness"`
}

type NumberPlateResult struct {
	Text       string          `json:"text"`
	IsValid    bool            `json:"isValid"`
	FilePath   string          `json:"filePath,omitempty"`
	Confidence float64         `json:"confidence,omitempty"`
	Vehicle    *models.Vehicle `json:"vehicle,omitempty"`
}

type ArduinoService struct {
	config         *localconfig.Config
	rekognition    *RekognitionService
	vehicleService *VehicleService
}

func NewArduinoService(cfg *localconfig.Config, vehicleService *VehicleService) (*ArduinoService, error) {
	rekognitionSvc, err := NewRekognitionService(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create rekognition service: %w", err)
	}

	return &ArduinoService{
		config:         cfg,
		rekognition:    rekognitionSvc,
		vehicleService: vehicleService,
	}, nil
}

func (s *ArduinoService) ExtractNumberPlate(data multipart.File) (*NumberPlateResult, error) {
	ctx := context.Background()

	// Use Rekognition service
	result, err := s.rekognition.DetectText(ctx, data)
	if err != nil {
		return nil, err
	}

	// Clean and validate the text
	text := cleanNumberPlateText(result.Text)
	isValid := isValidNumberPlate(text)

	// Find the vehicle by plate number
	vehicle, err := s.vehicleService.FindOne(ctx, bson.M{"plate_number": text})
	if err != nil {
		if err == repositories.ErrNotFound {
			return nil, ErrVehicleNotFound
		}
		return nil, err
	}

	// Vehicle exists, return the result without error
	return &NumberPlateResult{
		Text:       text,
		IsValid:    isValid,
		Confidence: result.Confidence,
		Vehicle:    vehicle,
	}, nil
}

func (s *ArduinoService) SaveImageToFolder(data multipart.File, filename string) (string, error) {
	// Define the folder to save images
	saveDir := "./uploads"

	// Create the directory if it doesn't exist
	err := os.MkdirAll(saveDir, os.ModePerm)
	if err != nil {
		return "", err
	}

	// Define the full path to save the image
	savePath := filepath.Join(saveDir, filename)

	// Create a new file to save the image
	destFile, err := os.Create(savePath)
	if err != nil {
		return "", err
	}
	defer destFile.Close()

	// Copy the image data to the new file
	_, err = io.Copy(destFile, data)
	if err != nil {
		return "", err
	}

	return savePath, nil
}

// Data processing methods

func cleanNumberPlateText(text string) string {
	// Remove whitespace and convert to uppercase
	reg := regexp.MustCompile(`[^A-Z0-9]`)
	text = strings.ToUpper(strings.TrimSpace(text))
	return reg.ReplaceAllString(text, "")
}

func isValidNumberPlate(text string) bool {
	if len(text) < 5 || len(text) > 10 {
		return false
	}

	// Basic pattern for number plates (customize based on your country's format)
	// This example assumes a format like "ABC123" or "AB12CDE"
	pattern := `^[A-Z0-9]{5,10}$`
	matched, err := regexp.MatchString(pattern, text)
	if err != nil {
		return false
	}

	// Ensure there's at least one letter and one number
	hasLetter := regexp.MustCompile(`[A-Z]`).MatchString(text)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(text)

	return matched && hasLetter && hasNumber
}

// FindVehicleByPlate is a helper method to find a vehicle by plate number
func (s *ArduinoService) FindVehicleByPlate(ctx context.Context, plateNumber string) (*models.Vehicle, error) {
	return s.vehicleService.FindOne(ctx, bson.M{"plate_number": plateNumber})
}
