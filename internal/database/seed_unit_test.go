package database

import (
	"math/rand"
	"testing"

	"github.com/google/uuid"
	"github.com/lucas-dev/in-memory-db/internal/model"
)

func TestSeedOptionsNormalize_UsesDefaults(t *testing.T) {
	got := (SeedOptions{}).normalize()

	if got.Users != defaultSeedUsers {
		t.Fatalf("Users = %d, want %d", got.Users, defaultSeedUsers)
	}
	if got.Categories != defaultSeedCategories {
		t.Fatalf("Categories = %d, want %d", got.Categories, defaultSeedCategories)
	}
	if got.Products != defaultSeedProducts {
		t.Fatalf("Products = %d, want %d", got.Products, defaultSeedProducts)
	}
	if got.Orders != defaultSeedOrders {
		t.Fatalf("Orders = %d, want %d", got.Orders, defaultSeedOrders)
	}
}

func TestSeedOptionsNormalize_RespectsProvidedValues(t *testing.T) {
	in := SeedOptions{
		Users:      30,
		Categories: 6,
		Products:   80,
		Orders:     60,
		Force:      true,
	}
	got := in.normalize()

	if got.Users != in.Users || got.Categories != in.Categories || got.Products != in.Products || got.Orders != in.Orders {
		t.Fatalf("normalize() changed provided values unexpectedly: got %+v want %+v", got, in)
	}
	if !got.Force {
		t.Fatalf("normalize() should preserve Force=true")
	}
}

func TestPickUniqueProducts_ReturnsUniqueSubset(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	products := []model.Product{
		{BaseModel: model.BaseModel{ID: uuid.New()}},
		{BaseModel: model.BaseModel{ID: uuid.New()}},
		{BaseModel: model.BaseModel{ID: uuid.New()}},
		{BaseModel: model.BaseModel{ID: uuid.New()}},
	}

	selected := pickUniqueProducts(products, 3, rng)
	if len(selected) != 3 {
		t.Fatalf("len(selected) = %d, want %d", len(selected), 3)
	}

	seen := make(map[uuid.UUID]struct{}, len(selected))
	for _, p := range selected {
		if _, ok := seen[p.ID]; ok {
			t.Fatalf("duplicate product ID found: %s", p.ID)
		}
		seen[p.ID] = struct{}{}
	}
}

func TestPickUniqueProducts_ReturnsAllWhenRequestedMoreThanAvailable(t *testing.T) {
	rng := rand.New(rand.NewSource(7))
	products := []model.Product{
		{BaseModel: model.BaseModel{ID: uuid.New()}},
		{BaseModel: model.BaseModel{ID: uuid.New()}},
	}

	selected := pickUniqueProducts(products, 10, rng)
	if len(selected) != len(products) {
		t.Fatalf("len(selected) = %d, want %d", len(selected), len(products))
	}
}

func TestCanFulfill(t *testing.T) {
	p1 := uuid.New()
	p2 := uuid.New()

	required := map[uuid.UUID]int{
		p1: 2,
		p2: 1,
	}
	stock := map[uuid.UUID]int{
		p1: 5,
		p2: 1,
	}

	if !canFulfill(required, stock) {
		t.Fatalf("canFulfill() = false, want true")
	}

	stock[p2] = 0
	if canFulfill(required, stock) {
		t.Fatalf("canFulfill() = true, want false")
	}
}

func TestRandomOrderStatus_OnlyReturnsAllowedValues(t *testing.T) {
	rng := rand.New(rand.NewSource(99))
	allowed := map[model.OrderStatus]struct{}{
		model.OrderStatusPending:   {},
		model.OrderStatusPaid:      {},
		model.OrderStatusShipped:   {},
		model.OrderStatusCancelled: {},
	}

	for i := 0; i < 1000; i++ {
		status := randomOrderStatus(rng)
		if _, ok := allowed[status]; !ok {
			t.Fatalf("randomOrderStatus() returned invalid status: %q", status)
		}
	}
}
