package bankofthailand

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestDebtSecuritiesResponseDecoding(t *testing.T) {
	jsonData := `{
		"result": {
			"api": "Bond Auction",
			"timestamp": "2017-10-03 08:53:51",
			"data": {
				"data_detail": [
					{
						"auction_date": "2017-09-26",
						"debt_securities_type": "Government Bonds",
						"thaibma_symbol": "LB233A",
						"isin_code": "TH0623033303",
						"auction_nm_th": "พันธบัตรรัฐบาลในปีงบประมาณ พ.ศ.2560 ครั้งที่ 1-Test API",
						"cfi_code": "DBFTFR",
						"re_open_from_th": "พันธบัตรรัฐบาลในปีงบประมาณ พ.ศ.2551 ครั้งที่ 4",
						"coupon_rate": "5.5",
						"time_to_maturity": "5.46 Yrs",
						"payment_date": "2017-09-28",
						"start_date_of_interest_earning_period": "2017-09-13",
						"maturity_date": "2023-03-13",
						"issue_amount_ncb_cb": "2000.0000000",
						"accepted_amount_ncb_cb": "2000.0000000",
						"accepted_amount_ncb": "",
						"accepted_amount_cb": "2000.0000000",
						"greenshoe_option_amount": "400.0000000",
						"pao_amount": "",
						"over_allotment_amount": "",
						"grand_total_amount": "2400.0000000",
						"accepted_lowest_yield": "1.7070000",
						"accepted_highest_yield": "1.7090000",
						"weighted_average_accepted_yield": "1.7077000",
						"bid_coverage_ratio": "2.2000000",
						"auction_st": "Approve"
					}
				]
			}
		}
	}`

	var resp DebtSecuritiesResponse
	if err := json.NewDecoder(strings.NewReader(jsonData)).Decode(&resp); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}

	if resp.Result.API != "Bond Auction" {
		t.Errorf("expected API 'Bond Auction', got %s", resp.Result.API)
	}
	if len(resp.Result.Data.DataDetail) != 1 {
		t.Fatalf("expected 1 record, got %d", len(resp.Result.Data.DataDetail))
	}
	rec := resp.Result.Data.DataDetail[0]
	if rec.AuctionDate != "2017-09-26" {
		t.Errorf("expected auction_date 2017-09-26, got %s", rec.AuctionDate)
	}
	if rec.DebtSecuritiesType != "Government Bonds" {
		t.Errorf("expected type Government Bonds, got %s", rec.DebtSecuritiesType)
	}
	if rec.ThaiBMASymbol != "LB233A" {
		t.Errorf("expected thaibma_symbol LB233A, got %s", rec.ThaiBMASymbol)
	}
	if rec.ISINCode != "TH0623033303" {
		t.Errorf("expected isin_code TH0623033303, got %s", rec.ISINCode)
	}
	if rec.CouponRate != "5.5" {
		t.Errorf("expected coupon_rate 5.5, got %s", rec.CouponRate)
	}
	if rec.AcceptedLowestYield != "1.7070000" {
		t.Errorf("expected accepted_lowest_yield 1.7070000, got %s", rec.AcceptedLowestYield)
	}
	if rec.WeightedAverageAcceptedYield != "1.7077000" {
		t.Errorf("expected weighted_average_accepted_yield 1.7077000, got %s", rec.WeightedAverageAcceptedYield)
	}
	if rec.BidCoverageRatio != "2.2000000" {
		t.Errorf("expected bid_coverage_ratio 2.2000000, got %s", rec.BidCoverageRatio)
	}
	if rec.AuctionStatus != "Approve" {
		t.Errorf("expected auction_st Approve, got %s", rec.AuctionStatus)
	}
}

func TestDebtSecuritiesResponseDecoding_EmptyResultSet(t *testing.T) {
	jsonData := `{
		"result": {
			"api": "Bond Auction",
			"timestamp": "2017-10-03 08:53:51",
			"data": {
				"data_detail": []
			}
		}
	}`

	var resp DebtSecuritiesResponse
	if err := json.NewDecoder(strings.NewReader(jsonData)).Decode(&resp); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}

	if len(resp.Result.Data.DataDetail) != 0 {
		t.Errorf("expected 0 records, got %d", len(resp.Result.Data.DataDetail))
	}
}
