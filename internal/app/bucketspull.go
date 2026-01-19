package app

import (
	"sync"

	"github.com/MaximBayurov/rate-limiter/internal/bucket"
)

type bucketsPull struct {
	capacity int
	buckets  map[string]bucket.Bucket
	mx       sync.Mutex
}

func newBucketsPull(capacity int) bucketsPull {
	return bucketsPull{
		capacity: capacity,
		buckets:  make(map[string]bucket.Bucket, 0),
		mx:       sync.Mutex{},
	}
}

func (p *bucketsPull) ClearEmptyBuckets() {
	p.mx.Lock()
	defer p.mx.Unlock()

	var key string
	var b bucket.Bucket
	for key, b = range p.buckets {
		if b.CurrentLoad() == 0 {
			delete(p.buckets, key)
		}
	}
}

func (p *bucketsPull) Allow(key string) bool {
	p.mx.Lock()
	defer p.mx.Unlock()

	var b bucket.Bucket
	var ok bool
	if b, ok = p.buckets[key]; !ok {
		b = bucket.NewLeakyBucket(p.capacity)
		p.buckets[key] = b
	}
	return b.Allow()
}

func (p *bucketsPull) DeleteBucket(key string) {
	p.mx.Lock()
	defer p.mx.Unlock()

	if _, ok := p.buckets[key]; !ok {
		return
	}
	delete(p.buckets, key)
}
