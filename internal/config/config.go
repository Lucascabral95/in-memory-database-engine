package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	Port        string
	Environment string
	JWTSecret   string

	RedisTCPEnabled bool
	RedisTCPPort    string
}

func LoadConfig() *Config {
	_ = godotenv.Load()

	cfg := &Config{
		DatabaseURL: getEnvOrDefault("DATABASE_URL", "postgresql://neondb_owner:npg_1MVEsxYtiQ9C@ep-icy-leaf-aiy805im-pooler.c-4.us-east-1.aws.neon.tech/neondb?sslmode=require&channel_binding=require"),
		Port:        getEnvOrDefault("PORT", "8080"),
		Environment: getEnvOrDefault("ENV", "development"),
		JWTSecret:   getEnvOrDefault("JWT_SECRET", "sd,fsdnlfksdmlkf"),

		RedisTCPEnabled: getEnvBoolOrDefault("REDIS_TCP_ENABLED", false),
		RedisTCPPort:    getEnvOrDefault("REDIS_TCP_PORT", "6379"),
	}

	if err := cfg.Validate(); err != nil {
		log.Fatalf("Configuration validation failed: %v", err)
	}

	return cfg
}

func (c *Config) Validate() error {
	if c.DatabaseURL == "" {
		return fmt.Errorf("DATABASE_URL is required but not set. Set it in .env file or as environment variable")
	}

	if c.RedisTCPEnabled {
		if c.RedisTCPPort == "" {
			return fmt.Errorf("REDIS_TCP_PORT is required when REDIS_TCP_ENABLED=true")
		}

		if _, err := strconv.Atoi(strings.TrimPrefix(c.RedisTCPPort, ":")); err != nil {
			return fmt.Errorf("REDIS_TCP_PORT must be numeric: %w", err)
		}
	}

	return nil
}

func (c *Config) GetDSN() string {
	return c.DatabaseURL
}

func (c *Config) GetServerAddr() string {
	return ":" + c.Port
}

func (c *Config) GetRedisTCPAddr() string {
	if strings.HasPrefix(c.RedisTCPPort, ":") {
		return c.RedisTCPPort
	}
	return ":" + c.RedisTCPPort
}

func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvBoolOrDefault(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	parsed, err := strconv.ParseBool(value)
	if err != nil {
		log.Printf("Warning: invalid boolean for %s (%s), using default %t", key, value, defaultValue)
		return defaultValue
	}

	return parsed
}
