package main

import (
	"context"
	"log"

	"github.com/fascari/go-concurrency-patterns/internal/goroutines"
)

func runGoroutines(ctx context.Context) error {
	res, err := goroutines.Run(ctx, goroutines.Params{
		Input:     []int64{1, 2, 3, 4, 5, 6, 7, 8},
		UseModern: true,
	})
	if err != nil {
		return err
	}

	log.Printf("sum=%d", res.Sum)
	return nil
}
