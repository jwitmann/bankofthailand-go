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
	"strings"
	"time"
)

const (
	defaultBaseURL = "https://gateway.api.bot.or.th/financial-institutions-holidays"
	defaultTimeout = 30 * time.Second
)

// Endpoint keys used for token and rate-limiter selection.
const (
	endpointOthers         = "others"
	endpointExchangeRates  = "exchange_rates"
	endpointInterestRates  = "interest_rates"
	endpointStatistics     = "statistics"
	endpointDebtSecurities = "debt_security_auction"
	endpointLicenseCheck   = "license_check"
)

// endpointPatterns maps endpoint keys to URL path substrings.
var endpointPatterns = map[string][]string{
	endpointOthers: {
		"financial-institutions-holidays",
	},
	endpointExchangeRates: {
		"Stat-ExchangeRate",
		"Stat-ReferenceRate",
		"Stat-SpotRate",
		"Stat-SwapPoint",
		"Stat-ThaiBahtImpliedInterestRate",
	},
	endpointInterestRates: {
		"PolicyRate",
		"BIBOR",
		"DepositRate",
		"LoanRate",
		"Stat-InterbankTransactionRate",
	},
	endpointStatistics: {
		"categorylist",
		"serieslist",
		"observations",
		"search-series",
	},
	endpointDebtSecurities: {
		"BondAuction",
	},
	endpointLicenseCheck: {
		"BotLicenseCheckAPI",
	},
}

type Client struct {
	httpClient       *http.Client
	baseURL          string
	appID            string
	token            string
	tokens           map[string]string
	configPath       string
	rateLimiter      RateLimiter
	endpointLimiters map[string]RateLimiter
	retryPolicy      *RetryPolicy
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

	if client.token == "" && client.tokens == nil {
		if token := os.Getenv("BOT_API_TOKEN"); token != "" {
			client.token = token
		} else if err := client.loadConfig(); err != nil {
			return nil, fmt.Errorf("failed to load credentials: %w", err)
		}
	}

	if client.rateLimiter == nil && client.endpointLimiters == nil {
		client.endpointLimiters = map[string]RateLimiter{
			endpointOthers:         NewRateLimiterForHolidays(),
			endpointExchangeRates:  NewRateLimiterForExchangeRates(),
			endpointInterestRates:  NewRateLimiterForInterestRates(),
			endpointStatistics:     NewRateLimiterForStatistics(),
			endpointDebtSecurities: NewRateLimiterForDebtSecurities(),
			endpointLicenseCheck:   NewRateLimiterForLicenseCheck(),
		}
	}

	if client.retryPolicy == nil {
		client.retryPolicy = DefaultRetryPolicy()
	}

	return client, nil
}

func (c *Client) loadConfig() error {
	data, err := os.ReadFile(c.configPath)
	if err != nil {
		return fmt.Errorf("no credentials found: create %s: %w", c.configPath, err)
	}

	var cfg struct {
		AppID  string            `json:"app_id"`
		Token  string            `json:"token"`
		Tokens map[string]string `json:"tokens"`
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}

	c.appID = cfg.AppID

	if len(cfg.Tokens) > 0 {
		c.tokens = cfg.Tokens
		c.token = cfg.Token
		return nil
	}

	if cfg.Token != "" {
		c.token = cfg.Token
		return nil
	}

	return fmt.Errorf("config missing tokens map or token field")
}

// endpointKeyForURL returns the endpoint key (others, exchange_rates, etc.)
// that matches the given URL path, or "" if no match.
func endpointKeyForURL(urlStr string) string {
	u, err := url.Parse(urlStr)
	if err != nil {
		return ""
	}
	pathLower := strings.ToLower(u.Path)
	for key, patterns := range endpointPatterns {
		for _, p := range patterns {
			if strings.Contains(pathLower, strings.ToLower(p)) {
				return key
			}
		}
	}
	return ""
}

// tokenForURL returns the token to use for the given URL.
// It prefers a per-endpoint token from the tokens map, falling back to the
// default token set via WithToken.
func (c *Client) tokenForURL(urlStr string) string {
	if key := endpointKeyForURL(urlStr); key != "" {
		if t, ok := c.tokens[key]; ok && t != "" {
			return t
		}
	}
	return c.token
}

// HasEndpointToken reports whether a per-endpoint token is configured for the
// given endpoint key (e.g. "exchange_rates", "interest_rates").
func (c *Client) HasEndpointToken(endpoint string) bool {
	if c.tokens == nil {
		return false
	}
	t, ok := c.tokens[endpoint]
	return ok && t != ""
}

// EndpointTokenKeys returns the configured per-endpoint token keys.
func (c *Client) EndpointTokenKeys() []string {
	if c.tokens == nil {
		return nil
	}
	keys := make([]string, 0, len(c.tokens))
	for k := range c.tokens {
		keys = append(keys, k)
	}
	return keys
}

func (c *Client) buildRequest(ctx context.Context, method, urlStr string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, urlStr, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	token := c.tokenForURL(urlStr)
	if token == "" {
		return nil, fmt.Errorf("no token available for %s", urlStr)
	}

	req.Header.Set("Authorization", token)
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

func (c *Client) requestGet(ctx context.Context, baseURL, path string, query url.Values, result interface{}) error {
	u, err := url.Parse(baseURL + path)
	if err != nil {
		return fmt.Errorf("invalid url: %w", err)
	}
	if query != nil {
		u.RawQuery = query.Encode()
	}

	resp, err := c.GetURL(ctx, u.String())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		return ErrNoContent
	}

	return json.NewDecoder(resp.Body).Decode(result)
}

// rateLimiterForURL returns the appropriate rate limiter for the given URL.
// If a custom rateLimiter was set via WithRateLimiter, that takes precedence.
// Otherwise, it selects based on the API path segment.
func (c *Client) rateLimiterForURL(urlStr string) RateLimiter {
	if c.rateLimiter != nil {
		return c.rateLimiter
	}
	key := endpointKeyForURL(urlStr)
	if key == "" {
		return nil
	}
	if limiter, ok := c.endpointLimiters[key]; ok {
		return limiter
	}
	return nil
}

func setQuery(v url.Values, key, value string) {
	if value != "" {
		v.Set(key, value)
	}
}

func getEndpoint[T any](ctx context.Context, c *Client, baseURL, path string, query url.Values, errMsg string) (*T, error) {
	var result T
	if err := c.requestGet(ctx, baseURL, path, query, &result); err != nil {
		if errMsg != "" {
			return nil, fmt.Errorf("%s: %w", errMsg, err)
		}
		return nil, err
	}
	return &result, nil
}

func (c *Client) do(req *http.Request) (*http.Response, error) {
	if limiter := c.rateLimiterForURL(req.URL.String()); limiter != nil {
		if err := limiter.Wait(req.Context()); err != nil {
			return nil, fmt.Errorf("rate limiter: %w", err)
		}
	}

	resp, err := c.doWithRetry(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, NewAPIError(resp)
	}

	return resp, nil
}

func (c *Client) doWithRetry(req *http.Request) (*http.Response, error) {
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

	return resp, nil
}
