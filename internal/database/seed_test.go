package database

import (
	"os"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/lucas-dev/in-memory-db/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type databaseTestDB = *gorm.DB

func newDatabaseTestTx(t *testing.T) databaseTestDB {
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

func TestSeedDatabase_CreatesRequestedVolumes(t *testing.T) {
	tx := newDatabaseTestTx(t)

	opts := SeedOptions{
		Users:      22,
		Categories: 7,
		Products:   55,
		Orders:     41,
	}

	if err := SeedDatabase(tx, opts); err != nil {
		t.Fatalf("SeedDatabase() error = %v", err)
	}

	assertCount(t, tx, &model.User{}, 22)
	assertCount(t, tx, &model.Category{}, 7)
	assertCount(t, tx, &model.Product{}, 55)
	assertCount(t, tx, &model.Order{}, 41)

	var itemCount int64
	if err := tx.Model(&model.OrderItem{}).Count(&itemCount).Error; err != nil {
		t.Fatalf("count order items error: %v", err)
	}
	if itemCount < 41 {
		t.Fatalf("order_items = %d, want >= %d", itemCount, 41)
	}
}

func TestSeedDatabase_WithoutForceFailsOnExistingData(t *testing.T) {
	tx := newDatabaseTestTx(t)

	existing := model.User{
		Email:     "already-there@example.com",
		Password:  "hashed-password",
		FirstName: "Already",
		LastName:  "There",
	}
	if err := tx.Create(&existing).Error; err != nil {
		t.Fatalf("setup existing user error: %v", err)
	}

	err := SeedDatabase(tx, SeedOptions{})
	if err == nil {
		t.Fatalf("SeedDatabase() error = nil, want non-nil")
	}
	if !strings.Contains(err.Error(), "seed cancelled: table users already has data") {
		t.Fatalf("SeedDatabase() error = %q, want users table collision", err.Error())
	}
}

func TestSeedDatabase_WithForcePurgesAndSeedsAgain(t *testing.T) {
	tx := newDatabaseTestTx(t)

	oldUser := model.User{
		Email:     "old-user@example.com",
		Password:  "hashed-password",
		FirstName: "Old",
		LastName:  "User",
	}
	if err := tx.Create(&oldUser).Error; err != nil {
		t.Fatalf("create old user error: %v", err)
	}

	oldCategory := model.Category{Name: "Old Category"}
	if err := tx.Create(&oldCategory).Error; err != nil {
		t.Fatalf("create old category error: %v", err)
	}

	oldProduct := model.Product{
		Name:        "Old Product",
		SKU:         "OLD-" + uuid.NewString(),
		Description: "old",
		Price:       10,
		Stock:       10,
		CategoryID:  &oldCategory.ID,
	}
	if err := tx.Create(&oldProduct).Error; err != nil {
		t.Fatalf("create old product error: %v", err)
	}

	oldOrder := model.Order{
		UserID:      oldUser.ID,
		TotalAmount: 10,
		Status:      model.OrderStatusPending,
	}
	if err := tx.Create(&oldOrder).Error; err != nil {
		t.Fatalf("create old order error: %v", err)
	}

	oldItem := model.OrderItem{
		OrderID:       oldOrder.ID,
		ProductID:     oldProduct.ID,
		Quantity:      1,
		PriceAtMoment: 10,
	}
	if err := tx.Create(&oldItem).Error; err != nil {
		t.Fatalf("create old order item error: %v", err)
	}

	oldMovement := model.StockMovement{
		ProductID: oldProduct.ID,
		Quantity:  -1,
		Reason:    model.StockMovementReasonSale,
	}
	if err := tx.Create(&oldMovement).Error; err != nil {
		t.Fatalf("create old stock movement error: %v", err)
	}

	if err := SeedDatabase(tx, SeedOptions{
		Users:      1,
		Categories: 1,
		Products:   1,
		Orders:     1,
		Force:      true,
	}); err != nil {
		t.Fatalf("SeedDatabase(force) error = %v", err)
	}

	assertCount(t, tx, &model.User{}, 1)
	assertCount(t, tx, &model.Category{}, 1)
	assertCount(t, tx, &model.Product{}, 1)
	assertCount(t, tx, &model.Order{}, 1)
	assertCount(t, tx, &model.OrderItem{}, 1)

	var oldUserCount int64
	if err := tx.Model(&model.User{}).Where("email = ?", oldUser.Email).Count(&oldUserCount).Error; err != nil {
		t.Fatalf("count old user by email error: %v", err)
	}
	if oldUserCount != 0 {
		t.Fatalf("old seeded user still exists, count = %d", oldUserCount)
	}
}

func TestRound2(t *testing.T) {
	got := round2(12.345)
	want := 12.35
	if got != want {
		t.Fatalf("round2(12.345) = %.2f, want %.2f", got, want)
	}
}

func assertCount(t *testing.T, tx databaseTestDB, entity interface{}, want int64) {
	t.Helper()

	var got int64
	if err := tx.Model(entity).Count(&got).Error; err != nil {
		t.Fatalf("count %T error: %v", entity, err)
	}
	if got != want {
		t.Fatalf("count %T = %d, want %d", entity, got, want)
	}
}
