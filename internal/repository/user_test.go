package repository

import (
	"testing"

	"github.com/google/uuid"
	"github.com/lucas-dev/in-memory-db/internal/model"
)

func TestUserRepository_RegisterFindUpdateDelete(t *testing.T) {
	tx := newRepositoryTestTx(t)
	repo := NewUserRepository(tx)

	user := &model.User{
		Email:     "repo-user-" + uuid.NewString() + "@example.com",
		Password:  "hashed-password",
		FirstName: "Repo",
		LastName:  "Tester",
	}

	if err := repo.RegisterUser(user); err != nil {
		t.Fatalf("RegisterUser() error = %v", err)
	}

	gotByEmail, err := repo.FindUserByEmail(user.Email)
	if err != nil {
		t.Fatalf("FindUserByEmail() error = %v", err)
	}
	if gotByEmail.ID != user.ID {
		t.Fatalf("FindUserByEmail() ID = %s, want %s", gotByEmail.ID, user.ID)
	}

	user.FirstName = "Updated"
	if err := repo.UpdateUser(user); err != nil {
		t.Fatalf("UpdateUser() error = %v", err)
	}

	updated, err := repo.FindUserByEmail(user.Email)
	if err != nil {
		t.Fatalf("FindUserByEmail() after update error = %v", err)
	}
	if updated.FirstName != "Updated" {
		t.Fatalf("updated FirstName = %s, want %s", updated.FirstName, "Updated")
	}

	if err := repo.DeleteUser(user.Email); err != nil {
		t.Fatalf("DeleteUser() error = %v", err)
	}

	if _, err := repo.FindUserByEmail(user.Email); err == nil {
		t.Fatalf("FindUserByEmail() after delete expected error, got nil")
	}
}

func TestUserRepository_FindUserByID_WithUUIDPrimaryKeyReturnsErrorForUintID(t *testing.T) {
	tx := newRepositoryTestTx(t)
	repo := NewUserRepository(tx)

	user := &model.User{
		Email:     "repo-user-id-" + uuid.NewString() + "@example.com",
		Password:  "hashed-password",
		FirstName: "ID",
		LastName:  "Mismatch",
	}
	if err := repo.RegisterUser(user); err != nil {
		t.Fatalf("RegisterUser() error = %v", err)
	}

	if _, err := repo.FindUserByID(1); err == nil {
		t.Fatalf("FindUserByID() expected error because model uses UUID primary key")
	}
}
