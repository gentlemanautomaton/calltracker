package calltracker_test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/gentlemanautomaton/calltracker"
)

// In this example we spawn a large number of goroutines that are blocking on
// simluated calls of varying duration. While the calls are running we monitor
// the tracker to observe the backlog.
func Example() {
	const (
		// How many calls should we make in parallel?
		calls = 1000
		// What's the lower bound of a typical call duration?
		lower = int64(time.Millisecond * 20)
		// What's the upper bound of a typical call duration?
		upper = int64(time.Second * 4)
		// What's the maximum duration of a call?
		max = int64(time.Second * 8)
	)

	// Create trackers simply by declaring them.
	// The zero value is safe to use.
	// The tracker's methods are thread-safe.
	var t calltracker.Tracker

	// Simulate some number of simultaneous calls.
	go func() {
		fmt.Printf("Starting first call (1 of %d)\n", calls)
		for i := 0; i < calls; i++ {
			// Register the start of a new call with the tracker as early as possible.
			call := t.Add()
			// For this example, arbitrarily decide how long the call should take.
			d := time.Duration(rand.Int63n(upper-lower) + lower)
			// Occassionally thrown in a call that has a really long or short
			// duration.
			if rand.Intn(50) == 0 {
				d = time.Duration(rand.Int63n(max))
			}
			// Make the call in its own goroutine.
			go func() {
				// Notify the tracker when we're finished.
				defer call.Done()
				// Simulate some long-running process.
				time.Sleep(d)
			}()

			// Sleep a random amount after every call to stagger the rate
			time.Sleep(time.Duration(rand.Int63n(lower / 2)))
		}
		fmt.Printf("Starting last call (%d of %d)\n", calls, calls)
	}()

	// Monitor the call statistics until all of the calls have finished.
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		// Wait until it's time to report.
		<-ticker.C
		// Grab a snapshot of the current call backlog.
		v := t.Value()
		// Report values
		fmt.Printf("%d calls outstanding. Delay min: %v max: %v: combined: %v\n", v.Len(), v.Min(), v.Max(), v.Elapsed())
		// Check whether all calls have finished.
		if v.Len() == 0 {
			break
		}
	}

	fmt.Println("Done.")
}

func TestExample(t *testing.T) {
	if !testing.Verbose() {
		return
	}
	Example()
}
