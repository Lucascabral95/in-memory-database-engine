package storage

import (
	"fmt"
	"sync"
	"time"
)

type Entry struct {
	Value      []byte
	Expiration int64
}

type MemoryStore struct {
	mu    sync.RWMutex
	store map[string]Entry
}

func NewMemoryStore() *MemoryStore {
	s := &MemoryStore{
		store: make(map[string]Entry),
	}

	go s.startCleanup()
	return s
}

func (m *MemoryStore) Set(key string, value []byte, ttlSeconds int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	var exp int64
	if ttlSeconds > 0 {
		exp = time.Now().Add(time.Duration(ttlSeconds) * time.Second).UnixNano()
	}

	m.store[key] = Entry{
		Value:      cloneBytes(value),
		Expiration: exp,
	}
}

func (m *MemoryStore) Get(key string) ([]byte, bool) {
	m.mu.RLock()
	entry, exists := m.store[key]
	if !exists {
		m.mu.RUnlock()
		return nil, false
	}

	if !isExpired(entry.Expiration) {
		valueCopy := cloneBytes(entry.Value)
		m.mu.RUnlock()
		return valueCopy, true
	}
	m.mu.RUnlock()

	m.mu.Lock()
	defer m.mu.Unlock()

	entry, exists = m.store[key]
	if !exists {
		return nil, false
	}

	if isExpired(entry.Expiration) {
		delete(m.store, key)
		return nil, false
	}

	return cloneBytes(entry.Value), true
}

func (m *MemoryStore) Del(key string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.store, key)
}

func (m *MemoryStore) startCleanup() {
	ticker := time.NewTicker(1 * time.Minute)
	for range ticker.C {
		m.mu.Lock()
		for k, v := range m.store {
			if isExpired(v.Expiration) {
				delete(m.store, k)
				fmt.Printf("[GC] Limpiando carrito expirado: %s\n", k)
			}
		}
		m.mu.Unlock()
	}
}

func isExpired(expiration int64) bool {
	return expiration > 0 && time.Now().UnixNano() > expiration
}

func cloneBytes(input []byte) []byte {
	if input == nil {
		return nil
	}

	out := make([]byte, len(input))
	copy(out, input)
	return out
}
