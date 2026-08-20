package main

import (
	"context"
	"errors"
	"fmt"
)

const (
	PatternGoroutines Pattern = "goroutines"
	PatternChannels   Pattern = "channels"
)

var (
	errPatternRequired = errors.New("--pattern is required")
	runners            = map[Pattern]func(context.Context) error{
		PatternGoroutines: runGoroutines,
		PatternChannels:   runChannels,
	}
)

type (
	Pattern string
)

func ParsePattern(s string) (Pattern, error) {
	if s == "" {
		return "", errPatternRequired
	}

	p := Pattern(s)
	if _, ok := runners[p]; !ok {
		return "", fmt.Errorf("unknown pattern: %q", s)
	}
	return p, nil
}

func (p Pattern) run(ctx context.Context) error {
	return runners[p](ctx)
}
