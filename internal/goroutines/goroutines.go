package goroutines

import (
	"context"
	"sync"
)

func Run(ctx context.Context, p Params) (Result, error) {
	if err := ctx.Err(); err != nil {
		return Result{}, err
	}

	var (
		mu  sync.Mutex
		wg  sync.WaitGroup
		sum int64
	)

	add := func(v int64) {
		mu.Lock()
		sum += v
		mu.Unlock()
	}

	launch := newLauncher(&wg, add, p.UseModern)
	for _, v := range p.Input {
		launch(v)
	}

	wg.Wait()
	return Result{Sum: sum}, nil
}

func newLauncher(wg *sync.WaitGroup, add func(int64), useModern bool) func(int64) {
	if useModern {
		return func(v int64) {
			wg.Go(func() {
				add(v)
			})
		}
	}
	return func(v int64) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			add(v)
		}()
	}
}
