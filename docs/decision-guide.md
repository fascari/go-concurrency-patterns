# Decision guide

Use this guide when choosing which concurrency primitive to reach for. It assumes a single binary written in idiomatic Go.

## Start with the data relationship

- One goroutine owns the data and passes it to another: use a **channel**.
- Multiple goroutines read and write the same state: use a **mutex**.
- A goroutine must wait for many others to finish: use a **WaitGroup**.
- Many goroutines must finish and the first error matters: use an **errgroup**.
- Work has a deadline or can be canceled: use a **context**.

## Channels vs mutexes

Use channels when:
- Goroutines communicate by passing values.
- You want backpressure (`ch <- v` blocks when the buffer is full).
- Ownership of a value transfers from producer to consumer.

Use mutexes when:
- State is naturally shared, like a cache or a counter.
- Multiple goroutines perform small, frequent updates.
- Converting shared state into message passing would obscure the logic.

## Worker pool vs pipeline

Use a **worker pool** when you have many independent tasks and a fixed number of workers. A **pipeline** fits when each task has sequential stages and you want each stage to run concurrently.

## errgroup vs WaitGroup

Use `errgroup.Group` when the failure of any goroutine should cancel the others and you need to surface the first error. Choose `sync.WaitGroup` when goroutines are independent and errors can be handled locally.

## Context

Always pass `ctx` as the first argument of functions that may block or perform I/O. Do not store `context.Context` in a struct. Never pass `context.TODO()` in production code without a clear plan to replace it.

## synctest

Use `testing/synctest` for tests that involve timeouts, timers, or goroutine scheduling. It gives deterministic control over time and lets the test observe blocked goroutines. Do not use it to replace `-race`.

## Anti-patterns

- Sharing memory without synchronization.
- Closing a channel from a receiver.
- Ignoring `ctx.Done()` in long-running loops.
- Using `time.Sleep` to fix a race.
- Mixing mutex-protected state with channel-based synchronization for the same invariant.

## References

- Effective Go: Concurrency, https://go.dev/doc/effective_go#concurrency
- Sameer Ajmani: Pipelines and Cancellation, https://go.dev/blog/pipelines
- Go Memory Model, https://go.dev/ref/mem
