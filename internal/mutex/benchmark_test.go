package mutex

import (
	"sync"
	"testing"
)

const (
	benchmarkAccount       = "account"
	benchmarkReadsPerWrite = 10
)

type (
	benchmarkStore interface {
		Load(key string) (int64, bool)
		Put(key string, balance int64)
	}

	mutexStore struct {
		mu       sync.Mutex
		balances map[string]int64
	}
)

func BenchmarkStoreReadHeavyContention(b *testing.B) {
	b.Run("Mutex", func(b *testing.B) {
		runReadHeavyContention(b, newMutexStore())
	})

	b.Run("RWMutex", func(b *testing.B) {
		runReadHeavyContention(b, New())
	})
}

func newMutexStore() *mutexStore {
	return new(mutexStore{
		balances: make(map[string]int64),
	})
}

func (s *mutexStore) Load(key string) (int64, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	balance, ok := s.balances[key]
	return balance, ok
}

func (s *mutexStore) Put(key string, balance int64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.balances[key] = balance
}

func runReadHeavyContention(b *testing.B, s benchmarkStore) {
	s.Put(benchmarkAccount, 0)
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		var balance int64
		for pb.Next() {
			for range benchmarkReadsPerWrite {
				balance, _ = s.Load(benchmarkAccount)
			}
			s.Put(benchmarkAccount, balance+1)
		}
	})
}
