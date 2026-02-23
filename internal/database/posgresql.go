package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/lucas-dev/in-memory-db/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	maxRetries  = 5
	pingTimeout = 5 * time.Second
	retryDelay  = 2 * time.Second
)

func InitDB(ctx context.Context, cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.GetDSN()), &gorm.Config{
		PrepareStmt: true,
		Logger:      getLogger(cfg),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	sqlDB.SetMaxOpenConns(cfg.DbMaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.DbMaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.DbConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(cfg.DbConnMaxIdleTime)

	for attempt := 1; attempt <= maxRetries; attempt++ {
		ctxPing, cancel := context.WithTimeout(ctx, pingTimeout)
		err := sqlDB.PingContext(ctxPing)
		cancel()

		if err == nil {
			log.Printf("Database connected successfully on attempt %d", attempt)
			return db, nil
		}

		log.Printf("Ping attempt %d/%d failed: %v", attempt, maxRetries, err)

		if attempt < maxRetries {
			waitOrExit(ctx, retryDelay)
		}
	}

	return nil, fmt.Errorf("failed to connect to database after %d attempts", maxRetries)
}

// waitOrExit espera el delay o cancela si el contexto expira
func waitOrExit(ctx context.Context, delay time.Duration) {
	select {
	case <-ctx.Done():
		log.Fatalf("context cancelled while waiting for retry: %v", ctx.Err())
	case <-time.After(delay):
	}
}

func getLogger(cfg *config.Config) logger.Interface {
	if cfg.IsDevelopment() {
		return logger.Default.LogMode(logger.Info)
	}
	return logger.Default.LogMode(logger.Error)
}

func CloseDB(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}
	return sqlDB.Close()
}