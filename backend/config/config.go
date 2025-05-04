package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Port string
	}
	MongoDB struct {
		URI  string
		NAME string
	}
	ARDUNIO_API_KEY string
	AWS             struct {
		Region          string
		AccessKeyID     string
		SecretAccessKey string
	}
	ALPR struct {
		APIKey string
	}
	S3BucketName string `mapstructure:"S3_BUCKET_NAME"`
}

func Load() *Config {
	// Load .env file
	godotenv.Load()

	cfg := &Config{}

	// Server configuration
	cfg.Server.Port = getEnv("SERVER_PORT", "8080")

	// MongoDB configuration
	cfg.MongoDB.URI = getEnv("MONGODB_URI", "mongodb://localhost:27017")
	cfg.MongoDB.NAME = getEnv("MONGODB_NAME", "Test")

	// Arduino API key
	cfg.ARDUNIO_API_KEY = getEnv("ARDUNIO_API_KEY", "1234")

	// AWS configuration
	cfg.AWS.Region = getEnv("AWS_REGION", "us-east-1")
	cfg.AWS.AccessKeyID = getEnv("AWS_ACCESS_KEY_ID", "")
	cfg.AWS.SecretAccessKey = getEnv("AWS_SECRET_ACCESS_KEY", "")

	// S3 bucket configuration
	cfg.S3BucketName = getEnv("S3_BUCKET_NAME", "parkme-uploads")

	return cfg
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	// Set default values
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("MONGODB_URI", "mongodb://localhost:27017")
	viper.SetDefault("DATABASE_NAME", "parkme")
	viper.SetDefault("JWT_SECRET", "your-secret-key")
	viper.SetDefault("JWT_EXPIRY", "24h")
	viper.SetDefault("API_KEY", "your-api-key")
	viper.SetDefault("S3_BUCKET_NAME", "parkme-uploads")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %v", err)
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %v", err)
	}

	return &config, nil
}
