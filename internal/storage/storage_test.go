package storage

import (
	"bytes"
	"testing"
	"time"
)

func TestNewMemoryStore(t *testing.T) {
	ms := NewMemoryStore()
	if ms == nil {
		t.Fatal("NewMemoryStore devolvió nil")
	}
	if ms.store == nil {
		t.Error("El mapa interno store no debería ser nil")
	}
}

func TestMemoryStore_SetAndGet(t *testing.T) {
	ms := NewMemoryStore()
	key := "test-key"
	value := []byte("test-value")

	ms.Set(key, value, 0)

	got, found := ms.Get(key)
	if !found {
		t.Fatal("Se esperaba encontrar la clave, pero no se encontró")
	}

	if !bytes.Equal(got, value) {
		t.Errorf("Valor esperado %s, obtenido %s", value, got)
	}
}

func TestMemoryStore_Get_NonExistent(t *testing.T) {
	ms := NewMemoryStore()

	_, found := ms.Get("missing-key")
	if found {
		t.Error("Se esperaba que la clave no existiera, pero fue encontrada")
	}
}

func TestMemoryStore_Del(t *testing.T) {
	ms := NewMemoryStore()
	key := "to-delete"
	ms.Set(key, []byte("value"), 0)

	ms.Del(key)

	_, found := ms.Get(key)
	if found {
		t.Error("La clave todavía existe después de llamar a Del")
	}
}

func TestMemoryStore_TTL_Expiration(t *testing.T) {
	ms := NewMemoryStore()
	key := "expiring-key"
	value := []byte("will-expire")

	ms.Set(key, value, 1)

	_, found := ms.Get(key)
	if !found {
		t.Fatal("La clave debería existir inmediatamente después del Set")
	}

	time.Sleep(1100 * time.Millisecond)

	_, found = ms.Get(key)
	if found {
		t.Error("La clave no debería existir después de que el TTL haya expirado")
	}

	ms.mu.RLock()
	_, exists := ms.store[key]
	ms.mu.RUnlock()

	if exists {
		t.Error("La clave debería haber sido eliminada del mapa interno")
	}
}

func TestMemoryStore_TTL_NoExpiration(t *testing.T) {
	ms := NewMemoryStore()
	key := "non-expiring"
	value := []byte("forever")

	ms.Set(key, value, 0)

	time.Sleep(200 * time.Millisecond)

	_, found := ms.Get(key)
	if !found {
		t.Error("La clave sin TTL no debería expirar")
	}
}

func TestMemoryStore_Immutability(t *testing.T) {
	ms := NewMemoryStore()
	key := "immutable"
	original := []byte("original")

	ms.Set(key, original, 0)

	got, _ := ms.Get(key)

	got[0] = 'X'

	gotAgain, _ := ms.Get(key)

	if bytes.Equal(original, got) {
		t.Error("Modificar el valor devuelto no debería afectar al original (se devolvió la misma referencia)")
	}

	if string(gotAgain) != "original" {
		t.Errorf("El valor almacenado cambió inesperadamente a: %s", gotAgain)
	}
}

func TestIsExpired(t *testing.T) {
	if isExpired(0) {
		t.Error("isExpired(0) debería ser falso")
	}

	future := time.Now().Add(1 * time.Hour).UnixNano()
	if isExpired(future) {
		t.Error("isExpired(future) debería ser falso")
	}

	past := time.Now().Add(-1 * time.Second).UnixNano()
	if !isExpired(past) {
		t.Error("isExpired(past) debería ser true")
	}
}

func TestCloneBytes(t *testing.T) {
	original := []byte("test")
	clone := cloneBytes(original)

	if !bytes.Equal(original, clone) {
		t.Error("El clon no tiene el mismo contenido que el original")
	}

	if &original[0] == &clone[0] {
		t.Error("El clon y el original apuntan a la misma dirección de memoria")
	}
}
