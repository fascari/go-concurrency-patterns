// Command concurrency dispatches one scenario by --pattern flag.
// Example: go run ./cmd/concurrency --pattern goroutines
package main

import (
	"context"
	"flag"
	"log"

)

func main() {
	pattern := flag.String("pattern", "", "scenario to run")
	flag.Parse()

	if *pattern == "" {
		log.Fatal("--pattern is required")
	}

	ctx := context.Background()
	_ = ctx
	switch *pattern {
	default:
		log.Fatalf("unknown pattern: %q", *pattern)
	}
}
