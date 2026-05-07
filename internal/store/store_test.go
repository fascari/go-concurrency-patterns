package store

import "testing"

func TestStoreSetGetLen(t *testing.T) {
	s := New[int]()
	s.Set("a", 1)
	s.Set("b", 2)

	if got, ok := s.Get("a"); !ok || got != 1 {
		t.Fatalf("Get(a): got=%d ok=%v", got, ok)
	}
	if got := s.Len(); got != 2 {
		t.Fatalf("Len: got=%d want=2", got)
	}
}
