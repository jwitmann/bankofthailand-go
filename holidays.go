package bankofthailand

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

func (c *Client) GetHolidays(ctx context.Context, year int) ([]Holiday, error) {
	query := url.Values{}
	query.Set("year", strconv.Itoa(year))

	var result HolidaysResponse
	if err := c.requestJSONBase(ctx, "/", query, &result); err != nil {
		return nil, fmt.Errorf("failed to get holidays: %w", err)
	}
	return result.Result.Data, nil
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
