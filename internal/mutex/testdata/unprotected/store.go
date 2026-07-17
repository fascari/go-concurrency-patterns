package unprotected

type Store struct {
	balances map[string]int64
}

func New() *Store {
	return new(Store{
		balances: make(map[string]int64),
	})
}

func (s *Store) Load(key string) (int64, bool) {
	balance, ok := s.balances[key]
	return balance, ok
}

func (s *Store) Put(key string, balance int64) {
	s.balances[key] = balance
}
