package model

import (
	"reflect"
	"testing"

	"github.com/google/uuid"
)

func TestBaseModel_BeforeCreate_AssignsUUIDWhenNil(t *testing.T) {
	base := &BaseModel{}

	if err := base.BeforeCreate(nil); err != nil {
		t.Fatalf("BeforeCreate() error = %v, want nil", err)
	}

	if base.ID == uuid.Nil {
		t.Fatalf("BeforeCreate() ID = nil UUID, want generated UUID")
	}
}

func TestBaseModel_BeforeCreate_KeepsExistingUUID(t *testing.T) {
	existing := uuid.New()
	base := &BaseModel{ID: existing}

	if err := base.BeforeCreate(nil); err != nil {
		t.Fatalf("BeforeCreate() error = %v, want nil", err)
	}

	if base.ID != existing {
		t.Fatalf("BeforeCreate() changed existing ID = %s, want %s", base.ID, existing)
	}
}

func TestUserModelTags(t *testing.T) {
	typ := reflect.TypeOf(User{})

	emailField, ok := typ.FieldByName("Email")
	if !ok {
		t.Fatalf("User.Email field not found")
	}
	if got := emailField.Tag.Get("json"); got != "email" {
		t.Fatalf("User.Email json tag = %s, want %s", got, "email")
	}

	passwordField, ok := typ.FieldByName("Password")
	if !ok {
		t.Fatalf("User.Password field not found")
	}
	if got := passwordField.Tag.Get("json"); got != "-" {
		t.Fatalf("User.Password json tag = %s, want %s", got, "-")
	}
}
