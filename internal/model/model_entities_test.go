package model

import (
	"reflect"
	"testing"

	"github.com/google/uuid"
)

func TestOrderStatusConstants(t *testing.T) {
	if OrderStatusPending != "PENDING" {
		t.Fatalf("OrderStatusPending = %s, want %s", OrderStatusPending, "PENDING")
	}
	if OrderStatusPaid != "PAID" {
		t.Fatalf("OrderStatusPaid = %s, want %s", OrderStatusPaid, "PAID")
	}
	if OrderStatusShipped != "SHIPPED" {
		t.Fatalf("OrderStatusShipped = %s, want %s", OrderStatusShipped, "SHIPPED")
	}
	if OrderStatusCancelled != "CANCELLED" {
		t.Fatalf("OrderStatusCancelled = %s, want %s", OrderStatusCancelled, "CANCELLED")
	}
}

func TestStockMovementReasonConstants(t *testing.T) {
	if StockMovementReasonSale != "SALE" {
		t.Fatalf("StockMovementReasonSale = %s, want %s", StockMovementReasonSale, "SALE")
	}
	if StockMovementReasonRestock != "RESTOCK" {
		t.Fatalf("StockMovementReasonRestock = %s, want %s", StockMovementReasonRestock, "RESTOCK")
	}
	if StockMovementReasonAdjustment != "ADJUSTMENT" {
		t.Fatalf("StockMovementReasonAdjustment = %s, want %s", StockMovementReasonAdjustment, "ADJUSTMENT")
	}
}

func TestStockMovement_BeforeCreate_AssignsUUIDWhenNil(t *testing.T) {
	sm := &StockMovement{}

	if err := sm.BeforeCreate(nil); err != nil {
		t.Fatalf("BeforeCreate() error = %v, want nil", err)
	}
	if sm.ID == uuid.Nil {
		t.Fatalf("BeforeCreate() ID = nil UUID, want generated UUID")
	}
}

func TestProductModelTags(t *testing.T) {
	typ := reflect.TypeOf(Product{})

	nameField, ok := typ.FieldByName("Name")
	if !ok {
		t.Fatalf("Product.Name field not found")
	}
	if got := nameField.Tag.Get("json"); got != "name" {
		t.Fatalf("Product.Name json tag = %s, want %s", got, "name")
	}

	stockField, ok := typ.FieldByName("Stock")
	if !ok {
		t.Fatalf("Product.Stock field not found")
	}
	if got := stockField.Tag.Get("json"); got != "stock" {
		t.Fatalf("Product.Stock json tag = %s, want %s", got, "stock")
	}
}

func TestCartModel_DefaultZeroValue(t *testing.T) {
	var cart Cart

	if cart.UserID != uuid.Nil {
		t.Fatalf("zero Cart.UserID = %s, want nil UUID", cart.UserID)
	}
	if cart.Items != nil {
		t.Fatalf("zero Cart.Items = %#v, want nil slice", cart.Items)
	}
}

func TestUpdateOrderStatusRequestTag(t *testing.T) {
	typ := reflect.TypeOf(UpdateOrderStatusRequest{})
	field, ok := typ.FieldByName("Status")
	if !ok {
		t.Fatalf("UpdateOrderStatusRequest.Status field not found")
	}

	if got := field.Tag.Get("binding"); got == "" {
		t.Fatalf("binding tag should not be empty")
	}
}
