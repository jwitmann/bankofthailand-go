package bankofthailand

import (
	"context"
	"fmt"
	"net/url"
)

const debtSecuritiesBaseURL = "https://gateway.api.bot.or.th/BondAuction/bond_auction_v2/"

type DebtSecuritiesResponse struct {
	Result struct {
		API       string             `json:"api"`
		Timestamp string             `json:"timestamp"`
		Data      DebtSecuritiesData `json:"data"`
	} `json:"result"`
}

type DebtSecuritiesData struct {
	DataDetail []DebtSecuritiesRecord `json:"data_detail"`
}

type DebtSecuritiesRecord struct {
	AuctionDate                      string `json:"auction_date"`
	DebtSecuritiesType               string `json:"debt_securities_type"`
	ThaiBMASymbol                    string `json:"thaibma_symbol"`
	ISINCode                         string `json:"isin_code"`
	AuctionNameTh                    string `json:"auction_nm_th"`
	CFICode                          string `json:"cfi_code"`
	ReOpenFromTh                     string `json:"re_open_from_th"`
	CouponRate                       string `json:"coupon_rate"`
	TimeToMaturity                   string `json:"time_to_maturity"`
	PaymentDate                      string `json:"payment_date"`
	StartDateOfInterestEarningPeriod string `json:"start_date_of_interest_earning_period"`
	MaturityDate                     string `json:"maturity_date"`
	IssueAmountNCB_CB                string `json:"issue_amount_ncb_cb"`
	AcceptedAmountNCB_CB             string `json:"accepted_amount_ncb_cb"`
	AcceptedAmountNCB                string `json:"accepted_amount_ncb"`
	AcceptedAmountCB                 string `json:"accepted_amount_cb"`
	GreenshoeOptionAmount            string `json:"greenshoe_option_amount"`
	PAOAmount                        string `json:"pao_amount"`
	OverAllotmentAmount              string `json:"over_allotment_amount"`
	GrandTotalAmount                 string `json:"grand_total_amount"`
	AcceptedLowestYield              string `json:"accepted_lowest_yield"`
	AcceptedHighestYield             string `json:"accepted_highest_yield"`
	WeightedAverageAcceptedYield     string `json:"weighted_average_accepted_yield"`
	BidCoverageRatio                 string `json:"bid_coverage_ratio"`
	AuctionStatus                    string `json:"auction_st"`
}

func (r DebtSecuritiesRecord) AuctionName(loc Locale) string {
	return pickString(loc, r.AuctionNameTh, translateAuctionName(r.AuctionNameTh))
}

func (r DebtSecuritiesRecord) ReOpenFrom(loc Locale) string {
	return pickString(loc, r.ReOpenFromTh, translateAuctionName(r.ReOpenFromTh))
}

func (c *Client) GetDebtSecuritiesAuction(ctx context.Context, startPeriod, endPeriod string) (*DebtSecuritiesResponse, error) {
	query := url.Values{}
	query.Set("start_period", startPeriod)
	query.Set("end_period", endPeriod)

	var result DebtSecuritiesResponse
	if err := c.requestJSON(ctx, debtSecuritiesBaseURL, "/", query, &result); err != nil {
		return nil, fmt.Errorf("failed to get debt securities auction: %w", err)
	}
	return &result, nil
}
