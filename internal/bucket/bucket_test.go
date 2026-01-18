package bucket

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLeakyBucket_Allow(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name       string
		capacity   int
		attempts   int
		allowed    int
		disallowed int
		sleep      time.Duration
	}{
		{
			name:       "Непрерывно (без time.Sleep)",
			capacity:   10,
			attempts:   15,
			allowed:    10,
			disallowed: 5,
			sleep:      time.Duration(0),
		},
		{
			name:       "Успевают сгореть (с отказами)",
			capacity:   5,
			attempts:   12,
			allowed:    7,
			disallowed: 5,
			sleep:      6 * time.Second,
		},
		{
			name:       "Не успевают сгореть (с отказами)",
			capacity:   5,
			attempts:   10,
			allowed:    5,
			disallowed: 5,
			sleep:      6 * time.Second,
		},
		{
			name:       "Успевают сгорать без отказов",
			capacity:   10,
			attempts:   15,
			allowed:    15,
			disallowed: 0,
			sleep:      6 * time.Second,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			bucket := NewLeakyBucket(tc.capacity)

			var allowed, disallowed int
			for i := 1; i <= tc.attempts; i++ {
				ok := bucket.Allow()
				if ok {
					allowed++
				} else {
					disallowed++
				}
				time.Sleep(tc.sleep)
			}
			assert.Equal(t, tc.allowed, allowed, "allowed")
			assert.Equal(t, tc.disallowed, disallowed, "disallowed")
		})
	}
}

func TestLeakyBucket_Allow_Batch(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name string
		capacity,
		first,
		second,
		allowed,
		disallowed int
		interval time.Duration
	}{
		{
			name:       "Полностью заполняем сразу",
			capacity:   10,
			first:      10,
			second:     0,
			allowed:    10,
			disallowed: 0,
			interval:   time.Minute,
		},
		{
			name:       "Заполняем на половину",
			capacity:   10,
			first:      5,
			second:     5,
			allowed:    10,
			disallowed: 0,
			interval:   time.Minute,
		},
		{
			name:       "Заполняем на половину (с отказами)",
			capacity:   10,
			first:      5,
			second:     7,
			allowed:    10,
			disallowed: 2,
			interval:   time.Second * 30,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			bucket := NewLeakyBucket(tc.capacity)

			var allowed, disallowed int
			for i := 1; i <= tc.first; i++ {
				ok := bucket.Allow()
				if ok {
					allowed++
				} else {
					disallowed++
				}
			}
			time.Sleep(tc.interval)
			for i := 1; i <= tc.second; i++ {
				ok := bucket.Allow()
				if ok {
					allowed++
				} else {
					disallowed++
				}
			}
			assert.Equal(t, tc.allowed, allowed, "allowed")
			assert.Equal(t, tc.disallowed, disallowed, "disallowed")
		})
	}
}

func TestLeakyBucket_Allow_BatchWithSleep(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name       string
		capacity   int
		allowed    int
		disallowed int
		sleep      time.Duration
		first      int
		second     int
		interval   time.Duration
	}{
		{
			name:       "С отказами",
			capacity:   10,
			allowed:    13,
			disallowed: 7,
			sleep:      3 * time.Second,
			first:      10,
			second:     10,
			interval:   10 * time.Second,
		},
		{
			name:       "Без отказов",
			capacity:   10,
			allowed:    20,
			disallowed: 0,
			sleep:      3 * time.Second,
			first:      10,
			second:     10,
			interval:   30 * time.Second,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			bucket := NewLeakyBucket(tc.capacity)

			var allowed, disallowed int
			for i := 1; i <= tc.first; i++ {
				ok := bucket.Allow()
				if ok {
					allowed++
				} else {
					disallowed++
				}
				time.Sleep(tc.sleep)
			}
			time.Sleep(tc.interval)
			for i := 1; i <= tc.second; i++ {
				ok := bucket.Allow()
				if ok {
					allowed++
				} else {
					disallowed++
				}
				time.Sleep(tc.sleep)
			}
			assert.Equal(t, tc.allowed, allowed, "allowed")
			assert.Equal(t, tc.disallowed, disallowed, "disallowed")
		})
	}
}
