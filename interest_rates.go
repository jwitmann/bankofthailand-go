package bankofthailand

import (
	"context"
	"net/url"
)

const (
	policyRateBaseURL           = "https://gateway.api.bot.or.th/PolicyRate/v3/policy_rate"
	biborBaseURL                = "https://gateway.api.bot.or.th/BIBOR/v2"
	depositRateBaseURL          = "https://gateway.api.bot.or.th/DepositRate/v2"
	loanRateBaseURL             = "https://gateway.api.bot.or.th/LoanRate/v2"
	interbankTransactionBaseURL = "https://gateway.api.bot.or.th/Stat-InterbankTransactionRate/v2/INTRBNK_TXN_RATE"
)

type PolicyRateResponse struct {
	Result struct {
		API               string `json:"api"`
		Timestamp         string `json:"timestamp"`
		Data              string `json:"data"`
		AnnouncementDate  string `json:"announcement_date"`
		NewsTextEn        string `json:"news_text_en"`
		NewsTextTh        string `json:"news_text_th"`
		EffectiveDateTime string `json:"effective_datetime"`
	} `json:"result"`
}

type BIBORResponse struct {
	Result struct {
		API       string `json:"api"`
		Timestamp string `json:"timestamp"`
		Data      struct {
			DataDetail []BIBORData `json:"data_detail"`
		} `json:"data"`
	} `json:"result"`
}

type BIBORData struct {
	Period      string `json:"period"`
	BankNameTh  string `json:"bankname_th,omitempty"`
	BankNameEng string `json:"bankname_eng,omitempty"`
	BIBORON     string `json:"bibor_o_n"`
	BIBOR1W     string `json:"bibor_1_week"`
	BIBOR1M     string `json:"bibor_1_month"`
	BIBOR2M     string `json:"bibor_2_month"`
	BIBOR3M     string `json:"bibor_3_month"`
	BIBOR6M     string `json:"bibor_6_month"`
	BIBOR9M     string `json:"bibor_9_month"`
	BIBOR1Y     string `json:"bibor_1_year"`
}

type DepositRateResponse struct {
	Result struct {
		API       string `json:"api"`
		Timestamp string `json:"timestamp"`
		Data      struct {
			DataDetail []DepositRateData `json:"data_detail"`
		} `json:"data"`
	} `json:"result"`
}

type DepositRateData struct {
	Period          string `json:"period"`
	BankTypeNameTh  string `json:"bank_type_name_th"`
	BankTypeNameEng string `json:"bank_type_name_eng"`
	BankNameTh      string `json:"bank_name_th"`
	BankNameEng     string `json:"bank_name_eng"`
	SavingMin       string `json:"saving_min"`
	SavingMax       string `json:"saving_max"`
	Fix3MthsMin     string `json:"fix_3_mths_min"`
	Fix3MthsMax     string `json:"fix_3_mths_max"`
	Fix6MthsMin     string `json:"fix_6_mths_min"`
	Fix6MthsMax     string `json:"fix_6_mths_max"`
	Fix12MthsMin    string `json:"fix_12_mths_min"`
	Fix12MthsMax    string `json:"fix_12_mths_max"`
	Fix24MthsMin    string `json:"fix_24_mths_min"`
	Fix24MthsMax    string `json:"fix_24_mths_max"`
}

type LoanRateResponse struct {
	Result struct {
		API       string `json:"api"`
		Timestamp string `json:"timestamp"`
		Data      struct {
			DataDetail []LoanRateData `json:"data_detail"`
		} `json:"data"`
	} `json:"result"`
}

type LoanRateData struct {
	Period          string `json:"period"`
	BankTypeNameTh  string `json:"bank_type_name_th"`
	BankTypeNameEng string `json:"bank_type_name_eng"`
	BankNameTh      string `json:"bank_name_th"`
	BankNameEng     string `json:"bank_name_eng"`
	MOR             string `json:"mor"`
	MLR             string `json:"mlr"`
	MRR             string `json:"mrr"`
	CeilingRate     string `json:"ceiling_rate"`
	DefaultRate     string `json:"default_rate"`
	CreditCardMin   string `json:"creditcard_min"`
	CreditCardMax   string `json:"creditcard_max"`
}

type InterbankTransactionRateResponse struct {
	Result struct {
		API       string `json:"api"`
		Timestamp string `json:"timestamp"`
		Data      struct {
			DataHeader DataHeader                     `json:"data_header"`
			DataDetail []InterbankTransactionRateData `json:"data_detail"`
		} `json:"data"`
	} `json:"result"`
}

type InterbankTransactionRateData struct {
	Period                      string `json:"period"`
	TermTypeNameTh              string `json:"term_type_name_th"`
	TermTypeNameEng             string `json:"term_type_name_eng"`
	MinInterestRate             string `json:"min_interest_rate"`
	MaxInterestRate             string `json:"max_interest_rate"`
	ModeInterestRate            string `json:"mode_interest_rate"`
	WeightedAverageInterestRate string `json:"weighted_average_interest_rate"`
}

func (c *Client) GetPolicyRate(ctx context.Context) (*PolicyRateResponse, error) {
	return getEndpoint[PolicyRateResponse](ctx, c, policyRateBaseURL, "/", nil, "failed to get policy rate")
}

func (c *Client) GetBIBOR(ctx context.Context, startPeriod, endPeriod, bank string) (*BIBORResponse, error) {
	query := url.Values{}
	query.Set("start_period", startPeriod)
	query.Set("end_period", endPeriod)
	setQuery(query, "bank", bank)

	return getEndpoint[BIBORResponse](ctx, c, biborBaseURL, "/bibor_rate/", query, "failed to get BIBOR")
}

func (c *Client) GetBIBORAverage(ctx context.Context, startPeriod, endPeriod string) (*BIBORResponse, error) {
	query := url.Values{}
	query.Set("start_period", startPeriod)
	query.Set("end_period", endPeriod)

	return getEndpoint[BIBORResponse](ctx, c, biborBaseURL, "/bibor_avg_rate/", query, "failed to get BIBOR average")
}

func (c *Client) GetDepositRate(ctx context.Context, startPeriod, endPeriod string) (*DepositRateResponse, error) {
	query := url.Values{}
	query.Set("start_period", startPeriod)
	query.Set("end_period", endPeriod)

	return getEndpoint[DepositRateResponse](ctx, c, depositRateBaseURL, "/deposit_rate/", query, "failed to get deposit rate")
}

func (c *Client) GetLoanRate(ctx context.Context, startPeriod, endPeriod string) (*LoanRateResponse, error) {
	query := url.Values{}
	query.Set("start_period", startPeriod)
	query.Set("end_period", endPeriod)

	return getEndpoint[LoanRateResponse](ctx, c, loanRateBaseURL, "/loan_rate/", query, "failed to get loan rate")
}

func (c *Client) GetInterbankTransactionRate(ctx context.Context, startPeriod, endPeriod, termType string) (*InterbankTransactionRateResponse, error) {
	query := url.Values{}
	query.Set("start_period", startPeriod)
	query.Set("end_period", endPeriod)
	setQuery(query, "term_type", termType)

	return getEndpoint[InterbankTransactionRateResponse](ctx, c, interbankTransactionBaseURL, "/", query, "failed to get interbank transaction rate")
}
