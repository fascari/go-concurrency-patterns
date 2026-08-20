package channels

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRun_ShouldDeliverAllQuotesForEachBufferMode(t *testing.T) {
	t.Parallel()

	quotes := []Quote{
		{
			Symbol: "AAPL",
			Price:  190,
		},
		{
			Symbol: "GOOG",
			Price:  175,
		},
		{
			Symbol: "MSFT",
			Price:  420,
		},
	}

	modes := []struct {
		name   string
		buffer int
	}{
		{
			name:   "unbuffered",
			buffer: 0,
		},
		{
			name:   "buffered",
			buffer: 4,
		},
	}

	for _, mode := range modes {
		t.Run(mode.name, func(t *testing.T) {
			t.Parallel()

			got, err := Run(t.Context(), Params{
				Quotes: quotes,
				Buffer: mode.buffer,
			})
			require.NoError(t, err)
			require.Equal(t, Result{Received: quotes}, got)
		})
	}
}

func TestRun_ShouldReturnEmptyResultForEmptyQuotes(t *testing.T) {
	t.Parallel()

	got, err := Run(t.Context(), Params{
		Quotes: nil,
		Buffer: 0,
	})
	require.NoError(t, err)
	require.Equal(t, Result{Received: nil}, got)
}

func TestRun_ShouldRejectCanceledContextBeforeStartingWork(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	got, err := Run(ctx, Params{
		Quotes: []Quote{
			{
				Symbol: "AAPL",
				Price:  1,
			},
		},
		Buffer: 0,
	})
	require.Equal(t, Result{}, got)
	require.True(t, errors.Is(err, context.Canceled))
}
