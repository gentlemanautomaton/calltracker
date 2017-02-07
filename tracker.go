package calltracker

import (
	"fmt"
	"sync"
	"time"
)

// TrackedCall is a call for which completion can be signaled.
type TrackedCall struct {
	Call
	t *Tracker
}

// Done should be called when a tracked call has been completed.
//
// It is safe to call Done() more than once, but it is not safe to call it
// simulataneously from separate goroutines.
func (tc *TrackedCall) Done() {
	if tc.t != nil {
		tc.t.remove(tc.number)
		tc.t = nil
	}
}

// Tracker tracks a set of outstanding calls. Its zero value is safe for use.
// It is safe to call its functions from more than one goroutine
// simulataneously.
type Tracker struct {
	mutex sync.RWMutex
	queue Queue
	next  uint64 // The next ID to hand out.
}

// Add will add a call to the tracker. The call's start time will be the time
// at which Add() was called.
//
// The returned value is a tracked call that allows the caller to indicate when
// the call has finished. The caller must call its Done() function to remove the
// call from the tracker's list of tracked calls.
func (t *Tracker) Add() (call TrackedCall) {
	when := time.Now()

	t.mutex.Lock()
	call = TrackedCall{
		Call: Call{
			number: t.next,
			when:   when,
		},
		t: t,
	}
	t.queue.Push(call.Call)
	t.next++
	t.mutex.Unlock()

	return
}

// Value returns a point-in-time view of the tracked calls.
func (t *Tracker) Value() (value Value) {
	t.mutex.RLock()
	value = make(Value, len(t.queue))
	copy(value, t.queue)
	t.mutex.RUnlock()
	return
}

func (t *Tracker) remove(number uint64) {
	t.mutex.Lock()
	i := t.queue.Find(number)
	if i >= 0 {
		t.queue.Remove(i)
	}
	t.mutex.Unlock()
	fmt.Printf("Removed %v for number %d\n", i, number)
}
