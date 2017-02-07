package calltracker

import (
	"sync"
	"time"
)

// Subscriber is a function that receivers tracker event notifications.
type Subscriber func(update Update)

// Update is an update to the tracker's value
type Update struct {
	Sequence uint64
	Value    Value
}

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
	mutex       sync.RWMutex
	queue       Queue
	next        uint64 // The next ID to hand out.
	sequence    uint64 // The next update sequence number to hand out.
	subscribers []Subscriber
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
	t.notify()
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

// Subscribe will add the given subscriber function as a listener that will
// receive updates whenever the tracker's value changes.
func (t *Tracker) Subscribe(s Subscriber) {
	t.mutex.Lock()
	t.subscribers = append(t.subscribers, s)
	t.mutex.Unlock()
}

func (t *Tracker) remove(number uint64) {
	t.mutex.Lock()
	i := t.queue.Index(number)
	if i >= 0 {
		t.queue.Remove(i)
		t.notify()
	}
	t.mutex.Unlock()
}

// notify will send an event to the subscribers of the tracker informing them
// of the value change.
//
// notify does not lock the tracker. It is the caller's responsibility to hold
// a lock on the tracker for the duration of the call.
func (t *Tracker) notify() {
	if len(t.subscribers) == 0 {
		return
	}

	subscribers := make([]Subscriber, len(t.subscribers))
	copy(subscribers, t.subscribers)

	var update Update
	update.Value = make(Value, len(t.queue))
	copy(update.Value, t.queue)

	for _, s := range subscribers {
		go s(update)
	}

	t.sequence++
}
