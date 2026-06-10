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
		now := time.Now()
		elapsed := now.Sub(r.lastRefill).Seconds()
		r.tokens += elapsed * r.refillRate
		if r.tokens > r.capacity {
			r.tokens = r.capacity
		}
		r.lastRefill = now
		r.tokens--
		return nil
	case <-ctx.Done():
		return fmt.Errorf("rate limiter wait cancelled: %w", ctx.Err())
	}
}

type NoOpRateLimiter struct{}

func (n *NoOpRateLimiter) Wait(ctx context.Context) error {
	return nil
}

const (
	secondsPerHour = 3600

	RateLimitHolidays       = 100
	RateLimitExchangeRates  = 200
	RateLimitInterestRates  = 200
	RateLimitStatistics     = 2000
	RateLimitDebtSecurities = 200
	RateLimitLicenseCheck   = 100
)

type RateLimitInfo struct {
	CallsPerHour int
	Quota        string
}

func GetRateLimitInfo(endpoint string) RateLimitInfo {
	switch endpoint {
	case "holidays":
		return RateLimitInfo{CallsPerHour: RateLimitHolidays, Quota: "unlimited"}
	case "exchange_rates", "reference_rate", "spot_rate", "swap_point", "implied_rate":
		return RateLimitInfo{CallsPerHour: RateLimitExchangeRates, Quota: "unlimited"}
	case "policy_rate", "bibor", "deposit_rate", "loan_rate", "interbank_rate":
		return RateLimitInfo{CallsPerHour: RateLimitInterestRates, Quota: "unlimited"}
	case "category_list", "series_list", "observations", "search":
		return RateLimitInfo{CallsPerHour: RateLimitStatistics, Quota: "unlimited"}
	case "debt_security_auction":
		return RateLimitInfo{CallsPerHour: RateLimitDebtSecurities, Quota: "unlimited"}
	case "license_check":
		return RateLimitInfo{CallsPerHour: RateLimitLicenseCheck, Quota: "unlimited"}
	default:
		return RateLimitInfo{CallsPerHour: RateLimitHolidays, Quota: "unlimited"}
	}
}

func NewHourlyRateLimiter(callsPerHour int) *TokenBucketRateLimiter {
	refillRate := float64(callsPerHour) / secondsPerHour
	return NewTokenBucketRateLimiter(1, refillRate)
}

func NewRateLimiterForHolidays() *TokenBucketRateLimiter {
	return NewHourlyRateLimiter(RateLimitHolidays)
}

func NewRateLimiterForExchangeRates() *TokenBucketRateLimiter {
	return NewHourlyRateLimiter(RateLimitExchangeRates)
}

func NewRateLimiterForInterestRates() *TokenBucketRateLimiter {
	return NewHourlyRateLimiter(RateLimitInterestRates)
}

func NewRateLimiterForStatistics() *TokenBucketRateLimiter {
	return NewHourlyRateLimiter(RateLimitStatistics)
}

func NewRateLimiterForDebtSecurities() *TokenBucketRateLimiter {
	return NewHourlyRateLimiter(RateLimitDebtSecurities)
}

func NewRateLimiterForLicenseCheck() *TokenBucketRateLimiter {
	return NewHourlyRateLimiter(RateLimitLicenseCheck)
}
