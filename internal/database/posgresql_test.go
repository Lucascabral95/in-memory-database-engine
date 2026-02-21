package database

import (
	"context"
	"os"
	"testing"

	"github.com/lucas-dev/in-memory-db/internal/config"
)

func newDatabaseTestConfig(t *testing.T) *config.Config {
	t.Helper()

	if os.Getenv("RUN_DB_TESTS") != "1" {
		t.Skip("integration DB tests disabled. Set RUN_DB_TESTS=1 to run database tests")
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		t.Skip("DATABASE_URL is empty; skipping database integration tests")
	}

	return &config.Config{
		DatabaseURL: dsn,
		Environment: "test",
	}
}

func TestInitDB_Success(t *testing.T) {
	cfg := newDatabaseTestConfig(t)

	db, err := InitDB(context.Background(), cfg)
	if err != nil {
		t.Fatalf("InitDB() error = %v", err)
	}
	if db == nil {
		t.Fatalf("InitDB() returned nil db")
	}

	t.Cleanup(func() {
		_ = CloseDB(db)
	})
}

func TestInitDB_InvalidDSN_ReturnsError(t *testing.T) {
	cfg := &config.Config{
		DatabaseURL: "postgresql://user:pass@invalid-host-xyz:5432/db?sslmode=disable",
		Environment: "test",
	}

	db, err := InitDB(context.Background(), cfg)
	if err == nil {
		t.Fatalf("InitDB() error = nil, want non-nil")
	}
	if db != nil {
		t.Fatalf("InitDB() db = %v, want nil on error", db)
	}
}

func TestGetLogger_ReturnsInterfaceForEachEnv(t *testing.T) {
	devCfg := &config.Config{Environment: "development"}
	prodCfg := &config.Config{Environment: "production"}

	if getLogger(devCfg) == nil {
		t.Fatalf("getLogger(development) returned nil")
	}
	if getLogger(prodCfg) == nil {
		t.Fatalf("getLogger(production) returned nil")
	}
}

func TestCloseDB_Success(t *testing.T) {
	cfg := newDatabaseTestConfig(t)

	db, err := InitDB(context.Background(), cfg)
	if err != nil {
		t.Fatalf("InitDB() error = %v", err)
	}

	if err := CloseDB(db); err != nil {
		t.Fatalf("CloseDB() error = %v", err)
	}
}
