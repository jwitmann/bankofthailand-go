package bankofthailand

import (
	"context"
	"fmt"
	"net/url"
)

const (
	categoryListBaseURL = "https://gateway.api.bot.or.th/categorylist"
	observationsBaseURL = "https://gateway.api.bot.or.th/observations"
	searchBaseURL       = "https://gateway.api.bot.or.th/search-series"
)

type CategoryListResponse struct {
	Result struct {
		API       string     `json:"api"`
		Timestamp string     `json:"timestamp"`
		Category  []Category `json:"category"`
	} `json:"result"`
}

type Category struct {
	Category       string `json:"category"`
	DescriptionTh  string `json:"description_th"`
	DescriptionEng string `json:"description_eng"`
}

type SeriesListResponse struct {
	Result struct {
		API       string   `json:"api"`
		Timestamp string   `json:"timestamp"`
		Series    []Series `json:"series"`
	} `json:"result"`
}

type Series struct {
	Category         string `json:"category"`
	SeriesCode       string `json:"series_code"`
	SeriesNameTh     string `json:"series_name_th"`
	SeriesNameEng    string `json:"series_name_eng"`
	ObservationStart string `json:"observation_start"`
	ObservationEnd   string `json:"observation_end"`
	LastUpdateDate   string `json:"last_update_date"`
}

type ObservationsResponse struct {
	Result struct {
		API       string              `json:"api"`
		Timestamp string              `json:"timestamp"`
		Series    []ObservationSeries `json:"series"`
	} `json:"result"`
}

type ObservationSeries struct {
	SeriesCode     string        `json:"series_code"`
	SeriesNameTh   string        `json:"series_name_th"`
	SeriesNameEng  string        `json:"series_name_eng"`
	UnitTh         string        `json:"unit_th"`
	UnitEng        string        `json:"unit_eng"`
	SeriesType     string        `json:"series_type"`
	Frequency      string        `json:"frequency"`
	LastUpdateDate string        `json:"last_update_date"`
	Observations   []Observation `json:"observations"`
}

type Observation struct {
	PeriodStart string `json:"period_start"`
	Value       string `json:"value"`
}

type SearchResponse struct {
	Result struct {
		API           string         `json:"api"`
		Timestamp     string         `json:"timestamp"`
		SeriesDetails []SeriesDetail `json:"series_details"`
	} `json:"result"`
}

type SeriesDetail struct {
	SeriesCode         string `json:"series_code"`
	ObservationStart   string `json:"observation_start"`
	ObservationEnd     string `json:"observation_end"`
	SeriesNameTh       string `json:"series_name_th"`
	SeriesNameEng      string `json:"series_name_eng"`
	SeriesCategories   string `json:"series_categories"`
	Frequency          string `json:"frequency"`
	FrequencyShort     string `json:"frequency_short"`
	UnitTh             string `json:"unit_th"`
	DataType           string `json:"data_type"`
	SeasonalAdjustment string `json:"seasonal_adjustment_flag"`
	LastUpdatedDate    string `json:"last_updated_date"`
	SourceOfDataTh     string `json:"source_of_data_th"`
	SourceOfDataEng    string `json:"source_of_data_eng"`
	LagTime            string `json:"lag_time"`
	ReleaseScheduleTh  string `json:"release_schedule_th"`
	ReleaseScheduleEng string `json:"release_schedule_eng"`
	AnnotationTh       string `json:"annotation_th"`
	AnnotationEng      string `json:"annotation_eng"`
	DescriptionTh      string `json:"description_th"`
	DescriptionEng     string `json:"description_eng"`
}

func (c *Client) GetCategoryList(ctx context.Context) (*CategoryListResponse, error) {
	var result CategoryListResponse
	if err := c.requestJSON(ctx, categoryListBaseURL, "/category_list/", nil, &result); err != nil {
		return nil, fmt.Errorf("failed to get category list: %w", err)
	}
	return &result, nil
}

func (c *Client) GetSeriesList(ctx context.Context, category string) (*SeriesListResponse, error) {
	query := url.Values{}
	query.Set("category", category)

	var result SeriesListResponse
	if err := c.requestJSON(ctx, categoryListBaseURL, "/series_list/", query, &result); err != nil {
		return nil, fmt.Errorf("failed to get series list: %w", err)
	}
	return &result, nil
}

func (c *Client) GetObservations(ctx context.Context, seriesCode, startPeriod, endPeriod, sortBy string) (*ObservationsResponse, error) {
	query := url.Values{}
	query.Set("series_code", seriesCode)
	query.Set("start_period", startPeriod)
	if endPeriod != "" {
		query.Set("end_period", endPeriod)
	}
	if sortBy != "" {
		query.Set("sort_by", sortBy)
	}

	var result ObservationsResponse
	if err := c.requestJSON(ctx, observationsBaseURL, "/", query, &result); err != nil {
		return nil, fmt.Errorf("failed to get observations: %w", err)
	}
	return &result, nil
}

func (c *Client) SearchSeries(ctx context.Context, keyword string) (*SearchResponse, error) {
	query := url.Values{}
	query.Set("keyword", keyword)

	var result SearchResponse
	if err := c.requestJSON(ctx, searchBaseURL, "/", query, &result); err != nil {
		return nil, fmt.Errorf("failed to search series: %w", err)
	}
	return &result, nil
}
