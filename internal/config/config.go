package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost            string
	DBPort            string
	DBUser            string
	DBPassword        string
	DBName            string
	DBSSLMode         string
	ServerHost        string
	ServerPort        string
	LogLevel          string
	LogFormat         string
	GenderizeAPIURL   string
	AgifyAPIURL       string
	NationalizeAPIURL string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, loading from environment variables")
	}

	return &Config{
		DBHost:            getEnv("DB_HOST", "localhost"),
		DBPort:            getEnv("DB_PORT", "5432"),
		DBUser:            getEnv("DB_USER", "postgres"),
		DBPassword:        getEnv("DB_PASSWORD", ""),
		DBName:            getEnv("DB_NAME", "effective_mobile_db"),
		DBSSLMode:         getEnv("DB_SSLMODE", "disable"),
		ServerHost:        getEnv("SERVER_HOST", "0.0.0.0"),
		ServerPort:        getEnv("SERVER_PORT", "8080"),
		LogLevel:          getEnv("LOG_LEVEL", "info"),
		LogFormat:         getEnv("LOG_FORMAT", "text"),
		GenderizeAPIURL:   getEnv("GENDERIZE_API_URL", "https://api.genderize.io"),
		AgifyAPIURL:       getEnv("AGIFY_API_URL", "https://api.agify.io"),
		NationalizeAPIURL: getEnv("NATIONALIZE_API_URL", "https://api.nationalize.io"),
	}
}

func getEnv(key, defaultVal string) string {
	if val, exists := os.LookupEnv(key); exists {
		return val
	}
	return defaultVal
}
