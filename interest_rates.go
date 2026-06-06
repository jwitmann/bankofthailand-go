package bankofthailand

import (
	"context"
	"encoding/json"
	"fmt"
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
	resp, err := c.GetURL(ctx, policyRateBaseURL+"/")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result PolicyRateResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode policy rate response: %w", err)
	}

	return &result, nil
}

func (c *Client) GetBIBOR(ctx context.Context, startPeriod, endPeriod, bank string) (*BIBORResponse, error) {
	query := url.Values{}
	query.Set("start_period", startPeriod)
	query.Set("end_period", endPeriod)
	if bank != "" {
		query.Set("bank", bank)
	}

	u, _ := url.Parse(biborBaseURL + "/bibor_rate/")
	u.RawQuery = query.Encode()

	resp, err := c.GetURL(ctx, u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result BIBORResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode BIBOR response: %w", err)
	}

	return &result, nil
}

func (c *Client) GetBIBORAverage(ctx context.Context, startPeriod, endPeriod string) (*BIBORResponse, error) {
	query := url.Values{}
	query.Set("start_period", startPeriod)
	query.Set("end_period", endPeriod)

	u, _ := url.Parse(biborBaseURL + "/bibor_avg_rate/")
	u.RawQuery = query.Encode()

	resp, err := c.GetURL(ctx, u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result BIBORResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode BIBOR avg response: %w", err)
	}

	return &result, nil
}

func (c *Client) GetDepositRate(ctx context.Context, startPeriod, endPeriod string) (*DepositRateResponse, error) {
	query := url.Values{}
	query.Set("start_period", startPeriod)
	query.Set("end_period", endPeriod)

	u, _ := url.Parse(depositRateBaseURL + "/deposit_rate/")
	u.RawQuery = query.Encode()

	resp, err := c.GetURL(ctx, u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result DepositRateResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode deposit rate response: %w", err)
	}

	return &result, nil
}

func (c *Client) GetLoanRate(ctx context.Context, startPeriod, endPeriod string) (*LoanRateResponse, error) {
	query := url.Values{}
	query.Set("start_period", startPeriod)
	query.Set("end_period", endPeriod)

	u, _ := url.Parse(loanRateBaseURL + "/loan_rate/")
	u.RawQuery = query.Encode()

	resp, err := c.GetURL(ctx, u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result LoanRateResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode loan rate response: %w", err)
	}

	return &result, nil
}

func (c *Client) GetInterbankTransactionRate(ctx context.Context, startPeriod, endPeriod, termType string) (*InterbankTransactionRateResponse, error) {
	query := url.Values{}
	query.Set("start_period", startPeriod)
	query.Set("end_period", endPeriod)
	if termType != "" {
		query.Set("term_type", termType)
	}

	u, _ := url.Parse(interbankTransactionBaseURL + "/")
	u.RawQuery = query.Encode()

	resp, err := c.GetURL(ctx, u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result InterbankTransactionRateResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode interbank transaction rate response: %w", err)
	}

	return &result, nil
}
