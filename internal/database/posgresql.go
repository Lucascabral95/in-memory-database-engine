package database

import (
	"context"
	"fmt"
	"time"

	"github.com/lucas-dev/in-memory-db/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB(ctx context.Context, cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.GetDSN()), &gorm.Config{
		PrepareStmt: true,
		Logger:      getLogger(cfg),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	sqlDb, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	sqlDb.SetMaxOpenConns(cfg.DbMaxOpenConns)
	sqlDb.SetMaxIdleConns(cfg.DbMaxIdleConns)
	sqlDb.SetConnMaxLifetime(cfg.DbConnMaxLifetime)
	sqlDb.SetConnMaxIdleTime(cfg.DbConnMaxIdleTime)

	ctxPing, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := sqlDb.PingContext(ctxPing); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
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
