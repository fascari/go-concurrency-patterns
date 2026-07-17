package unprotected

import (
	"sync"
	"testing"
)

func TestStore_ShouldExposeConcurrentMapRace(t *testing.T) {
	t.Parallel()

	s := New()

	var start sync.WaitGroup
	start.Add(1)

	var workers sync.WaitGroup
	workers.Go(func() {
		start.Wait()
		s.Put("account", 1)
	})
	workers.Go(func() {
		start.Wait()
		s.Load("account")
	})

	start.Done()
	workers.Wait()
}
