package database

import (
	"fmt"
	"log"

	"github.com/lucas-dev/in-memory-db/internal/config"
	"github.com/lucas-dev/in-memory-db/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB(cfg *config.Config) *gorm.DB {
	log.Println("=> Initializing database connection...")

	db, err := gorm.Open(postgres.Open(cfg.GetDSN()), &gorm.Config{
		PrepareStmt: false,
		Logger:      getLogger(cfg),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("=> Database connection established")

	if err := runMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	log.Println("=> Database migrations completed")

	return db
}

func runMigrations(db *gorm.DB) error {
	log.Println("=> Running database migrations...")

	return db.AutoMigrate(
		&model.User{},
		&model.Category{},
		&model.Product{},
		&model.Order{},
		&model.OrderItem{},
		&model.StockMovement{},
	)
}

func getLogger(cfg *config.Config) logger.Interface {
	if cfg.IsDevelopment() {
		return logger.Default.LogMode(logger.Info)
	}
	return logger.Default.LogMode(logger.Silent)
}

func CloseDB(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}
	return sqlDB.Close()
}
