//go:build racefix

package racefix

// Race uses a buffered completion channel that does not order concurrent
// increments of n relative to each other.
func Race(workers int) int {
	var n int
	signals := make(chan struct{}, workers)

	for range workers {
		go func() {
			n++
			signals <- struct{}{}
		}()
	}

	for range workers {
		<-signals
	}
	return n
}
