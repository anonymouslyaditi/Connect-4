package main

import (
	"os"
	"strconv"
)

// Config holds all configuration for the application
type Config struct {
	ServerPort      string
	DataDir         string
	KafkaEnabled    bool
	KafkaBrokers    []string
	KafkaTopic      string
	DBEnabled       bool
	DBHost          string
	DBPort          int
	DBUser          string
	DBPassword      string
	DBName          string
	MatchTimeout    int // seconds to wait for matchmaking
	ReconnectTimeout int // seconds to allow reconnection
}

// LoadConfig loads configuration from environment variables with defaults
func LoadConfig() *Config {
	return &Config{
		ServerPort:       getEnv("SERVER_PORT", "8080"),
		DataDir:          getEnv("DATA_DIR", "data"),
		KafkaEnabled:     getEnvBool("KAFKA_ENABLED", false),
		KafkaBrokers:     getEnvSlice("KAFKA_BROKERS", []string{"localhost:9092"}),
		KafkaTopic:       getEnv("KAFKA_TOPIC", "game-analytics"),
		DBEnabled:        getEnvBool("DB_ENABLED", false),
		DBHost:           getEnv("DB_HOST", "localhost"),
		DBPort:           getEnvInt("DB_PORT", 5432),
		DBUser:           getEnv("DB_USER", "postgres"),
		DBPassword:       getEnv("DB_PASSWORD", "postgres"),
		DBName:           getEnv("DB_NAME", "connect4"),
		MatchTimeout:     getEnvInt("MATCH_TIMEOUT", 10),
		ReconnectTimeout: getEnvInt("RECONNECT_TIMEOUT", 30),
	}
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
	if val := os.Getenv(key); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return defaultVal
}

func getEnvBool(key string, defaultVal bool) bool {
	if val := os.Getenv(key); val != "" {
		if b, err := strconv.ParseBool(val); err == nil {
			return b
		}
	}
	return defaultVal
}

func getEnvSlice(key string, defaultVal []string) []string {
	if val := os.Getenv(key); val != "" {
		// Simple comma-separated parsing
		result := []string{}
		current := ""
		for _, c := range val {
			if c == ',' {
				if current != "" {
					result = append(result, current)
					current = ""
				}
			} else {
				current += string(c)
			}
		}
		if current != "" {
			result = append(result, current)
		}
		return result
	}
	return defaultVal
}

