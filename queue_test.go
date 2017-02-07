package calltracker

import "testing"

func TestPush(t *testing.T) {
	const items = 100
	var q Queue
	for i := uint64(0); i < items; i++ {
		q.Push(Call{number: i})
	}
	for i := uint64(0); i < items; i++ {
		if q[i].Number() != i {
			t.Errorf("Queue index %d contains call number %d. Expected %d", i, q[i].Number(), i)
		}
	}
}

func TestRemove(t *testing.T) {
	const items = 21
	removing := []int{7, 8, 1, 1, 12, 11, 14, 13, 0, 0, 1, 2, 0}
	expected := []uint64{6, 10, 11, 12, 13, 14, 17, 18}
	var q Queue
	for i := uint64(0); i < items; i++ {
		q.Push(Call{number: i})
	}
	for _, r := range removing {
		q.Remove(r)
	}
	if len(q) != len(expected) {
		t.Fatalf("Queue contains %d items when %d were expected.", len(q), len(expected))
	}
	for i, number := range expected {
		if q[i].Number() != number {
			t.Errorf("Queue index %d contains call number %d. Expected %d", i, q[i].Number(), number)
		}
	}
}
