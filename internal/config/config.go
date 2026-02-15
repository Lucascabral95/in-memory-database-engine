package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	Port        string
	Environment string
	JWTSecret   string
}

func LoadConfig() *Config {
	_ = godotenv.Load()

	cfg := &Config{
		DatabaseURL: getEnvOrDefault("DATABASE_URL", "postgresql://neondb_owner:npg_1MVEsxYtiQ9C@ep-icy-leaf-aiy805im-pooler.c-4.us-east-1.aws.neon.tech/neondb?sslmode=require&channel_binding=require"),
		Port:        getEnvOrDefault("PORT", "8080"),
		Environment: getEnvOrDefault("ENV", "development"),
		JWTSecret:   getEnvOrDefault("JWT_SECRET", "sd,fsdnlfksdmlkf"),
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
	return nil
}

func (c *Config) GetDSN() string {
	return c.DatabaseURL
}

func (c *Config) GetServerAddr() string {
	return ":" + c.Port
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
