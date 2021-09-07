package supervisor

// RoundRobin keeps track of the currently selected item in a round robin order.
type RoundRobin struct {
	current uint // index of the last dispatched/current value
}

// NewRoundRobin creates and returns a *RoundRobin.
func NewRoundRobin() *RoundRobin {
	return &RoundRobin{
		0,
	}
}

// Next generates the next round robin value with the supplied length of objects
// to choose from.
func (r *RoundRobin) Next(length uint) uint {
	// Avoid divide by zero.
	if length == 0 {
		return 0
	}

	i := r.current

	r.current = i + 1%length

	return i
}
