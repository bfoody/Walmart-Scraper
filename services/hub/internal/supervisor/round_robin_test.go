package supervisor

import "testing"

var expected = []uint{
	0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20,
}

// TestRoundRobin runs a simple test to make sure the round-robin algorithm works.
func TestRoundRobin(t *testing.T) {
	rr := NewRoundRobin()
	n := 50

	for _, item := range expected {
		val := rr.Next(uint(n))
		if val != item {
			t.Fatalf("expected %d, got %d", item, val)
		}
	}
}
