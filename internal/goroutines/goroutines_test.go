package goroutines

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRun_ShouldReturnCompleteSumForEachWaitGroupStyle(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name  string
		input []int64
		want  int64
	}{
		{
			name:  "should return complete sum for command input",
			input: []int64{1, 2, 3, 4, 5, 6, 7, 8},
			want:  36,
		},
		{
			name:  "should return correct sum for mixed signed values",
			input: []int64{10, -5, 3, -2},
			want:  6,
		},
		{
			name:  "should return zero for nil input",
			input: nil,
			want:  0,
		},
		{
			name:  "should return zero for empty input",
			input: []int64{},
			want:  0,
		},
	}

	modes := []struct {
		name      string
		useModern bool
	}{
		{
			name:      "classic",
			useModern: false,
		},
		{
			name:      "modern",
			useModern: true,
		},
	}

	for _, mode := range modes {
		t.Run(mode.name, func(t *testing.T) {
			t.Parallel()
			for _, tc := range cases {
				t.Run(tc.name, func(t *testing.T) {
					t.Parallel()
					ctx := t.Context()
					p := Params{
						Input:     tc.input,
						UseModern: mode.useModern,
					}
					got, err := Run(ctx, p)
					require.NoError(t, err)
					require.Equal(t, Result{Sum: tc.want}, got)
				})
			}
		})
	}
}

func TestRun_ShouldRejectCanceledContextBeforeStartingWork(t *testing.T) {
	t.Parallel()

	modes := []struct {
		name      string
		useModern bool
	}{
		{
			name:      "classic",
			useModern: false,
		},
		{
			name:      "modern",
			useModern: true,
		},
	}

	for _, mode := range modes {
		t.Run(mode.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithCancel(t.Context())
			cancel()
			p := Params{
				Input:     []int64{1, 2, 3},
				UseModern: mode.useModern,
			}
			got, err := Run(ctx, p)
			require.Equal(t, Result{}, got)
			require.True(t, errors.Is(err, context.Canceled))
		})
	}
}
