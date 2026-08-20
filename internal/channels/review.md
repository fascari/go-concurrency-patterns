# Channels: code analysis

## Overview

The producer sends quotes and closes the channel. `Run` drains that stream with `for range` on a receive-only view. `Params.Buffer` only changes when send waits for receive.

## What to look for

### 1. Ownership

`produce` takes `chan<- Quote` and `defer close(quotes)`. The consumer takes `<-chan Quote` and never closes.

### 2. Buffer

Capacity 0 makes send and receive wait for each other. Capacity N lets the producer run ahead by at most N values.

### 3. Cancellation

Producer sends use `select` on `ctx.Done()`. `defer close` still runs, so `consume` is not left blocked, and `errgroup` returns the cancel error after the range ends.

### 4. A buffered signal does not order shared writes

`internal/channels/racefix` increments a shared counter and only uses a buffered channel as a completion signal. `go test -race -tags racefix ./internal/channels/racefix` still reports the race.

## References

- Donovan and Kernighan, *The Go Programming Language*, §8.4-8.6
- *Effective Go*, Channels
- Cox-Buday, *Concurrency in Go*, chapter 3
- Harsanyi, *100 Go Mistakes and How to Avoid Them*, mistakes 67-72
- Go Memory Model, https://go.dev/ref/mem
- `golang.org/x/sync/errgroup`
