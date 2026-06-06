package bankofthailand

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

const (
	defaultBaseURL = "https://gateway.api.bot.or.th/financial-institutions-holidays"
	defaultTimeout = 30 * time.Second
)

type Client struct {
	httpClient  *http.Client
	baseURL     string
	appID       string
	token       string
	configPath  string
	rateLimiter RateLimiter
	retryPolicy *RetryPolicy
}

func NewClient(options ...Option) (*Client, error) {
	client := &Client{
		httpClient: &http.Client{Timeout: defaultTimeout},
		baseURL:    defaultBaseURL,
		configPath: filepath.Join("config", "bot-keys.json"),
	}

	for _, opt := range options {
		opt(client)
	}

	if client.token == "" {
		if err := client.loadConfig(); err != nil {
			return nil, fmt.Errorf("failed to load credentials: %w", err)
		}
	}

	if client.rateLimiter == nil {
		client.rateLimiter = NewTokenBucketRateLimiter(5, 1)
	}

	if client.retryPolicy == nil {
		client.retryPolicy = DefaultRetryPolicy()
	}

	return client, nil
}

func (c *Client) loadConfig() error {
	token := os.Getenv("BOT_API_TOKEN")

	if token != "" {
		c.token = token
		return nil
	}

	data, err := os.ReadFile(c.configPath)
	if err != nil {
		return fmt.Errorf("no credentials found: set BOT_API_TOKEN env var or create %s: %w", c.configPath, err)
	}

	var cfg struct {
		AppID string `json:"app_id"`
		Token string `json:"token"`
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}

	if cfg.Token == "" {
		return fmt.Errorf("config missing token")
	}

	c.appID = cfg.AppID
	c.token = cfg.Token
	return nil
}

func (c *Client) buildRequest(ctx context.Context, method, urlStr string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, urlStr, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", c.token)
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return req, nil
}

func (c *Client) newRequest(ctx context.Context, method, path string, query url.Values, body io.Reader) (*http.Request, error) {
	u, err := url.Parse(c.baseURL + path)
	if err != nil {
		return nil, fmt.Errorf("invalid url: %w", err)
	}

	if query != nil {
		u.RawQuery = query.Encode()
	}

	return c.buildRequest(ctx, method, u.String(), body)
}

func (c *Client) Get(ctx context.Context, path string, query url.Values) (*http.Response, error) {
	req, err := c.newRequest(ctx, http.MethodGet, path, query, nil)
	if err != nil {
		return nil, err
	}
	return c.do(req)
}

func (c *Client) GetURL(ctx context.Context, urlStr string) (*http.Response, error) {
	req, err := c.buildRequest(ctx, http.MethodGet, urlStr, nil)
	if err != nil {
		return nil, err
	}
	return c.do(req)
}

func (c *Client) do(req *http.Request) (*http.Response, error) {
	if c.rateLimiter != nil {
		if err := c.rateLimiter.Wait(req.Context()); err != nil {
			return nil, fmt.Errorf("rate limiter: %w", err)
		}
	}

	var resp *http.Response
	var err error

	attempts := c.retryPolicy.MaxRetries + 1
	for i := 0; i < attempts; i++ {
		if i > 0 {
			time.Sleep(c.retryPolicy.Backoff(i))
		}

		resp, err = c.httpClient.Do(req)
		if err != nil {
			if !c.retryPolicy.ShouldRetry(err) {
				return nil, err
			}
			continue
		}

		if resp.StatusCode < 500 {
			break
		}

		resp.Body.Close()

		if !c.retryPolicy.ShouldRetryStatus(resp.StatusCode) {
			return nil, NewAPIError(resp)
		}
	}

	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, NewAPIError(resp)
	}

	return resp, nil
}
