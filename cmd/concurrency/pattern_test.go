package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParsePattern_ShouldAcceptKnownPatterns(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   string
		want Pattern
	}{
		{
			name: "should parse goroutines",
			in:   "goroutines",
			want: PatternGoroutines,
		},
		{
			name: "should parse channels",
			in:   "channels",
			want: PatternChannels,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := ParsePattern(tt.in)
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestParsePattern_ShouldRejectMissingPattern(t *testing.T) {
	t.Parallel()

	got, err := ParsePattern("")
	require.Empty(t, got)
	require.ErrorIs(t, err, errPatternRequired)
}

func TestParsePattern_ShouldRejectUnknownPattern(t *testing.T) {
	t.Parallel()

	got, err := ParsePattern("mutex")
	require.Empty(t, got)
	require.ErrorContains(t, err, `unknown pattern: "mutex"`)
}
