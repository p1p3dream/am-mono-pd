package distsync

import (
	"math/rand/v2"
	"sync"
	"time"

	"abodemine/lib/errors"
)

// Backoff implements a concurrency-safe
// exponential backoff strategy with jitter.
type Backoff struct {
	InitialInterval time.Duration
	JitterFactor    float64
	MaxInterval     time.Duration
	MaxRetries      int

	retries int
	mu      sync.Mutex
}

// Next returns the duration to wait before the next retry.
func (b *Backoff) Next() (time.Duration, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.retries++

	if b.retries > b.MaxRetries {
		return 0, &errors.Object{
			Id:     "b0c1f2a4-3d8e-4b8e-9f79-ebf9def3e131",
			Code:   errors.Code_RESOURCE_EXHAUSTED,
			Detail: "Max retries exceeded.",
		}
	}

	delay := b.InitialInterval * time.Duration(1<<b.retries)

	if delay > b.MaxInterval {
		return b.MaxInterval, nil
	}

	if b.JitterFactor < 0 {
		return 0, &errors.Object{
			Id:     "907e8d95-1127-4a8f-af37-026a9741fe73",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Jitter factor must be greater than 0.",
		}
	}

	jitterFactor := b.JitterFactor

	if jitterFactor == 0 {
		jitterFactor = 0.4
	}

	jitter := time.Duration(rand.Float64()*jitterFactor*float64(delay)) -
		time.Duration(jitterFactor/2*float64(delay))

	delay += jitter

	// Check again to prevent exceeding max interval.
	if delay > b.MaxInterval {
		return b.MaxInterval, nil
	}

	return delay, nil
}

// Retries returns the number of retries that have been attempted.
func (b *Backoff) Retries() int {
	b.mu.Lock()
	v := b.retries
	b.mu.Unlock()

	return v
}

// Reset resets the backoff to the initial state.
func (b *Backoff) Reset() {
	b.mu.Lock()
	b.retries = 0
	b.mu.Unlock()
}

// Sleep sleeps for the duration calculated by Next().
func (b *Backoff) Sleep() error {
	delay, err := b.Next()
	if err != nil {
		return errors.Forward(err, "e4767739-af03-42e6-9b0b-2f9c1b94d950")
	}
	time.Sleep(delay)
	return nil
}
