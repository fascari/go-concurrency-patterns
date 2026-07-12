package store

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStore_ShouldSetGetAndCountEntries(t *testing.T) {
	s := New[int]()
	s.Set("a", 1)
	s.Set("b", 2)

	got, ok := s.Get("a")
	require.True(t, ok)
	require.Equal(t, 1, got)
	require.Equal(t, 2, s.Len())
}
