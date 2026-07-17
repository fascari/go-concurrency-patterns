package mutex

import (
	"errors"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	sourceAccount             = "source"
	destinationAccount        = "destination"
	missingSourceAccount      = "missing-source"
	missingDestinationAccount = "missing-destination"
	testSourceBalance         = int64(100)
	testDestinationBalance    = int64(50)
	concurrentStartingBalance = int64(1_000)
	concurrentCallCount       = 1_000
)

func TestStore_ShouldStoreBalancesForEveryInitializationMode(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		newStore func() *Store
	}{
		{
			name:     "should store balances from New",
			newStore: New,
		},
		{
			name: "should store balances from zero value",
			newStore: func() *Store {
				return new(Store)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := tt.newStore()
			s.Put(sourceAccount, testSourceBalance)

			balance, ok := s.Load(sourceAccount)
			require.True(t, ok)
			require.Equal(t, testSourceBalance, balance)
		})
	}
}

func TestStoreTransfer_ShouldMoveBalanceBetweenAccounts(t *testing.T) {
	t.Parallel()

	s := New()
	s.Put(sourceAccount, testSourceBalance)
	s.Put(destinationAccount, testDestinationBalance)

	err := s.Transfer(sourceAccount, destinationAccount, 40)
	require.NoError(t, err)

	sourceBalance, sourceExists := s.Load(sourceAccount)
	destinationBalance, destinationExists := s.Load(destinationAccount)
	require.True(t, sourceExists)
	require.True(t, destinationExists)
	require.Equal(t, int64(60), sourceBalance)
	require.Equal(t, int64(90), destinationBalance)
}

func TestStoreTransfer_ShouldLeaveBalancesUnchangedForEveryError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		source      string
		destination string
		amount      int64
		wantErr     error
	}{
		{
			name:        "should reject zero amount",
			source:      sourceAccount,
			destination: destinationAccount,
			amount:      0,
			wantErr:     ErrInvalidAmount,
		},
		{
			name:        "should reject negative amount",
			source:      sourceAccount,
			destination: destinationAccount,
			amount:      -1,
			wantErr:     ErrInvalidAmount,
		},
		{
			name:        "should reject self transfer",
			source:      sourceAccount,
			destination: sourceAccount,
			amount:      1,
			wantErr:     ErrSelfTransfer,
		},
		{
			name:        "should reject missing source",
			source:      missingSourceAccount,
			destination: destinationAccount,
			amount:      1,
			wantErr:     ErrAccountNotFound,
		},
		{
			name:        "should reject missing destination",
			source:      sourceAccount,
			destination: missingDestinationAccount,
			amount:      1,
			wantErr:     ErrAccountNotFound,
		},
		{
			name:        "should reject insufficient funds",
			source:      sourceAccount,
			destination: destinationAccount,
			amount:      testSourceBalance + 1,
			wantErr:     ErrInsufficientFunds,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := New()
			s.Put(sourceAccount, testSourceBalance)
			s.Put(destinationAccount, testDestinationBalance)

			sourceBalanceBefore, sourceExistsBefore := s.Load(tt.source)
			destinationBalanceBefore, destinationExistsBefore := s.Load(tt.destination)

			err := s.Transfer(tt.source, tt.destination, tt.amount)
			require.Error(t, err)
			require.True(t, errors.Is(err, tt.wantErr))

			sourceBalanceAfter, sourceExistsAfter := s.Load(tt.source)
			destinationBalanceAfter, destinationExistsAfter := s.Load(tt.destination)
			require.Equal(t, sourceBalanceBefore, sourceBalanceAfter)
			require.Equal(t, sourceExistsBefore, sourceExistsAfter)
			require.Equal(t, destinationBalanceBefore, destinationBalanceAfter)
			require.Equal(t, destinationExistsBefore, destinationExistsAfter)
		})
	}
}

func TestStoreTransfer_ShouldPreserveTotalBalanceUnderConcurrentCalls(t *testing.T) {
	t.Parallel()

	s := New()
	s.Put(sourceAccount, concurrentStartingBalance)
	s.Put(destinationAccount, concurrentStartingBalance)

	var failures atomic.Int64
	var wg sync.WaitGroup
	for i := range concurrentCallCount {
		wg.Go(func() {
			source, destination := sourceAccount, destinationAccount
			if i%2 != 0 {
				source, destination = destination, source
			}

			if err := s.Transfer(source, destination, 1); err != nil {
				failures.Add(1)
			}
		})
	}
	wg.Wait()

	require.Zero(t, failures.Load())

	sourceBalance, sourceExists := s.Load(sourceAccount)
	destinationBalance, destinationExists := s.Load(destinationAccount)
	require.True(t, sourceExists)
	require.True(t, destinationExists)
	require.Equal(t, concurrentStartingBalance, sourceBalance)
	require.Equal(t, concurrentStartingBalance, destinationBalance)
	require.Equal(t, concurrentStartingBalance*2, sourceBalance+destinationBalance)
}

func TestStore_ShouldSafelyInitializeZeroValueUnderConcurrentCalls(t *testing.T) {
	t.Parallel()

	var s Store
	var wg sync.WaitGroup
	for i := range concurrentCallCount {
		wg.Go(func() {
			s.Put(strconv.Itoa(i), int64(i))
		})
	}
	wg.Wait()

	for i := range concurrentCallCount {
		balance, ok := s.Load(strconv.Itoa(i))
		require.True(t, ok)
		require.Equal(t, int64(i), balance)
	}
}

func TestSyncOnce_ShouldExecuteInitializerExactlyOnceAcrossConcurrentCalls(t *testing.T) {
	t.Parallel()

	var once sync.Once
	var initializations atomic.Int64
	var wg sync.WaitGroup
	for range concurrentCallCount {
		wg.Go(func() {
			once.Do(func() {
				initializations.Add(1)
			})
		})
	}
	wg.Wait()

	require.Equal(t, int64(1), initializations.Load())
}
