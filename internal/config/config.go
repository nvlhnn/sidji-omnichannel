package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	App      AppConfig
	Database DatabaseConfig
	Redis    RedisConfig
	Meta     MetaConfig
	AWS      AWSConfig
	AI       AIConfig
}


type AppConfig struct {
	Env         string
	Port        int
	Secret      string
	FrontendURL string
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
}

type MetaConfig struct {
	AppID                     string
	AppSecret                 string
	VerifyToken               string
}

type AWSConfig struct {
	Region          string
	AccessKeyID     string
	SecretAccessKey string
	S3Bucket        string
}

type AIConfig struct {
	GeminiAPIKey string
	OpenAIAPIKey string
	Provider     string // "gemini" or "openai"
}

func Load() (*Config, error) {
	// Load .env file if exists
	godotenv.Load()

	return &Config{
		App: AppConfig{
			Env:         getEnv("APP_ENV", "development"),
			Port:        getEnvInt("APP_PORT", 8080),
			Secret:      getEnv("APP_SECRET", "change-this-secret"),
			FrontendURL: getEnv("FRONTEND_URL", "http://localhost:3000"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvInt("DB_PORT", 5432),
			User:     getEnv("DB_USER", "sidji"),
			Password: getEnv("DB_PASSWORD", "sidji123"),
			Name:     getEnv("DB_NAME", "sidji"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnvInt("REDIS_PORT", 6379),
			Password: getEnv("REDIS_PASSWORD", ""),
		},
		Meta: MetaConfig{
			AppID:                     getEnv("META_APP_ID", ""),
			AppSecret:                 getEnv("META_APP_SECRET", ""),
			VerifyToken:               getEnv("META_VERIFY_TOKEN", ""),
		},
		AWS: AWSConfig{
			Region:          getEnv("AWS_REGION", "ap-southeast-1"),
			AccessKeyID:     getEnv("AWS_ACCESS_KEY_ID", ""),
			SecretAccessKey: getEnv("AWS_SECRET_ACCESS_KEY", ""),
			S3Bucket:        getEnv("AWS_S3_BUCKET", "sidji-omnichannel-media"),
		},
		AI: AIConfig{
			GeminiAPIKey: getEnv("GEMINI_API_KEY", ""),
			OpenAIAPIKey: getEnv("OPENAI_API_KEY", ""),
			Provider:     getEnv("AI_PROVIDER", "gemini"),
		},
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}
