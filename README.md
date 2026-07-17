# Go concurrency patterns

This repository works through Go concurrency one topic at a time, pairing small examples with race-tested code, benchmarks when they help answer a performance question, and review notes that explain the synchronization choices. The unit table defines the intended study sequence, while the repository tree shows which units have code today.

## What this is

Runnable units keep the code narrow enough to expose the behavior of one primitive without hiding it behind application infrastructure, while text-only units collect the mental models and decision guidance shared across examples. Tests cover normal behavior, error paths, and concurrent execution where the distinction matters, while benchmarks and static-analysis fixtures appear only when they explain behavior that ordinary tests don't expose.

These examples are teaching code rather than production services, so they favor visible synchronization and focused tests over persistence, configuration, and operational concerns.

## Requirements

- Go 1.26.1
- `mise` for the task shortcuts shown below
- `golangci-lint` for `mise run lint`

## Quick start

Install the configured Go version, run the race-tested suite, execute the currently registered CLI pattern, or start the linter through `mise`:

```bash
mise install
mise run test
mise run pattern goroutines
mise run lint
```

The `pattern` task forwards its argument to the dispatcher under `cmd/concurrency`, while `test` and `lint` discover the packages present in the repository, so adding a unit doesn't require another README command.

## Units

The table describes the curriculum rather than implementation status, which keeps it stable as packages are added over time.

| Unit | Topic | Location |
|---|---|---|
| 0 | Foundations | `docs/foundations.md` |
| 1 | Goroutines and `WaitGroup.Go` | `internal/goroutines/` |
| 2 | Shared state and mutexes | `internal/mutex/` |
| 3 | Channels | `internal/channels/` |
| 4 | `select` | `internal/selects/` |
| 5 | Context cancellation | `internal/contexts/` |
| 6 | Worker pool | `internal/transfers/` |
| 7 | Pipeline | `internal/deposits/` |
| 8 | Fan-out and fan-in | `internal/pricefeeds/` |
| 9 | `errgroup` | `internal/snapshot/` |
| 10 | Goroutine leaks | `internal/outbox/` |
| 11 | Testing concurrency | `internal/synctest/` |
| 12 | Semaphores | `internal/semaphores/` |
| 13 | Atomics | `internal/atomics/` |
| 14 | Decision guide | `docs/decision-guide.md` |

## References

- *The Go Memory Model*, https://go.dev/ref/mem
- *Effective Go: Concurrency*, https://go.dev/doc/effective_go#concurrency
- Rob Pike, *Concurrency Is Not Parallelism*, https://go.dev/blog/waza-talk
- Cox-Buday, *Concurrency in Go*, Chapters 1-5
- Donovan and Kernighan, *The Go Programming Language*, Chapters 8-9
- Harsanyi, *100 Go Mistakes and How to Avoid Them*, Chapters 8-9
- Sameer Ajmani, *Go Concurrency Patterns: Pipelines and Cancellation*, https://go.dev/blog/pipelines
- Go 1.25 release notes, `sync.WaitGroup.Go`, `go vet` WaitGroup checks, and `testing/synctest`, https://go.dev/doc/go1.25
- Anthony Zhiyanov, *Gist of Go: Concurrency*, https://antonz.org/go-concurrency/
- Hoare, *Communicating Sequential Processes*, CACM, 1978
