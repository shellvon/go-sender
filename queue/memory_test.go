package queue_test

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/shellvon/go-sender/queue"
)

type testItem struct {
	id    int
	ready time.Time
}

func (t *testItem) Compare(other *testItem) bool {
	return t.id < other.id
}

func (t *testItem) SetScheduledAt(tm time.Time) { t.ready = tm }
func (t *testItem) GetScheduledAt() *time.Time  { return &t.ready }

func TestMemoryQueue_BasicFlow(t *testing.T) {
	q := queue.NewMemoryQueue[*testItem](10)
	item := &testItem{id: 1}
	err := q.Enqueue(context.Background(), item)
	if err != nil {
		t.Fatalf("enqueue failed: %v", err)
	}
	itemOut, err := q.Dequeue(context.Background())
	if err != nil || itemOut == nil || itemOut.id != item.id {
		t.Errorf("dequeue got %v, want %v, err=%v", itemOut, item, err)
	}
}

func TestMemoryQueue_FullAndEmpty(t *testing.T) {
	q := queue.NewMemoryQueue[*testItem](1)
	_ = q.Enqueue(context.Background(), &testItem{id: 1})
	err := q.Enqueue(context.Background(), &testItem{id: 2})
	if !errors.Is(err, queue.ErrQueueFull) {
		t.Errorf("expected ErrQueueFull, got %v", err)
	}
	_, _ = q.Dequeue(context.Background())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	_, err = q.Dequeue(ctx)
	if err == nil {
		t.Error("expected dequeue timeout error")
	}
}

func TestMemoryQueue_Closed(t *testing.T) {
	q := queue.NewMemoryQueue[*testItem](1)
	_ = q.Close()
	err := q.Enqueue(context.Background(), &testItem{id: 1})
	if !errors.Is(err, queue.ErrQueueClosed) {
		t.Errorf("expected ErrQueueClosed, got %v", err)
	}
	_, err = q.Dequeue(context.Background())
	if !errors.Is(err, queue.ErrQueueClosed) {
		t.Errorf("expected ErrQueueClosed on dequeue, got %v", err)
	}
	_ = q.Close() // 幂等
}

func TestMemoryQueue_EnqueueDelayed(t *testing.T) {
	q := queue.NewMemoryQueue[*testItem](2)
	item := &testItem{id: 1}
	err := q.EnqueueDelayed(context.Background(), item, 20*time.Millisecond)
	if err != nil {
		t.Fatalf("enqueue delayed failed: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	_, err = q.Dequeue(ctx)
	if err == nil {
		t.Error("should not dequeue before delay")
	}
	time.Sleep(15 * time.Millisecond)
	itemOut, err := q.Dequeue(context.Background())
	if err != nil || itemOut == nil || itemOut.id != item.id {
		t.Errorf("dequeue after delay failed: %v, %v", itemOut, err)
	}
}

func TestMemoryQueue_Concurrency(t *testing.T) {
	q := queue.NewMemoryQueue[*testItem](0)
	wg := sync.WaitGroup{}
	for i := range 50 {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			_ = q.Enqueue(context.Background(), &testItem{id: i})
		}(i)
	}
	wg.Wait()
	count := 0
	for !q.IsEmpty() {
		_, _ = q.Dequeue(context.Background())
		count++
	}
	if count != 50 {
		t.Errorf("expected 50 items, got %d", count)
	}
}

func TestMemoryQueue_Getters(t *testing.T) {
	q := queue.NewMemoryQueue[*testItem](2)
	if q.Size() != 0 {
		t.Error("new queue size should be 0")
	}
	if q.MaxSize() != 2 {
		t.Error("MaxSize not correct")
	}
	if q.IsFull() {
		t.Error("empty queue should not be full")
	}
	_ = q.Enqueue(context.Background(), &testItem{id: 1})
	_ = q.Enqueue(context.Background(), &testItem{id: 2})
	if !q.IsFull() {
		t.Error("queue should be full after enqueues")
	}
	_, _ = q.Dequeue(context.Background())
	if q.Size() != 1 {
		t.Error("size after dequeue should be 1")
	}
}
