//go:build racefix

package racefix

import "testing"

func TestRace_ShouldExposeConcurrentCounterRace(t *testing.T) {
	t.Parallel()

	_ = Race(8)
}
