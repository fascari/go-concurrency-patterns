# Goroutines and WaitGroup: code analysis

## Overview

Unit 1 shows two ways to synchronize goroutines with `sync.WaitGroup`, the classic `Add`/`Done` ritual and the `WaitGroup.Go` helper from Go 1.25 that folds the boilerplate into a single call. Each input element launches one goroutine, they sum values concurrently, and the parent waits for completion before reading the result.

A goroutine isn't an OS thread, but a schedulable task managed entirely by the Go runtime, starting with a stack of roughly 2KB that grows or shrinks on demand through stack copying. The runtime uses an M:N scheduler with three abstractions (G for goroutine, M for machine which is an OS thread, P for processor which is a logical context holding a run queue) and multiplexes thousands of goroutines onto a small pool of OS threads via work stealing and handoff, which makes goroutines cheap enough to start by the thousand (see [Understanding Go's Scheduler](https://www.linkedin.com/pulse/understanding-gos-scheduler-how-goroutine-management-works-ascari-o6kuf/)).

## What to look for

### 1. Two WaitGroup patterns

The `newLauncher` function picks between two approaches, where in the classic path the caller invokes `wg.Add(1)` before the goroutine starts and `defer wg.Done()` inside it, giving explicit control over when the counter increments, while the modern path uses `wg.Go(func() { ... })`, which does exactly the same thing internally (calling `wg.Add(1)`, launching the closure, and calling `wg.Done` on completion) but removes the opportunity to forget any of those steps. Both produce the same result and the choice is stylistic, depending on what the team finds more readable.

`wg.Add` must happen before the goroutine starts to avoid a timing race where `wg.Add(1)` running inside the goroutine body means the scheduler might not have scheduled that goroutine yet by the time the parent reaches `wg.Wait()`, so `Wait()` could see a counter of zero and return prematurely. In the classic path `Add` is called right before `go func`, which is correct, and `WaitGroup.Go` handles this ordering internally by calling `Add` before `go`.

The loop iterates over `p.Input` with `for range`, which in Go 1.22+ creates a fresh variable per iteration so each goroutine closure captures its own value. Before Go 1.22 the loop reused the same variable address, and launching goroutines inside the loop would capture the last value every time, a bug that's shipped in production code more times than anyone wants to admit. Go 1.22 removed this trap by giving each iteration its own variable.

### 2. Shared state and mutex

The goroutines all write to the same `sum` variable through a closure, which is a classic race condition, so a `sync.Mutex` guards each write because multiple goroutines access the same memory concurrently without a happens-before relation between them, and the race detector would flag the access immediately. Could `atomic.AddInt64` replace it? Yes, it'd be faster, but the mutex is more general and makes the critical section visible in a way that matters for teaching.

### 3. Context checking

`Run` checks `ctx.Err()` before starting any work, which is the standard Go pattern for respecting cancellation at function entry, and `TestRun_ShouldRejectCanceledContextBeforeStartingWork` verifies this. Once the goroutines launch, `Run` doesn't propagate cancellation to them, making this an all-or-nothing admission policy that's fine for bounded input but would need a `select` over `ctx.Done()` if the work were long-running.

## References

- Donovan and Kernighan, *The Go Programming Language*, 8.1-8.4
- Cox-Buday, *Concurrency in Go*, Ch 3 (sync primitives)
- Harsanyi, *100 Go Mistakes and How to Avoid Them*, mistakes 60-61
- Go 1.25 release notes: `sync.WaitGroup.Go`, `go vet waitgroup`
- Go 1.22 release notes: loop variable change
- Go Memory Model, https://go.dev/ref/mem
- Zhiyanov, *Gist of Go: Concurrency*, https://antonz.org/go-concurrency/
- Gist of Go: Concurrency: https://antonz.org/go-concurrency/