package main

import (
	"context"
	"flag"
	"log"

	"github.com/fascari/go-concurrency-patterns/internal/goroutines"
)

func main() {
	pattern := flag.String("pattern", "", "scenario to run")
	flag.Parse()

	if *pattern == "" {
		log.Fatal("--pattern is required")
	}

	ctx := context.Background()
	switch *pattern {
	case "goroutines":
		res, err := goroutines.Run(ctx, goroutines.Params{
			Input:     []int64{1, 2, 3, 4, 5, 6, 7, 8},
			UseModern: true,
		})
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("sum=%d", res.Sum)
	default:
		log.Fatalf("unknown pattern: %q", *pattern)
	}
}
