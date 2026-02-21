package config

import (
	"strings"
	"testing"
	"time"
)

func TestConfigValidate_RequiresDatabaseURL(t *testing.T) {
	cfg := &Config{
		DatabaseURL: "",
	}

	err := cfg.Validate()
	if err == nil {
		t.Fatalf("Validate() error = nil, want non-nil")
	}
}

func TestConfigValidate_RedisEnabledRequiresPort(t *testing.T) {
	cfg := &Config{
		DatabaseURL:     "postgres://example",
		RedisTCPEnabled: true,
		RedisTCPPort:    "",
	}

	err := cfg.Validate()
	if err == nil {
		t.Fatalf("Validate() error = nil, want non-nil")
	}
}

func TestConfigValidate_RedisEnabledRequiresNumericPort(t *testing.T) {
	cfg := &Config{
		DatabaseURL:     "postgres://example",
		RedisTCPEnabled: true,
		RedisTCPPort:    "abc",
	}

	err := cfg.Validate()
	if err == nil {
		t.Fatalf("Validate() error = nil, want non-nil")
	}
	if !strings.Contains(err.Error(), "REDIS_TCP_PORT must be numeric") {
		t.Fatalf("Validate() error = %v, want numeric port message", err)
	}
}

func TestConfigValidate_Success(t *testing.T) {
	cfg := &Config{
		DatabaseURL:     "postgres://example",
		RedisTCPEnabled: true,
		RedisTCPPort:    "6379",
	}

	if err := cfg.Validate(); err != nil {
		t.Fatalf("Validate() error = %v, want nil", err)
	}
}

func TestConfigHelpers(t *testing.T) {
	cfg := &Config{
		DatabaseURL:  "postgres://example",
		Port:         "8080",
		Environment:  "development",
		RedisTCPPort: "6379",
	}

	if got := cfg.GetDSN(); got != "postgres://example" {
		t.Fatalf("GetDSN() = %s, want %s", got, "postgres://example")
	}

	if got := cfg.GetServerAddr(); got != ":8080" {
		t.Fatalf("GetServerAddr() = %s, want %s", got, ":8080")
	}

	if got := cfg.GetRedisTCPAddr(); got != ":6379" {
		t.Fatalf("GetRedisTCPAddr() = %s, want %s", got, ":6379")
	}

	cfg.RedisTCPPort = ":6380"
	if got := cfg.GetRedisTCPAddr(); got != ":6380" {
		t.Fatalf("GetRedisTCPAddr() = %s, want %s", got, ":6380")
	}

	if !cfg.IsDevelopment() {
		t.Fatalf("IsDevelopment() = false, want true")
	}
}

func TestGetEnvOrDefault(t *testing.T) {
	key := "CONFIG_TEST_ENV_FOO"
	t.Setenv(key, "")

	if got := getEnvOrDefault(key, "default-value"); got != "default-value" {
		t.Fatalf("getEnvOrDefault() = %s, want %s", got, "default-value")
	}

	t.Setenv(key, "env-value")
	if got := getEnvOrDefault(key, "default-value"); got != "env-value" {
		t.Fatalf("getEnvOrDefault() = %s, want %s", got, "env-value")
	}
}

func TestGetEnvBoolOrDefault(t *testing.T) {
	key := "CONFIG_TEST_ENV_BOOL"

	t.Setenv(key, "")
	if got := getEnvBoolOrDefault(key, true); !got {
		t.Fatalf("getEnvBoolOrDefault() = false, want true")
	}

	t.Setenv(key, "false")
	if got := getEnvBoolOrDefault(key, true); got {
		t.Fatalf("getEnvBoolOrDefault() = true, want false")
	}

	t.Setenv(key, "invalid")
	if got := getEnvBoolOrDefault(key, true); !got {
		t.Fatalf("getEnvBoolOrDefault() invalid value should return default")
	}
}

func TestGetEnvIntOrDefault(t *testing.T) {
	key := "CONFIG_TEST_ENV_INT"

	t.Setenv(key, "")
	if got := getEnvIntOrDefault(key, 20); got != 20 {
		t.Fatalf("getEnvIntOrDefault() = %d, want %d", got, 20)
	}

	t.Setenv(key, "33")
	if got := getEnvIntOrDefault(key, 20); got != 33 {
		t.Fatalf("getEnvIntOrDefault() = %d, want %d", got, 33)
	}

	t.Setenv(key, "invalid")
	if got := getEnvIntOrDefault(key, 20); got != 20 {
		t.Fatalf("getEnvIntOrDefault() invalid value should return default")
	}
}

func TestGetEnvDurationOrDefault(t *testing.T) {
	key := "CONFIG_TEST_ENV_DURATION"

	t.Setenv(key, "")
	if got := getEnvDurationOrDefault(key, 5*time.Minute); got != 5*time.Minute {
		t.Fatalf("getEnvDurationOrDefault() = %v, want %v", got, 5*time.Minute)
	}

	t.Setenv(key, "30s")
	if got := getEnvDurationOrDefault(key, 5*time.Minute); got != 30*time.Second {
		t.Fatalf("getEnvDurationOrDefault() = %v, want %v", got, 30*time.Second)
	}

	t.Setenv(key, "invalid")
	if got := getEnvDurationOrDefault(key, 5*time.Minute); got != 5*time.Minute {
		t.Fatalf("getEnvDurationOrDefault() invalid value should return default")
	}
}
