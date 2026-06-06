package bankofthailand

import (
	"errors"
	"fmt"
	"net/http"
)

// ErrNoContent is returned when the API returns HTTP 204 (No Content).
// This typically means the requested data is not yet available.
var ErrNoContent = errors.New("no content")

type APIError struct {
	StatusCode int
	Message    string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("bot api error: status=%d, message=%s", e.StatusCode, e.Message)
}

func NewAPIError(resp *http.Response) *APIError {
	return &APIError{
		StatusCode: resp.StatusCode,
		Message:    resp.Status,
	}
}
