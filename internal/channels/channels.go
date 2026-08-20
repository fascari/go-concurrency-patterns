package channels

import (
	"context"

	"golang.org/x/sync/errgroup"
)

func Run(ctx context.Context, p Params) (Result, error) {
	if err := ctx.Err(); err != nil {
		return Result{}, err
	}

	quotes := make(chan Quote, p.Buffer)
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return produce(ctx, quotes, p.Quotes)
	})

	received := consume(quotes)
	if err := g.Wait(); err != nil {
		return Result{}, err
	}
	return Result{Received: received}, nil
}

func produce(ctx context.Context, quotes chan<- Quote, items []Quote) error {
	defer close(quotes)

	for _, q := range items {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case quotes <- q:
		}
	}
	return nil
}

func consume(quotes <-chan Quote) []Quote {
	var received []Quote
	for q := range quotes {
		received = append(received, q)
	}
	return received
}
