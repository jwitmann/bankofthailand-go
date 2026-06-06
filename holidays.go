package bankofthailand

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

func (c *Client) GetHolidays(ctx context.Context, year int) ([]Holiday, error) {
	query := url.Values{}
	query.Set("year", strconv.Itoa(year))

	resp, err := c.Get(ctx, "/", query)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result []Holiday
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode holidays response: %w", err)
	}

	return result, nil
}

func (c *Client) GetHolidaysThaiFA(ctx context.Context, year int) (*ThaiFAResponse, error) {
	holidays, err := c.GetHolidays(ctx, year)
	if err != nil {
		return nil, err
	}

	thaifa := &ThaiFAResponse{
		Result: struct {
			API       string          `json:"api"`
			Timestamp string          `json:"timestamp"`
			Data      []ThaiFAHoliday `json:"data"`
		}{
			API:       "API_V2.FIHolidays",
			Timestamp: "",
			Data:      make([]ThaiFAHoliday, len(holidays)),
		},
	}

	for i, h := range holidays {
		thaifa.Result.Data[i] = ThaiFAHoliday{
			HolidayWeekDay:         h.HolidayWeekDay,
			HolidayWeekDayThai:     h.HolidayWeekDayThai,
			Date:                   h.Date,
			DateThai:               h.DateThai,
			HolidayDescription:     h.HolidayDescription,
			HolidayDescriptionThai: h.HolidayDescriptionThai,
		}
	}

	return thaifa, nil
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
