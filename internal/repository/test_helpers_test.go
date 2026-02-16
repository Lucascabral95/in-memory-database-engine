package repository

import (
	"os"
	"testing"

	"github.com/lucas-dev/in-memory-db/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type repositoryTestDB = *gorm.DB

func newRepositoryTestTx(t *testing.T) repositoryTestDB {
	t.Helper()

	if os.Getenv("RUN_DB_TESTS") != "1" {
		t.Skip("integration DB tests disabled. Set RUN_DB_TESTS=1 to run repository tests")
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		t.Skip("DATABASE_URL is empty; skipping repository integration tests")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("open test DB error: %v", err)
	}

	if err := db.AutoMigrate(
		&model.User{},
		&model.Category{},
		&model.Product{},
		&model.Order{},
		&model.OrderItem{},
		&model.StockMovement{},
	); err != nil {
		t.Fatalf("auto migrate test DB error: %v", err)
	}

	tx := db.Begin()
	if tx.Error != nil {
		t.Fatalf("begin tx error: %v", tx.Error)
	}

	t.Cleanup(func() {
		_ = tx.Rollback().Error
	})

	return tx
}
