package store

import "sync"

type (
	Store[T any] struct {
		mu sync.RWMutex
		m  map[string]T
	}
)

func New[T any]() *Store[T] {
	return &Store[T]{m: make(map[string]T)}
}

func (s *Store[T]) Get(k string) (T, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.m[k]
	return v, ok
}

func (s *Store[T]) Set(k string, v T) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.m[k] = v
}

func (s *Store[T]) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.m)
}
