package main

import (
	"context"
	"log"

	"github.com/fascari/go-concurrency-patterns/internal/channels"
)

func runChannels(ctx context.Context) error {
	res, err := channels.Run(ctx, channels.Params{
		Quotes: sampleQuotes(),
		Buffer: 0,
	})
	if err != nil {
		return err
	}

	log.Printf("received=%d", len(res.Received))
	return nil
}

func sampleQuotes() []channels.Quote {
	return []channels.Quote{
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
		{
			Symbol: "AMZN",
			Price:  185,
		},
		{
			Symbol: "META",
			Price:  510,
		},
		{
			Symbol: "NVDA",
			Price:  900,
		},
		{
			Symbol: "TSLA",
			Price:  250,
		},
		{
			Symbol: "AMD",
			Price:  160,
		},
		{
			Symbol: "INTC",
			Price:  35,
		},
		{
			Symbol: "ORCL",
			Price:  140,
		},
	}
}
