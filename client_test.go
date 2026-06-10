package bankofthailand

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestNewClient_WithToken(t *testing.T) {
	client, err := NewClient(
		WithToken("test-token"),
		WithBaseURL("https://example.com"),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	if client.token != "test-token" {
		t.Errorf("expected token=test-token, got %s", client.token)
	}
}

func TestNewClient_WithEnvVar(t *testing.T) {
	os.Setenv("BOT_API_TOKEN", "env-token")
	defer os.Unsetenv("BOT_API_TOKEN")

	client, err := NewClient(WithBaseURL("https://example.com"))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	if client.token != "env-token" {
		t.Errorf("expected token=env-token, got %s", client.token)
	}
}

func TestNewClient_MissingCredentials(t *testing.T) {
	os.Unsetenv("BOT_API_TOKEN")

	_, err := NewClient(
		WithBaseURL("https://example.com"),
		WithConfigPath("/nonexistent/config.json"),
	)
	if err == nil {
		t.Fatal("expected error for missing credentials")
	}
}

func TestClient_Get_WithAuthHeader(t *testing.T) {
	var receivedAuth string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedAuth = r.Header.Get("Authorization")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(`[{"Date":"2026-01-01"}]`)); err != nil {
			t.Errorf("failed to write response: %v", err)
		}
	}))
	defer server.Close()

	client, err := NewClient(
		WithToken("test-token"),
		WithBaseURL(server.URL),
		WithRateLimiter(&NoOpRateLimiter{}),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	resp, err := client.Get(context.Background(), "/", nil)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if receivedAuth != "test-token" {
		t.Errorf("expected Authorization=test-token, got %s", receivedAuth)
	}
}

func TestClient_Get_RetryOnServerError(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts < 3 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(`[{"Date":"2026-01-01"}]`)); err != nil {
			t.Errorf("failed to write response: %v", err)
		}
	}))
	defer server.Close()

	client, err := NewClient(
		WithToken("test"),
		WithBaseURL(server.URL),
		WithRateLimiter(&NoOpRateLimiter{}),
		WithRetryPolicy(&RetryPolicy{
			MaxRetries:    5,
			BaseDelay:     10 * time.Millisecond,
			MaxDelay:      100 * time.Millisecond,
			Multiplier:    2.0,
			RetryStatuses: []int{http.StatusServiceUnavailable},
		}),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	resp, err := client.Get(context.Background(), "/", nil)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if attempts != 3 {
		t.Errorf("expected 3 attempts, got %d", attempts)
	}
}

func TestClient_Get_APIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		if _, err := w.Write([]byte(`{"error":"unauthorized"}`)); err != nil {
			t.Errorf("failed to write response: %v", err)
		}
	}))
	defer server.Close()

	client, err := NewClient(
		WithToken("test"),
		WithBaseURL(server.URL),
		WithRateLimiter(&NoOpRateLimiter{}),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	_, err = client.Get(context.Background(), "/", nil)
	if err == nil {
		t.Fatal("expected error for 401 response")
	}

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", apiErr.StatusCode)
	}
}

func TestRateLimiter_Wait(t *testing.T) {
	limiter := NewTokenBucketRateLimiter(2, 10)
	ctx := context.Background()

	if err := limiter.Wait(ctx); err != nil {
		t.Fatalf("first wait failed: %v", err)
	}
	if err := limiter.Wait(ctx); err != nil {
		t.Fatalf("second wait failed: %v", err)
	}

	start := time.Now()
	if err := limiter.Wait(ctx); err != nil {
		t.Fatalf("third wait failed: %v", err)
	}
	elapsed := time.Since(start)

	if elapsed < 50*time.Millisecond {
		t.Errorf("expected some delay for third request, got %v", elapsed)
	}
}

func TestParseHolidayYear(t *testing.T) {
	tests := []struct {
		input    string
		expected int
		wantErr  bool
	}{
		{"2026", 2026, false},
		{"2000", 2000, false},
		{"2100", 2100, false},
		{"1999", 0, true},
		{"2101", 0, true},
		{"abc", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := ParseHolidayYear(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseHolidayYear(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if got != tt.expected {
				t.Errorf("ParseHolidayYear(%q) = %d, want %d", tt.input, got, tt.expected)
			}
		})
	}
}

func TestGetRateLimitInfo(t *testing.T) {
	tests := []struct {
		endpoint string
		want     RateLimitInfo
	}{
		{"holidays", RateLimitInfo{CallsPerHour: 100, Quota: "unlimited"}},
		{"exchange_rates", RateLimitInfo{CallsPerHour: 200, Quota: "unlimited"}},
		{"reference_rate", RateLimitInfo{CallsPerHour: 200, Quota: "unlimited"}},
		{"spot_rate", RateLimitInfo{CallsPerHour: 200, Quota: "unlimited"}},
		{"swap_point", RateLimitInfo{CallsPerHour: 200, Quota: "unlimited"}},
		{"implied_rate", RateLimitInfo{CallsPerHour: 200, Quota: "unlimited"}},
		{"policy_rate", RateLimitInfo{CallsPerHour: 200, Quota: "unlimited"}},
		{"bibor", RateLimitInfo{CallsPerHour: 200, Quota: "unlimited"}},
		{"deposit_rate", RateLimitInfo{CallsPerHour: 200, Quota: "unlimited"}},
		{"loan_rate", RateLimitInfo{CallsPerHour: 200, Quota: "unlimited"}},
		{"interbank_rate", RateLimitInfo{CallsPerHour: 200, Quota: "unlimited"}},
		{"category_list", RateLimitInfo{CallsPerHour: 2000, Quota: "unlimited"}},
		{"series_list", RateLimitInfo{CallsPerHour: 2000, Quota: "unlimited"}},
		{"observations", RateLimitInfo{CallsPerHour: 2000, Quota: "unlimited"}},
		{"search", RateLimitInfo{CallsPerHour: 2000, Quota: "unlimited"}},
		{"debt_security_auction", RateLimitInfo{CallsPerHour: 200, Quota: "unlimited"}},
		{"license_check", RateLimitInfo{CallsPerHour: 100, Quota: "unlimited"}},
		{"unknown", RateLimitInfo{CallsPerHour: 100, Quota: "unlimited"}},
	}

	for _, tt := range tests {
		t.Run(tt.endpoint, func(t *testing.T) {
			got := GetRateLimitInfo(tt.endpoint)
			if got.CallsPerHour != tt.want.CallsPerHour {
				t.Errorf("GetRateLimitInfo(%q).CallsPerHour = %d, want %d", tt.endpoint, got.CallsPerHour, tt.want.CallsPerHour)
			}
			if got.Quota != tt.want.Quota {
				t.Errorf("GetRateLimitInfo(%q).Quota = %s, want %s", tt.endpoint, got.Quota, tt.want.Quota)
			}
		})
	}
}

func TestRequestGet_NoContent(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client, err := NewClient(
		WithToken("test"),
		WithBaseURL(server.URL),
		WithRateLimiter(&NoOpRateLimiter{}),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	var result struct{}
	err = client.requestGet(context.Background(), server.URL, "/", nil, &result)
	if !errors.Is(err, ErrNoContent) {
		t.Errorf("expected ErrNoContent, got %v", err)
	}
}

func TestRequestGet_InvalidURL(t *testing.T) {
	client, err := NewClient(
		WithToken("test"),
		WithBaseURL("https://example.com"),
		WithRateLimiter(&NoOpRateLimiter{}),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	var result struct{}
	err = client.requestGet(context.Background(), "://invalid-url", "/", nil, &result)
	if err == nil {
		t.Fatal("expected error for invalid URL")
	}
}

func TestHourlyRateLimiter(t *testing.T) {
	limiter := NewTokenBucketRateLimiter(10, 1000)
	ctx := context.Background()

	if err := limiter.Wait(ctx); err != nil {
		t.Fatalf("first wait failed: %v", err)
	}

	start := time.Now()
	if err := limiter.Wait(ctx); err != nil {
		t.Fatalf("second wait failed: %v", err)
	}
	elapsed := time.Since(start)

	if elapsed > 2*time.Millisecond {
		t.Errorf("expected fast delay with capacity 10, got %v", elapsed)
	}
}
