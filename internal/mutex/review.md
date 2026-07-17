# Mutex, RWMutex, and sync.Once: code analysis

## Overview

Unit 2 uses a small in-memory balance store to show how `sync.RWMutex` protects shared state and how `sync.Once` gives the type a usable zero value, while `Transfer` demonstrates why an operation that spans several map accesses must own one critical section from validation through mutation.

`Load` and `Put` each synchronize a single map access, but implementing `Transfer` as a sequence of those methods would leave an unlocked gap between the read and the writes, so the current implementation locks once and updates both balances directly. A mutex isn't attached to the map by the runtime, which means correctness depends on every access following the same locking rule and preserving the same invariant.

## What to look for

### 1. Zero-value initialization and ownership

`New` returns `new(Store)` without allocating `balances` (`store.go:13-15`), so a value returned by the constructor behaves like `var s Store`, and `TestStore_ShouldStoreBalancesForEveryInitializationMode` verifies both forms. Every method that reaches the map calls `initialize`, although `Transfer` rejects a non-positive amount and a self-transfer first, which avoids initializing state for requests that can be rejected without reading it.

The first call to `sync.Once.Do` creates the map, concurrent callers wait for that initialization, and later calls use the completed fast path (`store.go:69-73`), while the memory guarantee provided by `Once` makes the initialized map visible before any waiting call returns. The concurrent zero-value test writes 1,000 distinct keys through one uninitialized `Store`, and the separate `sync.Once` test confirms that an initializer runs exactly once even when 1,000 goroutines call it.

All `Store` methods use pointer receivers, which avoids copying the `sync.RWMutex` and `sync.Once` during ordinary calls, because copying an initialized store would split its synchronization state while both map headers could still refer to the same backing map. The smaller fixture under `testdata/mutexcopy` copies a struct containing `sync.Mutex` and a preallocated map so `go vet` reports the assignment, although the fixture itself only locks the copy and doesn't start a concurrent access through the original value.

### 2. Lock scope and transfer atomicity

`Load` initializes the map, takes `RLock`, reads one balance, and returns the value with its presence flag (`store.go:17-25`), while `Put` initializes the map, takes the exclusive lock, and creates or replaces one balance (`store.go:27-34`). Concurrent readers can overlap, but every writer excludes both readers and other writers while it changes the map.

`Transfer` validates the amount and account pair before locking, then initializes the map, takes one exclusive lock, checks the source and destination accounts, verifies the available funds, and writes both balances before releasing the lock (`store.go:36-67`). Since all checks happen before either write, every error leaves the map unchanged, and callers that use the store methods cannot observe only one side of a successful transfer.

Calling `Load` or `Put` while `Transfer` already holds the write lock would deadlock because `sync.RWMutex` isn't reentrant, while calling them as separate top-level operations would release the lock between steps and lose atomicity, which is why `Transfer` accesses `balances` directly inside its own critical section.

### 3. Error behavior and API boundaries

`errors.go` defines four sentinel errors, and `Transfer` selects them in a fixed order by rejecting an invalid amount, rejecting identical accounts, checking the source, checking the destination, and finally checking for insufficient funds. The table-driven error test covers zero and negative amounts, self-transfer, each missing account, and insufficient funds, then uses `errors.Is` and before-and-after snapshots to confirm both the returned error and the absence of partial state changes.

`Load` reports a missing key through its boolean result, `Put` accepts any string key and any `int64` balance, and `Transfer` requires both accounts to exist before moving a positive amount. Empty keys, negative balances supplied through `Put`, and overflow when adding to the destination balance aren't rejected by the current implementation, and the tests don't claim guarantees for those cases.

### 4. RWMutex behavior used by the store

When a writer calls `Lock` while readers hold the mutex, new `RLock` calls wait until that writer acquires and releases the lock, which prevents a steady stream of readers from postponing the writer indefinitely. The store doesn't attempt a read-to-write upgrade, recursive locking, or speculative `TryLock` calls, and every path uses a blocking lock with a `defer` placed immediately after acquisition.

### 5. What the tests establish

The production tests cover constructor and zero-value initialization, a successful transfer, every declared transfer error, concurrent transfers, concurrent initialization, and direct `sync.Once` behavior. The concurrent transfer test launches 1,000 transfers that alternate direction, expects every call to succeed, and verifies both exact account balances and the preserved total, so it exercises the multi-account invariant under contention rather than checking only for a missing data race.

Running the production package with the race detector passes, and ordinary `go vet` reports no issue, while the intentionally unsafe packages stay under `testdata` so recursive package discovery doesn't include them in a normal `./internal/mutex/...` run.

### 6. Benchmark scope and result

The benchmark defines one interface for the `sync.Mutex` baseline and the production `sync.RWMutex` store, then runs both through `b.RunParallel` with ten reads followed by one write in each cycle (`benchmark_test.go:13-69`). Its initial `Put` happens before `b.ResetTimer`, so map allocation isn't timed, although each production operation still pays the completed `sync.Once.Do` fast path.

Because each benchmark cycle reads and writes in separate calls, it doesn't perform an atomic increment and doesn't inspect the final balance, which makes the result a measurement of lock throughput under a read-heavy workload rather than a correctness test. Five runs on an Intel Core i7-8665U with Go 1.26.1 produced these ranges:

```text
Mutex      782.0 to 908.1 ns/op   0 B/op   0 allocs/op
RWMutex    463.5 to 680.6 ns/op   0 B/op   0 allocs/op
```

`RWMutex` was faster in these runs because readers could overlap, but the result isn't a general ranking, since shorter uncontended sections, more frequent writes, different core counts, and scheduler behavior can reduce the gain or make a plain mutex faster.

### 7. Why the design uses one coarse lock

`sync.Map` works well when entries are written once and read many times, or when goroutines usually operate on separate keys, but a transfer is one read-modify-write operation spanning two keys, and separate `sync.Map.Load` and `sync.Map.Store` calls cannot make that operation atomic. Adding an outer lock would still be necessary and would remove the main reason to choose `sync.Map` here.

Per-account locks could allow unrelated transfers to proceed together, but they would also need stable lock ownership and deterministic account ordering to avoid deadlocks, while the current coarse lock keeps the invariant visible and matches the package's measured scope without extra coordination.

### 8. Race detector and go vet fixtures

The production checks pass with the following commands:

```bash
go test -race -count=1 ./internal/mutex
go vet ./internal/mutex
```

The fixture checks fail on purpose when invoked directly:

```bash
go test -race -count=1 ./internal/mutex/testdata/unprotected
go vet ./internal/mutex/testdata/mutexcopy
```

The unprotected fixture starts one `Put` and one `Load` together, which makes the race detector report concurrent map access, while the copied-lock fixture makes `go vet` report the assignment that copies a value containing `sync.Mutex`. The race detector can only report conflicts reached during execution, whereas `go vet` examines source patterns without needing the copied-lock path to produce a runtime race, so the two tools cover different failure modes.

## References

- Donovan and Kernighan, *The Go Programming Language*, Chapter 9
- Cox-Buday, *Concurrency in Go*, Chapter 3 (synchronization primitives)
- Harsanyi, *100 Go Mistakes and How to Avoid Them*, mistakes 57-63
- Go `sync` package documentation, https://pkg.go.dev/sync
- Go memory model, https://go.dev/ref/mem
- Go race detector, https://go.dev/doc/articles/race_detector
- Zhiyanov, *Gist of Go: Concurrency*, https://antonz.org/go-concurrency/
