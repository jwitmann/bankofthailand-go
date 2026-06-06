package bankofthailand

import (
	"net"
	"net/http"
	"time"
)

type RetryPolicy struct {
	MaxRetries    int
	BaseDelay     time.Duration
	MaxDelay      time.Duration
	Multiplier    float64
	RetryStatuses []int
}

func DefaultRetryPolicy() *RetryPolicy {
	return &RetryPolicy{
		MaxRetries:    3,
		BaseDelay:     500 * time.Millisecond,
		MaxDelay:      30 * time.Second,
		Multiplier:    2.0,
		RetryStatuses: []int{http.StatusTooManyRequests, http.StatusInternalServerError, http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout},
	}
}

func (r *RetryPolicy) Backoff(attempt int) time.Duration {
	delay := float64(r.BaseDelay) * pow(r.Multiplier, float64(attempt-1))
	if delay > float64(r.MaxDelay) {
		delay = float64(r.MaxDelay)
	}
	return time.Duration(delay)
}

func (r *RetryPolicy) ShouldRetry(err error) bool {
	if err == nil {
		return false
	}

	if _, ok := err.(net.Error); ok {
		return true
	}

	return false
}

func (r *RetryPolicy) ShouldRetryStatus(statusCode int) bool {
	for _, code := range r.RetryStatuses {
		if code == statusCode {
			return true
		}
	}
	return false
}

func pow(base, exp float64) float64 {
	result := 1.0
	for i := 0; i < int(exp); i++ {
		result *= base
	}
	return result
}
