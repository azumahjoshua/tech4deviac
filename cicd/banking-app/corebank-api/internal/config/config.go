package config

import (
	"log"
	"os"

	// "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/joho/godotenv"
)

type AWSConfig struct {
	Region          string
	AccessKeyID     string
	SecretAccessKey string
}

func LoadConfig() AWSConfig {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	return AWSConfig{
		Region:          getEnv("AWS_REGION", "us-east-1"),
		AccessKeyID:     getEnv("AWS_ACCESS_KEY_ID", ""),
		SecretAccessKey: getEnv("AWS_SECRET_ACCESS_KEY", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}