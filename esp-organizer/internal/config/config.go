package config

import (
	"log"
	"net/http"
	"os"

	"esp-organizer/internal/utils"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v2"
)

var isLoaded = false

// Config holds all configuration for the application
type Config struct {
	ServerPort     string `yaml:"server_port"`
	DBHost         string `yaml:"db_host"`
	DBPort         string `yaml:"db_port"`
	DBUser         string `yaml:"db_user"`
	DBPassword     string `yaml:"db_password"`
	DBName         string `yaml:"db_name"`
	JWTSecret      string `yaml:"jwt_secret"`
	LogLevel       string `yaml:"log_level"`
	MongoDBURI     string `yaml:"mongodb_uri"`
	PineconeAPIKey string `yaml:"pinecone_api_key"`
	LLMModel       string `yaml:"llm_model"`
	WeaviateURL    string `yaml:"weaviate_url"`
	AWSRegion      string `yaml:"aws_region"`
	AWSS3Bucket    string `yaml:"aws_s3_bucket"`
}

// LoadConfig loads configuration from environment variables and a YAML file
func LoadConfig(envFile, yamlFile string) (*Config, error) {
	// Load .env file if provided
	if envFile != "" {
		if err := godotenv.Load(envFile); err != nil {
			return nil, err
		}
	}

	config := &Config{
		ServerPort:  getEnv("SERVER_PORT", "8080"),
		DBHost:      getEnv("DB_HOST", "localhost"),
		DBPort:      getEnv("DB_PORT", "5432"),
		DBUser:      getEnv("DB_USER", ""),
		DBPassword:  getEnv("DB_PASSWORD", ""),
		DBName:      getEnv("DB_NAME", ""),
		JWTSecret:   getEnv("JWT_SECRET", ""),
		LogLevel:    getEnv("LOG_LEVEL", "info"),
		MongoDBURI:  getEnv("MONGODB_URI", ""),
		WeaviateURL: getEnv("WEAVIATE_URL", "http://localhost:8081"),
		AWSRegion:   getEnv("AWS_REGION", "eu-west-2"),
		AWSS3Bucket: getEnv("AWS_S3_BUCKET", "esp-new-organizer-immunology"),
	}

	if yamlFile != "" {
		file, err := os.Open(yamlFile)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		decoder := yaml.NewDecoder(file)
		if err := decoder.Decode(config); err != nil {
			return nil, err
		}
	}

	return config, nil
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}

// InitConfig initializes the configuration from files
func InitConfig() *Config {
	config, err := LoadConfig(".env", "config/config.yaml")
	if err != nil {
		log.Printf("Warning: .env file not found: %v", config)
		//define the empty error variable
		config, err = LoadConfig("../../.env", "../../config/config.yaml")
		if err != nil {
			log.Fatalf("Error loading config: %v", err)
		}
	}
	return config
}

func GetStringValue(properties map[string]interface{}, key string) string {
	if val, ok := properties[key].(string); ok {
		return val
	}
	return "Unknown"
}

func GetIntValue(properties map[string]interface{}, key string) int {
	if val, ok := properties[key].(float64); ok {
		return int(val)
	}
	if val, ok := properties[key].(int); ok {
		return val
	}
	return 0
}

func GetCertaintyValue(properties map[string]interface{}) float64 {
	if additional, ok := properties["_additional"].(map[string]interface{}); ok {
		if certainty, ok := additional["certainty"].(float64); ok {
			return certainty
		}
	}
	return 0.0
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// CorsMiddleware delegates to the centralized implementation
func CorsMiddleware(next http.Handler) http.Handler {
	return utils.CORSMiddleware(next)
}

func TruncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// LoadConfig finds the project root, loads the .env file, and expands variables.
// It ensures this operation is only performed once.
