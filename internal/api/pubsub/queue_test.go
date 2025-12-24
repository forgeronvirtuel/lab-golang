package pubsub

import "testing"

func TestNewQueue(t *testing.T) {
	q := NewQueue()
	if q == nil {
		t.Fatal("expected new queue, got nil")
	}
	if !q.IsEmpty() {
		t.Errorf("expected empty queue, got size %d", q.Size())
	}
	if q.Size() != 0 {
		t.Errorf("expected size 0, got %d", q.Size())
	}
}

func TestEnqueue(t *testing.T) {
	q := NewQueue()
	q.Enqueue([]byte("first"))
	if q.IsEmpty() {
		t.Errorf("expected non-empty queue after enqueue")
	}
	if q.Size() != 1 {
		t.Errorf("expected size 1 after one enqueue, got %d", q.Size())
	}

	q.Enqueue([]byte("second"))
	if q.Size() != 2 {
		t.Errorf("expected size 2 after two enqueues, got %d", q.Size())
	}
}

func TestDequeue(t *testing.T) {
	q := NewQueue()
	q.Enqueue([]byte("first"))
	q.Enqueue([]byte("second"))

	value, ok := q.Dequeue()
	if !ok {
		t.Errorf("expected successful dequeue")
	}
	if string(value) != "first" {
		t.Errorf("expected 'first', got '%s'", value)
	}
	if q.Size() != 1 {
		t.Errorf("expected size 1 after one dequeue, got %d", q.Size())
	}

	value, ok = q.Dequeue()
	if !ok {
		t.Errorf("expected successful dequeue")
	}
	if string(value) != "second" {
		t.Errorf("expected 'second', got '%s'", value)
	}
	if !q.IsEmpty() {
		t.Errorf("expected empty queue after dequeuing all elements")
	}

	_, ok = q.Dequeue()
	if ok {
		t.Errorf("expected unsuccessful dequeue from empty queue")
	}
}

func TestIsEmpty(t *testing.T) {
	q := NewQueue()
	if !q.IsEmpty() {
		t.Errorf("expected new queue to be empty")
	}

	q.Enqueue([]byte("item"))
	if q.IsEmpty() {
		t.Errorf("expected non-empty queue after enqueue")
	}

	q.Dequeue()
	if !q.IsEmpty() {
		t.Errorf("expected empty queue after dequeueing all items")
	}
}

func TestSize(t *testing.T) {
	q := NewQueue()
	if q.Size() != 0 {
		t.Errorf("expected size 0 for new queue, got %d", q.Size())
	}

	q.Enqueue([]byte("item1"))
	if q.Size() != 1 {
		t.Errorf("expected size 1 after one enqueue, got %d", q.Size())
	}

	q.Enqueue([]byte("item2"))
	if q.Size() != 2 {
		t.Errorf("expected size 2 after two enqueues, got %d", q.Size())
	}

	q.Dequeue()
	if q.Size() != 1 {
		t.Errorf("expected size 1 after one dequeue, got %d", q.Size())
	}

	q.Dequeue()
	if q.Size() != 0 {
		t.Errorf("expected size 0 after dequeuing all items, got %d", q.Size())
	}
}
