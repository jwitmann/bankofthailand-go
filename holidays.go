package bankofthailand

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strconv"
)

func (c *Client) GetHolidays(ctx context.Context, year int) ([]Holiday, error) {
	resp, err := c.GetHolidaysRaw(ctx, year)
	if err != nil {
		return nil, err
	}
	return resp.Result.Data, nil
}

func (c *Client) GetHolidaysRaw(ctx context.Context, year int) (*HolidaysResponse, error) {
	query := url.Values{}
	query.Set("year", strconv.Itoa(year))

	var result HolidaysResponse
	if err := c.requestGet(ctx, c.baseURL, "/", query, &result); err != nil {
		if errors.Is(err, ErrNoContent) {
			return &HolidaysResponse{Result: HolidaysResult{Data: []Holiday{}}}, nil
		}
		return nil, fmt.Errorf("failed to get holidays: %w", err)
	}
	return &result, nil
}

func ParseHolidayYear(s string) (int, error) {
	year, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("invalid year: %w", err)
	}
	if year < 2000 || year > 2100 {
		return 0, fmt.Errorf("year must be between 2000 and 2100")
	}
	return year, nil
}
