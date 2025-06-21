package queue

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/shellvon/go-sender/core"
)

var _ core.Schedulable = (*TestItem)(nil)

type TestItem struct {
	value       string
	scheduledAt *time.Time
	priority    int
}

func (t *TestItem) Compare(other *TestItem) bool {
	return t.priority < other.priority
}

func (t *TestItem) GetScheduledAt() *time.Time {
	return t.scheduledAt
}

func (t *TestItem) SetScheduledAt(ti time.Time) {
	t.scheduledAt = &ti
}

func TestNewMemoryQueue(t *testing.T) {
	queue := NewMemoryQueue[*TestItem](10)
	if queue == nil {
		t.Fatal("NewMemoryQueue returned nil")
	}
	if queue.MaxSize() != 10 {
		t.Errorf("Expected max size 10, got %d", queue.MaxSize())
	}
	if !queue.IsEmpty() {
		t.Error("New queue should be empty")
	}
}

func TestEnqueueDequeue(t *testing.T) {
	queue := NewMemoryQueue[*TestItem](10)
	ctx := context.Background()

	item := &TestItem{value: "test", priority: 1}
	err := queue.Enqueue(ctx, item)
	if err != nil {
		t.Errorf("Enqueue failed: %v", err)
	}

	if queue.Size() != 1 {
		t.Errorf("Expected size 1, got %d", queue.Size())
	}

	dequeued, err := queue.Dequeue(ctx)
	if err != nil {
		t.Errorf("Dequeue failed: %v", err)
	}
	if dequeued.value != item.value {
		t.Errorf("Expected value %s, got %s", item.value, dequeued.value)
	}
}

func TestEnqueueDelayed(t *testing.T) {
	queue := NewMemoryQueue[*TestItem](10)
	ctx := context.Background()

	item := TestItem{value: "delayed", priority: 1}
	err := queue.EnqueueDelayed(ctx, &item, 100*time.Millisecond)
	if err != nil {
		t.Errorf("EnqueueDelayed failed: %v", err)
	}

	now := time.Now()
	dequeued, err := queue.Dequeue(context.Background())

	if err != nil {
		t.Errorf("Dequeue failed: %v", err)
	}
	if dequeued.value != item.value {
		t.Errorf("Expected value %s, got %s", item.value, dequeued.value)
	}

	if !dequeued.scheduledAt.After(now) {
		t.Errorf("dequeue too early")
	}
}

func TestQueueFull(t *testing.T) {
	queue := NewMemoryQueue[*TestItem](1)
	ctx := context.Background()

	err := queue.Enqueue(ctx, &TestItem{value: "test1", priority: 1})
	if err != nil {
		t.Errorf("Enqueue failed: %v", err)
	}

	err = queue.Enqueue(ctx, &TestItem{value: "test2", priority: 2})
	if !errors.Is(err, ErrQueueFull) {
		t.Errorf("Expected ErrQueueFull, got %v", err)
	}

	if !queue.IsFull() {
		t.Error("Queue should be full")
	}
}

func TestQueueClose(t *testing.T) {
	queue := NewMemoryQueue[*TestItem](10)
	ctx := context.Background()

	err := queue.Close()
	if err != nil {
		t.Errorf("Close failed: %v", err)
	}

	err = queue.Enqueue(ctx, &TestItem{value: "test", priority: 1})
	if !errors.Is(err, ErrQueueClosed) {
		t.Errorf("Expected ErrQueueClosed, got %v", err)
	}

	_, err = queue.Dequeue(ctx)
	if !errors.Is(err, ErrQueueClosed) {
		t.Errorf("Expected ErrQueueClosed, got %v", err)
	}
}

func TestPriorityQueue(t *testing.T) {
	queue := NewMemoryQueue[*TestItem](10)
	ctx := context.Background()

	items := []*TestItem{
		{value: "low", priority: 3},
		{value: "high", priority: 1},
		{value: "medium", priority: 2},
	}

	for _, item := range items {
		err := queue.Enqueue(ctx, item)
		if err != nil {
			t.Errorf("Enqueue failed: %v", err)
		}
	}

	// Should dequeue in priority order: high, medium, low
	expectedOrder := []string{"high", "medium", "low"}
	for i, expected := range expectedOrder {
		item, err := queue.Dequeue(ctx)
		if err != nil {
			t.Errorf("Dequeue failed at step %d: %v", i, err)
		}
		if item.value != expected {
			t.Errorf("Expected value %s, got %s at step %d", expected, item.value, i)
		}
	}
}
