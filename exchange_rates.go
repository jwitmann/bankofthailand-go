package bankofthailand

import (
	"context"
	"net/url"
)

const (
	exchangeRateBaseURL  = "https://gateway.api.bot.or.th/Stat-ExchangeRate/v2"
	referenceRateBaseURL = "https://gateway.api.bot.or.th/Stat-ReferenceRate/v2"
	spotRateBaseURL      = "https://gateway.api.bot.or.th/Stat-SpotRate/v2/SPOTRATE"
	swapPointBaseURL     = "https://gateway.api.bot.or.th/Stat-SwapPoint/v2/SWAPPOINT"
	impliedRateBaseURL   = "https://gateway.api.bot.or.th/Stat-ThaiBahtImpliedInterestRate/v2/THB_IMPL_INT_RATE"
)

type ExchangeRateResponse struct {
	Result struct {
		API       string `json:"api"`
		Timestamp string `json:"timestamp"`
		Data      struct {
			DataHeader DataHeader         `json:"data_header"`
			DataDetail []ExchangeRateData `json:"data_detail"`
		} `json:"data"`
	} `json:"result"`
}

type DataHeader struct {
	ReportNameEng    string         `json:"report_name_eng"`
	ReportNameTh     string         `json:"report_name_th"`
	ReportUOQNameEng string         `json:"report_uoq_name_eng"`
	ReportUOQNameTh  string         `json:"report_uoq_name_th"`
	SourceOfData     []SourceOfData `json:"report_source_of_data"`
	Remarks          []Remark       `json:"report_remark"`
	LastUpdated      string         `json:"last_updated"`
}

type SourceOfData struct {
	SourceEng string `json:"source_of_data_eng"`
	SourceTh  string `json:"source_of_data_th"`
}

type Remark struct {
	RemarkEng string `json:"report_remark_eng"`
	RemarkTh  string `json:"report_remark_th"`
}

type ExchangeRateData struct {
	Period          string `json:"period"`
	CurrencyID      string `json:"currency_id"`
	CurrencyNameTh  string `json:"currency_name_th"`
	CurrencyNameEng string `json:"currency_name_eng"`
	BuyingSight     string `json:"buying_sight"`
	BuyingTransfer  string `json:"buying_transfer"`
	Selling         string `json:"selling"`
	MidRate         string `json:"mid_rate"`
}

type ReferenceRateResponse struct {
	Result struct {
		API       string `json:"api"`
		Timestamp string `json:"timestamp"`
		Data      struct {
			DataHeader DataHeader          `json:"data_header"`
			DataDetail []ReferenceRateData `json:"data_detail"`
		} `json:"data"`
	} `json:"result"`
}

type ReferenceRateData struct {
	Period string `json:"period"`
	Rate   string `json:"rate"`
}

type SpotRateResponse struct {
	Result struct {
		API       string `json:"api"`
		Timestamp string `json:"timestamp"`
		Data      struct {
			DataHeader DataHeader     `json:"data_header"`
			DataDetail []SpotRateData `json:"data_detail"`
		} `json:"data"`
	} `json:"result"`
}

type SpotRateData struct {
	Period    string `json:"period"`
	BidRate   string `json:"bid_rate"`
	OfferRate string `json:"offer_rate"`
}

type SwapPointResponse struct {
	Result struct {
		API       string `json:"api"`
		Timestamp string `json:"timestamp"`
		Data      struct {
			DataHeader DataHeader      `json:"data_header"`
			DataDetail []SwapPointData `json:"data_detail"`
		} `json:"data"`
	} `json:"result"`
}

type SwapPointData struct {
	Period          string `json:"period"`
	TermTypeNameTh  string `json:"term_type_name_th"`
	TermTypeNameEng string `json:"term_type_name_eng"`
	BidRate         string `json:"bid_rate"`
	OfferRate       string `json:"offer_rate"`
}

type ImpliedRateResponse struct {
	Result struct {
		API       string `json:"api"`
		Timestamp string `json:"timestamp"`
		Data      struct {
			DataHeader DataHeader        `json:"data_header"`
			DataDetail []ImpliedRateData `json:"data_detail"`
		} `json:"data"`
	} `json:"result"`
}

type ImpliedRateData struct {
	Period          string `json:"period"`
	RateTypeNameTh  string `json:"rate_type_name_th"`
	RateTypeNameEng string `json:"rate_type_name_eng"`
	InterestRate    string `json:"interest_rate"`
}

func (c *Client) GetDailyAverageExchangeRate(ctx context.Context, startPeriod, endPeriod, currency string) (*ExchangeRateResponse, error) {
	query := url.Values{}
	query.Set("start_period", startPeriod)
	query.Set("end_period", endPeriod)
	setQuery(query, "currency", currency)

	return getEndpoint[ExchangeRateResponse](ctx, c, exchangeRateBaseURL, "/DAILY_AVG_EXG_RATE/", query, "failed to get exchange rate")
}

func (c *Client) GetDailyReferenceRate(ctx context.Context, startPeriod, endPeriod string) (*ReferenceRateResponse, error) {
	query := url.Values{}
	query.Set("start_period", startPeriod)
	query.Set("end_period", endPeriod)

	return getEndpoint[ReferenceRateResponse](ctx, c, referenceRateBaseURL, "/DAILY_REF_RATE/", query, "failed to get reference rate")
}

func (c *Client) GetSpotRate(ctx context.Context, startPeriod, endPeriod string) (*SpotRateResponse, error) {
	query := url.Values{}
	query.Set("start_period", startPeriod)
	query.Set("end_period", endPeriod)

	return getEndpoint[SpotRateResponse](ctx, c, spotRateBaseURL, "/", query, "failed to get spot rate")
}

func (c *Client) GetSwapPoint(ctx context.Context, startPeriod, endPeriod, termType string) (*SwapPointResponse, error) {
	query := url.Values{}
	query.Set("start_period", startPeriod)
	query.Set("end_period", endPeriod)
	setQuery(query, "term_type", termType)

	return getEndpoint[SwapPointResponse](ctx, c, swapPointBaseURL, "/", query, "failed to get swap point")
}

func (c *Client) GetImpliedInterestRate(ctx context.Context, startPeriod, endPeriod, rateType string) (*ImpliedRateResponse, error) {
	query := url.Values{}
	query.Set("start_period", startPeriod)
	query.Set("end_period", endPeriod)
	setQuery(query, "rate_type", rateType)

	return getEndpoint[ImpliedRateResponse](ctx, c, impliedRateBaseURL, "/", query, "failed to get implied rate")
}
