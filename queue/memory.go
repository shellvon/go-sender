package queue

import (
	"container/heap"
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"github.com/shellvon/go-sender/core"
)

var (
	ErrQueueClosed = errors.New("queue is closed")
	ErrQueueFull   = errors.New("queue is full")
)

// MemoryQueue is a generic in-memory queue implementation.
type MemoryQueue[T core.Comparable[T]] struct {
	items      *PriorityQueue[T]
	mu         sync.RWMutex
	notifyChan chan struct{}
	closed     int32
	closeOnce  sync.Once
	maxSize    int
}

// NewMemoryQueue creates a new in-memory queue with the specified maximum size.
func NewMemoryQueue[T core.Comparable[T]](maxSize int) *MemoryQueue[T] {
	pq := &PriorityQueue[T]{}
	heap.Init(pq)

	return &MemoryQueue[T]{
		items:      pq,
		maxSize:    maxSize,
		notifyChan: make(chan struct{}, 1),
	}
}

func (mq *MemoryQueue[T]) Enqueue(ctx context.Context, item T) error {
	if atomic.LoadInt32(&mq.closed) == 1 {
		return ErrQueueClosed
	}

	mq.mu.Lock()
	defer mq.mu.Unlock()

	// Check capacity limit
	if mq.maxSize > 0 && mq.items.Len() >= mq.maxSize {
		return ErrQueueFull
	}

	heap.Push(mq.items, item)

	// Non-blocking notification
	select {
	case mq.notifyChan <- struct{}{}:
	default:
	}

	return nil
}

// Size returns the current number of elements in the queue.
func (mq *MemoryQueue[T]) Size() int {
	mq.mu.RLock()
	defer mq.mu.RUnlock()
	return mq.items.Len()
}

// MaxSize returns the maximum capacity of the queue, 0 means unlimited.
func (mq *MemoryQueue[T]) MaxSize() int {
	return mq.maxSize
}

// IsFull checks if the queue is full.
func (mq *MemoryQueue[T]) IsFull() bool {
	if mq.maxSize <= 0 {
		return false // Unlimited queue is never full
	}
	mq.mu.RLock()
	defer mq.mu.RUnlock()
	return mq.items.Len() >= mq.maxSize
}

// IsEmpty checks if the queue is empty.
func (mq *MemoryQueue[T]) IsEmpty() bool {
	mq.mu.RLock()
	defer mq.mu.RUnlock()
	return mq.items.Len() == 0
}

func (mq *MemoryQueue[T]) EnqueueDelayed(ctx context.Context, item T, delay time.Duration) error {
	if atomic.LoadInt32(&mq.closed) == 1 {
		return errors.New("queue is closed")
	}

	// If item implements Schedulable, set the scheduled time
	if schedulable, ok := any(item).(core.Schedulable); ok {
		scheduledTime := time.Now().Add(delay)
		schedulable.SetScheduledAt(scheduledTime)
	}

	return mq.Enqueue(ctx, item)
}

func (mq *MemoryQueue[T]) Dequeue(ctx context.Context) (T, error) {
	var zero T

	for {
		if atomic.LoadInt32(&mq.closed) == 1 {
			return zero, ErrQueueClosed
		}

		mq.mu.Lock()
		if mq.items.Len() > 0 {
			// Check if the top element is ready
			top := (*mq.items)[0]
			if mq.isReady(top) {
				item := heap.Pop(mq.items).(T)
				mq.mu.Unlock()
				return item, nil
			}
		}
		mq.mu.Unlock()

		// Wait for new items or context cancellation
		select {
		case <-ctx.Done():
			return zero, ctx.Err()
		case <-mq.notifyChan:
			// Continue checking
		case <-time.After(100 * time.Millisecond):
			// Periodically check delayed items
		}
	}
}

func (mq *MemoryQueue[T]) Close() error {
	mq.closeOnce.Do(func() {
		atomic.StoreInt32(&mq.closed, 1)
		close(mq.notifyChan)
	})
	return nil
}

// isReady checks if an item is ready to be dequeued.
func (mq *MemoryQueue[T]) isReady(item T) bool {
	// If item implements Schedulable, check the scheduled time
	if schedulable, ok := any(item).(core.Schedulable); ok {
		scheduledAt := schedulable.GetScheduledAt()
		if scheduledAt != nil {
			return time.Now().After(*scheduledAt)
		}
	}
	return true // No scheduled time or not Schedulable, immediately available
}

// PriorityQueue is a priority queue implementation.
type PriorityQueue[T core.Comparable[T]] []T

func (pq PriorityQueue[T]) Len() int { return len(pq) }

func (pq PriorityQueue[T]) Less(i, j int) bool {
	return pq[i].Compare(pq[j])
}

func (pq PriorityQueue[T]) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue[T]) Push(x interface{}) {
	*pq = append(*pq, x.(T))
}

func (pq *PriorityQueue[T]) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}
