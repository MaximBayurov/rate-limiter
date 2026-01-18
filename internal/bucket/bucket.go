package bucket

import (
	"container/list"
	"sync"
	"time"
)

type Bucket interface {
	Allow() bool
	CurrentLoad() int
}

type LeakyBucket struct {
	capacity int
	rate     time.Duration
	requests *list.List
	mu       sync.Mutex
}

func NewLeakyBucket(capacity int) *LeakyBucket {
	return &LeakyBucket{
		capacity: capacity,
		rate:     time.Minute,
		requests: list.New(),
	}
}

func (lb *LeakyBucket) Allow() bool {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	now := time.Now()
	lb.cleanup(now)

	if lb.requests.Len() >= lb.capacity {
		return false
	}

	lb.requests.PushBack(now)
	return true
}

func (lb *LeakyBucket) CurrentLoad() int {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	now := time.Now()
	lb.cleanup(now)

	return lb.requests.Len()
}

func (lb *LeakyBucket) cleanup(now time.Time) {
	for lb.requests.Len() > 0 {
		front := lb.requests.Front()
		if now.Sub(front.Value.(time.Time)) > lb.rate {
			lb.requests.Remove(front)
		} else {
			break
		}
	}
}
