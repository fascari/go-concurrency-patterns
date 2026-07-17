package mutex

import "sync"

type (
	Store struct {
		mu       sync.RWMutex
		once     sync.Once
		balances map[string]int64
	}
)

func New() *Store {
	return new(Store)
}

func (s *Store) Load(key string) (int64, bool) {
	s.initialize()

	s.mu.RLock()
	defer s.mu.RUnlock()

	balance, ok := s.balances[key]
	return balance, ok
}

func (s *Store) Put(key string, balance int64) {
	s.initialize()

	s.mu.Lock()
	defer s.mu.Unlock()

	s.balances[key] = balance
}

func (s *Store) Transfer(source, destination string, amount int64) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}
	if source == destination {
		return ErrSelfTransfer
	}

	s.initialize()

	s.mu.Lock()
	defer s.mu.Unlock()

	sourceBalance, ok := s.balances[source]
	if !ok {
		return ErrAccountNotFound
	}

	destinationBalance, ok := s.balances[destination]
	if !ok {
		return ErrAccountNotFound
	}

	if sourceBalance < amount {
		return ErrInsufficientFunds
	}

	s.balances[source] = sourceBalance - amount
	s.balances[destination] = destinationBalance + amount

	return nil
}

func (s *Store) initialize() {
	s.once.Do(func() {
		s.balances = make(map[string]int64)
	})
}
