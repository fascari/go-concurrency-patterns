# Go concurrency patterns

A literature-grounded study of Go concurrency primitives and patterns. Each unit isolates one concept, implements a minimal runnable example, and verifies it with the race detector.

## What this is

This repository contains annotated implementations of common Go concurrency patterns. The goal is to understand each primitive from first principles and see it used correctly in a small, self-contained program.

The code isn't production-grade. It's optimized for clarity, correct synchronization, and good test coverage.

## Requirements

- Go 1.26.1 or later

## Quick start

This project uses `mise` for task running.

```bash
# Run tests with the race detector
mise run test

# Run a concurrency pattern
mise run pattern goroutines

# Run the linter
mise run lint

# Format code
mise run fmt
```

## Units

| Unit | Topic | Package |
|---|---|---|---|
| 0 | Foundations | `docs/foundations.md` |
| 1 | Goroutines and WaitGroup.Go | `internal/goroutines/` |
| 2 | Shared state and mutexes | `internal/store/` |
| 3 | Channels | `internal/channels/` |
| 4 | select | `internal/selects/` |
| 5 | Context cancellation | `internal/contexts/` |
| 6 | Worker pool | `internal/transfers/` |
| 7 | Pipeline | `internal/deposits/` |
| 8 | Fan-out / fan-in | `internal/pricefeeds/` |
| 9 | errgroup | `internal/snapshot/` |
| 10 | Goroutine leaks | `internal/outbox/` |
| 11 | Testing concurrency | `internal/synctest/` |
| 12 | Semaphores | `internal/semaphores/` |
| 13 | Atomics | `internal/atomics/` |
| 14 | Decision guide | `docs/decision-guide.md` |

## References

- Go Memory Model, https://go.dev/ref/mem
- Effective Go: Concurrency, https://go.dev/doc/effective_go#concurrency
- Rob Pike: Concurrency is not Parallelism, https://go.dev/blog/waza-talk
- Concurrency in Go, Cox-Buday, Ch. 1–5
- The Go Programming Language, Donovan and Kernighan, Ch. 8–9
- 100 Go Mistakes and How to Avoid Them, Harsanyi, Ch. 8–9
- Sameer Ajmani: Pipelines and Cancellation, https://go.dev/blog/pipelines
- Go 1.25 release notes, `sync.WaitGroup.Go`, `go vet` waitgroup checks, `testing/synctest`
- Anthony Zhiyanov, *Gist of Go: Concurrency*, https://antonz.org/go-concurrency/
- Hoare: Communicating Sequential Processes (CACM, 1978)
