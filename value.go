package calltracker

import "time"

// Value is a point-in-time view of outstanding calls. The calls are ordered
// from oldest to newest.
type Value []Call

// Elapsed returns the total amount of time that has elapsed across all
// oustanding calls.
func (v Value) Elapsed() time.Duration {
	var total time.Duration
	for _, call := range v {
		total += call.Elapsed()
	}
	return total
}

// MinElapsed returns the smallest amount of time that has elapsed for any one
// of the outstanding calls. If there are no outstanding calls it will return
// a duration of zero.
func (v Value) MinElapsed() time.Duration {
	length := len(v)
	if length == 0 {
		return time.Duration(0)
	}

	now := time.Now()
	min := now.Sub(v[0].when)
	for i := 1; i < length; i++ {
		elapsed := now.Sub(v[i].when)
		if elapsed < min {
			min = elapsed
		}
	}
	return min
}

// MaxElapsed returns the greatest amount of time that has elapsed for any one
// of the outstanding calls. If there are no outstanding calls it will return
// a duration of zero.
func (v Value) MaxElapsed() time.Duration {
	length := len(v)
	if length == 0 {
		return time.Duration(0)
	}

	now := time.Now()
	max := now.Sub(v[0].when)
	for i := 1; i < length; i++ {
		elapsed := now.Sub(v[i].when)
		if elapsed > max {
			max = elapsed
		}
	}
	return max
}

// Len returns the total number of outstanding calls.
func (v Value) Len() int {
	return len(v)
}
