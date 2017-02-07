package calltracker

import "sort"

// Note: The code in this file is derived from example and source code
// in the container/heap package in the Go programming language standard
// library.

// Queue is an ordered queue of calls. The oldest calls are at the beginning,
// and the newest calls are at the end. Age is determined by comparing the
// number of each call.
type Queue []Call

// Push inserts the call c at the end of the queue without verifying the
// correctness of its ordering.
func (q *Queue) Push(c Call) {
	*q = append(*q, c)
}

// Remove removes the element at index i from the queue.
// The complexity is O(log(n)) where n = len(q).
func (q *Queue) Remove(i int) {
	*q = append((*q)[:i], (*q)[i+1:]...)
}

// Index returns the index of the first call with the given number in the queue.
// It will return -1 if such a call is not present in the queue.
func (q Queue) Index(number uint64) int {
	n := len(q)
	i := sort.Search(n, func(i int) bool { return q[i].number >= number })
	if i < n && q[i].number == number {
		return i
	}
	return -1
}
