package bankofthailand

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type RateLimiter interface {
	Wait(ctx context.Context) error
}

type TokenBucketRateLimiter struct {
	mu         sync.Mutex
	tokens     float64
	capacity   float64
	refillRate float64
	lastRefill time.Time
}

func NewTokenBucketRateLimiter(capacity int, refillPerSecond float64) *TokenBucketRateLimiter {
	return &TokenBucketRateLimiter{
		tokens:     float64(capacity),
		capacity:   float64(capacity),
		refillRate: refillPerSecond,
		lastRefill: time.Now(),
	}
}

func (r *TokenBucketRateLimiter) Wait(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(r.lastRefill).Seconds()
	r.tokens += elapsed * r.refillRate
	if r.tokens > r.capacity {
		r.tokens = r.capacity
	}
	r.lastRefill = now

	if r.tokens >= 1 {
		r.tokens--
		return nil
	}

	waitTime := time.Duration((1 - r.tokens) / r.refillRate * float64(time.Second))
	timer := time.NewTimer(waitTime)
	defer timer.Stop()

	select {
	case <-timer.C:
		r.tokens--
		if r.tokens < 0 {
			r.tokens = 0
		}
		return nil
	case <-ctx.Done():
		return fmt.Errorf("rate limiter wait cancelled: %w", ctx.Err())
	}
}

type NoOpRateLimiter struct{}

func (n *NoOpRateLimiter) Wait(ctx context.Context) error {
	return nil
}
