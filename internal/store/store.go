// Package store provides an in-memory key/value store guarded by
// sync.RWMutex. Intentionally NOT sync.Map: pedagogy needs explicit
// critical sections so scenarios can illustrate races and -race output.
package store

import "sync"

// Store is a goroutine-safe map keyed by string.
type Store[T any] struct {
	mu sync.RWMutex
	m  map[string]T
}

// New returns an empty Store.
func New[T any]() *Store[T] {
	return &Store[T]{m: make(map[string]T)}
}

// Get returns the value and a presence flag.
func (s *Store[T]) Get(k string) (T, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.m[k]
	return v, ok
}

// Set assigns v to key k.
func (s *Store[T]) Set(k string, v T) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.m[k] = v
}

// Len returns the current number of entries.
func (s *Store[T]) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.m)
}
