package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DB  *DBConfig
	AWS *AWSConfig
}

func LoadConfig() *Config {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}
	c := &Config{}
	c.DB = &DBConfig{
		Host:     getEnvWithDefault("DB_HOST", "localhost"),
		Port:     getEnvWithDefault("DB_PORT", "5432"),
		User:     getEnvWithDefault("DB_USER", "postgres"),
		Password: getEnvWithDefault("DB_PASSWORD", "postgres"),
		DBName:   getEnvWithDefault("DB_NAME", "chores"),
	}
	c.AWS = &AWSConfig{
		SNSTopicARN: getEnvWithDefault("AWS_SNS_TOPIC_ARN", ""),
	}
	return c
}

type AWSConfig struct {
	SNSTopicARN string
}

// DBConfig stores database connection parameters
type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

// GetDBConnectionString returns the PostgreSQL connection string
func (c *DBConfig) GetDBConnectionString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		c.User, c.Password, c.Host, c.Port, c.DBName)
}

// getEnvWithDefault gets an environment variable or returns a default value
func getEnvWithDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
