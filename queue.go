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

// Insert inserts the call c into the queue in its correct position. The
// complexity is O(log(n)) where n = len(q).
func (q *Queue) Insert(c Call) {
	*q = append(*q, c)
	q.up(len(*q) - 1)
}

// Remove removes the element at index i from the queue.
// The complexity is O(log(n)) where n = len(q).
func (q *Queue) Remove(i int) {
	n := len(*q) - 1
	if n != i {
		q.swap(i, n)
		q.down(i, n)
		q.up(i)
	}
	*q = (*q)[0:n]
}

// Find returns the index of the first call with the given number in the queue.
// It will return -1 if such a call is not present in the queue.
func (q Queue) Find(number uint64) int {
	n := len(q)
	i := sort.Search(n, func(i int) bool { return q[i].number <= number })
	if i < n && q[i].number == number {
		return i
	}
	return -1
}

func (q Queue) less(i, j int) bool {
	return q[i].number < q[j].number
}

func (q Queue) swap(i, j int) {
	q[i], q[j] = q[j], q[i]
}

func (q Queue) up(j int) {
	for {
		i := (j - 1) / 2 // parent
		if i == j || !q.less(j, i) {
			break
		}
		q.swap(i, j)
		j = i
	}
}

func (q Queue) down(i, n int) {
	for {
		j1 := 2*i + 1
		if j1 >= n || j1 < 0 { // j1 < 0 after int overflow
			break
		}
		j := j1 // left child
		if j2 := j1 + 1; j2 < n && !q.less(j1, j2) {
			j = j2 // = 2*i + 2  // right child
		}
		if !q.less(j, i) {
			break
		}
		q.swap(i, j)
		i = j
	}
}
