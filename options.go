package bankofthailand

import (
	"net/http"
	"time"
)

type Option func(*Client)

func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// WithBaseURL overrides the base URL for the Holidays endpoint only.
// Other endpoints (exchange rates, interest rates, statistics, etc.)
// use hard-coded base URLs from the Bank of Thailand API documentation.
func WithBaseURL(baseURL string) Option {
	return func(c *Client) {
		c.baseURL = baseURL
	}
}

func WithToken(token string) Option {
	return func(c *Client) {
		c.token = token
	}
}

func WithRateLimiter(limiter RateLimiter) Option {
	return func(c *Client) {
		c.rateLimiter = limiter
	}
}

func WithRetryPolicy(policy *RetryPolicy) Option {
	return func(c *Client) {
		c.retryPolicy = policy
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.httpClient.Timeout = timeout
	}
}

func WithConfigPath(path string) Option {
	return func(c *Client) {
		c.configPath = path
	}
}
