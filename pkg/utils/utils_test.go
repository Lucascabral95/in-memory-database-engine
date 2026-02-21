package utils

import (
	"strings"
	"testing"

	"github.com/lucas-dev/in-memory-db/internal/model"
	"golang.org/x/crypto/bcrypt"
)

func TestConvertAnyToInt_Success(t *testing.T) {
	input := "12345"
	result, err := ConvertAnyToInt(input)

	if err != nil {
		t.Errorf("Se esperaba error nil, se obtuvo: %v", err)
	}
	if result != 12345 {
		t.Errorf("Se esperaba 12345, se obtuvo %d", result)
	}
}

func TestConvertAnyToInt_Negative(t *testing.T) {
	input := "-500"
	result, err := ConvertAnyToInt(input)

	if err != nil {
		t.Errorf("Se esperaba error nil, se obtuvo: %v", err)
	}
	if result != -500 {
		t.Errorf("Se esperaba -500, se obtuvo %d", result)
	}
}

func TestConvertAnyToInt_InvalidString(t *testing.T) {
	input := "not-a-number"
	_, err := ConvertAnyToInt(input)

	if err == nil {
		t.Error("Se esperaba un error al convertir un string no numérico, pero fue nil")
	}
}

func TestGenerateUUID_Format(t *testing.T) {
	uuid := GenerateUUID()

	if len(uuid) != 36 {
		t.Errorf("Longitud de UUID incorrecta: esperada 36, obtenida %d", len(uuid))
	}

	if uuid[8] != '-' || uuid[13] != '-' || uuid[18] != '-' || uuid[23] != '-' {
		t.Error("El formato del UUID no coincide con el estándar")
	}
}

func TestGenerateUUID_Uniqueness(t *testing.T) {
	id1 := GenerateUUID()
	id2 := GenerateUUID()

	if id1 == id2 {
		t.Error("Se generaron dos UUIDs idénticos, deberían ser únicos")
	}
}

func TestGenerateSKU_Format(t *testing.T) {
	sku := GenerateSKU()

	if !strings.HasPrefix(sku, "SKU-") {
		t.Errorf("El SKU debe empezar con 'SKU-', se obtuvo: %s", sku)
	}

	if len(sku) != 40 {
		t.Errorf("Longitud de SKU incorrecta: esperada 40, obtenida %d", len(sku))
	}
}

func TestIsValidOrderStatus_Valid(t *testing.T) {
	validStatuses := []model.OrderStatus{
		model.OrderStatusPending,
		model.OrderStatusPaid,
		model.OrderStatusShipped,
		model.OrderStatusCancelled,
	}

	for _, status := range validStatuses {
		if !IsValidOrderStatus(status) {
			t.Errorf("El estado %v debería ser válido", status)
		}
	}
}

func TestIsValidOrderStatus_Invalid(t *testing.T) {
	invalidStatus := model.OrderStatus("INVALID_STATUS")

	if IsValidOrderStatus(invalidStatus) {
		t.Error("Un estado inválido debería retornar false")
	}
}

func TestHashPassword_Success(t *testing.T) {
	password := "mySecretPassword123"
	hash, err := HashPassword(password)

	if err != nil {
		t.Fatalf("Error inesperado hasheando password: %v", err)
	}

	if hash == "" {
		t.Error("El hash no debería estar vacío")
	}

	if hash == password {
		t.Error("El hash no debe ser igual al password en texto plano")
	}
}

func TestComparePassword_Correct(t *testing.T) {
	password := "mySecretPassword123"
	hash, _ := HashPassword(password)

	err := ComparePassword(hash, password)

	if err != nil {
		t.Errorf("La comparación debería ser exitosa para el password correcto: %v", err)
	}
}

func TestComparePassword_Incorrect(t *testing.T) {
	password := "mySecretPassword123"
	wrongPassword := "wrongPassword"
	hash, _ := HashPassword(password)

	err := ComparePassword(hash, wrongPassword)

	if err == nil {
		t.Error("Se esperaba un error al comparar passwords incorrectos, pero fue nil")
	}

	if err != bcrypt.ErrMismatchedHashAndPassword {
		t.Errorf("Se esperaba error 'ErrMismatchedHashAndPassword', se obtuvo: %v", err)
	}
}

func TestHashPassword_TooLong(t *testing.T) {
	password := strings.Repeat("a", 73)

	_, err := HashPassword(password)

	if err == nil {
		t.Fatal("se esperaba error para password mayor a 72 bytes")
	}
}

func TestComparePassword_InvalidHash(t *testing.T) {
	err := ComparePassword("invalid-hash", "password")

	if err == nil {
		t.Fatal("se esperaba error con hash invalido")
	}
}
