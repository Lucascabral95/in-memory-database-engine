package database

import (
	"os"
	"testing"

	"github.com/lucas-dev/in-memory-db/internal/config"
	"github.com/lucas-dev/in-memory-db/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func newDatabaseTestConnection(t *testing.T) *gorm.DB {
	t.Helper()

	if os.Getenv("RUN_DB_TESTS") != "1" {
		t.Skip("integration DB tests disabled. Set RUN_DB_TESTS=1 to run database tests")
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		t.Skip("DATABASE_URL is empty; skipping database integration tests")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("open test DB error: %v", err)
	}

	return db
}

func TestRunMigrations_CreatesTables(t *testing.T) {
	db := newDatabaseTestConnection(t)

	if err := runMigrations(db); err != nil {
		t.Fatalf("runMigrations() error = %v", err)
	}

	entities := []interface{}{
		&model.User{},
		&model.Category{},
		&model.Product{},
		&model.Order{},
		&model.OrderItem{},
		&model.StockMovement{},
	}

	for _, entity := range entities {
		if !db.Migrator().HasTable(entity) {
			t.Fatalf("expected table for %T to exist after runMigrations()", entity)
		}
	}
}

func TestCloseDB_Success(t *testing.T) {
	db := newDatabaseTestConnection(t)

	if err := CloseDB(db); err != nil {
		t.Fatalf("CloseDB() error = %v", err)
	}
}

func TestGetLogger_ReturnsInterfaceForEachEnv(t *testing.T) {
	devCfg := &config.Config{Environment: "development"}
	prodCfg := &config.Config{Environment: "production"}

	devLogger := getLogger(devCfg)
	if devLogger == nil {
		t.Fatalf("getLogger(development) returned nil")
	}

	prodLogger := getLogger(prodCfg)
	if prodLogger == nil {
		t.Fatalf("getLogger(production) returned nil")
	}
}
