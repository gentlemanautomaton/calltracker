package calltracker

import "time"

// Call is an oustanding call.
type Call struct {
	number uint64
	when   time.Time
}

// Number returns the number of calls that have taken place before this call.
func (c Call) Number() uint64 {
	return c.number
}

// When returns the time at which the call started.
func (c Call) When() time.Time {
	return c.when
}

// Elapsed returns the amount of time that has elapsed since the call started.
func (c Call) Elapsed() time.Duration {
	return time.Now().Sub(c.when)
}
