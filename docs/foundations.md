# Foundations

Mental models, primitives, and vocabulary that anchor every scenario in this study. Keep it text-only. Runnable demos live under `internal/`.

## Goroutines

A goroutine is a lightweight thread managed by the Go runtime. The `go` statement starts a function call in a new goroutine. Goroutines share memory, so synchronization is required when goroutines access the same variable concurrently.

## Communication: channels and select

Channels are typed conduits that let goroutines communicate and synchronize. Use channels to pass ownership of data, to signal completion, or to distribute work. The `select` statement waits on multiple channel operations and chooses one that is ready.

## Shared state: sync primitives

Use `sync.Mutex` or `sync.RWMutex` when goroutines must share mutable state that is not naturally owned by a single goroutine. Keep critical sections small and hold locks for the shortest time possible.

## Coordination: WaitGroup and errgroup

`sync.WaitGroup` waits for a collection of goroutines to finish. In Go 1.25, `WaitGroup.Go` removes the need for a wrapper closure. `golang.org/x/sync/errgroup` extends this pattern by collecting the first error returned by any goroutine.

## Cancellation: context

`context.Context` carries deadlines, cancellation signals, and request-scoped values. Pass `ctx` as the first argument of functions that do I/O or long-running work. Check `ctx.Err()` and respect cancellation promptly.

## Testing concurrency

Run every test with `-race`. Use `testing/synctest` for deterministic time and goroutine-aware testing of code that depends on timeouts, timers, or sleeps. Avoid tests that depend on specific scheduling order.

## Shared memory vs CSP

Go supports both models. Prefer channels when one goroutine produces data and another consumes it. Prefer mutexes when multiple goroutines mutate shared data structures. Do not mix the two models casually. A goroutine that reads a shared map should not also rely on a channel for the same invariant.

## Happens-before

The Go Memory Model defines which reads are guaranteed to observe which writes. Use synchronization primitives or channel sends and receives to establish happens-before relationships. Do not rely on undocumented side effects of goroutine scheduling.

## Race detector

`go test -race` instruments the program to detect data races. A race is any unsynchronized access to the same memory location where at least one access is a write. Treat every race as a bug.

## References

- Go Memory Model, https://go.dev/ref/mem
- Effective Go: Concurrency, https://go.dev/doc/effective_go#concurrency
- The Go Programming Language, Donovan and Kernighan, Ch 8-9
- Hoare: Communicating Sequential Processes (CACM, 1978)
