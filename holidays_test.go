package bankofthailand

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetHolidays(t *testing.T) {
	mockResponse := HolidaysResponse{
		Result: HolidaysResult{
			API: "API_V2.FIHolidays",
			Data: []Holiday{
				{
					Date:                   "2026-01-01",
					DateThai:               "01/01/2569",
					HolidayWeekDay:         "Thursday",
					HolidayWeekDayThai:     "วันพฤหัสบดี",
					HolidayDescription:     "New Year's Day",
					HolidayDescriptionThai: "วันขึ้นปีใหม่",
				},
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			t.Errorf("expected path /, got %s", r.URL.Path)
		}
		if r.URL.Query().Get("year") != "2026" {
			t.Errorf("expected year=2026, got %s", r.URL.Query().Get("year"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(mockResponse); err != nil {
			t.Errorf("failed to encode response: %v", err)
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

	resp, err := client.GetHolidays(context.Background(), 2026)
	if err != nil {
		t.Fatalf("GetHolidays failed: %v", err)
	}

	if len(resp) != 1 {
		t.Fatalf("expected 1 holiday, got %d", len(resp))
	}

	holiday := resp[0]
	if holiday.Date != "2026-01-01" {
		t.Errorf("expected date 2026-01-01, got %s", holiday.Date)
	}
	if holiday.HolidayDescription != "New Year's Day" {
		t.Errorf("expected description 'New Year's Day', got %s", holiday.HolidayDescription)
	}
}

func TestGetHolidays_APIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
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

	_, err = client.GetHolidays(context.Background(), 2026)
	if err == nil {
		t.Fatal("expected error for 404 response")
	}
}
