package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Env      string
	GRPC     GRPCConfig
	Postgres PostgresConfig
	TokenTTL time.Duration
}

type GRPCConfig struct {
	Port    int
	Timeout time.Duration
}

type PostgresConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")

	if configPath == "" {
		panic("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist: " + configPath)
	}

	if err := godotenv.Load(configPath); err != nil {
		panic("failed to load config: " + err.Error())
	}

	return &Config{
		Env: getStringEnv("ENV", "local"),
		GRPC: GRPCConfig{
			Port:    getIntEnv("GRPC_PORT", 55051),
			Timeout: getDurationEnv("GRPC_TIMEOUT", "10s"),
		},
		Postgres: PostgresConfig{
			Host:     getStringEnv("POSTGRES_HOST", "localhost"),
			Port:     getIntEnv("POSTGRES_PORT", 5432),
			User:     getStringEnv("POSTGRES_USER", ""),
			Password: getStringEnv("POSTGRES_PASSWORD", ""),
			Database: getStringEnv("POSTGRES_DB", ""),
		},
		TokenTTL: getDurationEnv("TOKEN_TTL", "1h"),
	}
}

func getStringEnv(key string, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	if defaultValue == "" {
		panic(key + " is not set")
	}

	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value, ok := os.LookupEnv(key); ok {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue string) time.Duration {
	val := getStringEnv(key, defaultValue)

	duration, err := time.ParseDuration(val)
	if err != nil {
		panic("failed to parse " + key + ": " + err.Error())
	}

	return duration
}
