package bankofthailand

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestExchangeRateResponseDecoding(t *testing.T) {
	jsonData := `{
		"result": {
			"api": "DailyAverageExchangeRate",
			"timestamp": "2026-06-06 12:00:00",
			"data": {
				"data_header": {
					"report_name_eng": "Average Exchange Rate",
					"report_name_th": "อัตราแลกเปลี่ยนเฉลี่ย",
					"report_uoq_name_eng": "Baht per 1 unit of foreign currency",
					"report_uoq_name_th": "บาทต่อ 1 หน่วยเงินตราต่างประเทศ",
					"report_source_of_data": [
						{"source_of_data_eng": "Bank of Thailand", "source_of_data_th": "ธนาคารแห่งประเทศไทย"}
					],
					"report_remark": [],
					"last_updated": "2026-06-05 18:00:00"
				},
				"data_detail": [
					{
						"period": "2026-06-05",
						"currency_id": "USD",
						"currency_name_th": "ดอลลาร์สหรัฐ",
						"currency_name_eng": "US DOLLAR",
						"buying_sight": "33.50",
						"buying_transfer": "33.55",
						"selling": "34.00",
						"mid_rate": "33.775"
					}
				]
			}
		}
	}`

	var resp ExchangeRateResponse
	if err := json.NewDecoder(strings.NewReader(jsonData)).Decode(&resp); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}

	if resp.Result.API != "DailyAverageExchangeRate" {
		t.Errorf("expected API DailyAverageExchangeRate, got %s", resp.Result.API)
	}
	if len(resp.Result.Data.DataDetail) != 1 {
		t.Fatalf("expected 1 rate, got %d", len(resp.Result.Data.DataDetail))
	}
	rate := resp.Result.Data.DataDetail[0]
	if rate.CurrencyID != "USD" {
		t.Errorf("expected currency USD, got %s", rate.CurrencyID)
	}
	if rate.MidRate != "33.775" {
		t.Errorf("expected mid_rate 33.775, got %s", rate.MidRate)
	}
	if resp.Result.Data.DataHeader.ReportNameEng != "Average Exchange Rate" {
		t.Errorf("expected report name 'Average Exchange Rate', got %s", resp.Result.Data.DataHeader.ReportNameEng)
	}
}

func TestReferenceRateResponseDecoding(t *testing.T) {
	jsonData := `{
		"result": {
			"api": "DailyReferenceRate",
			"timestamp": "2026-06-06 12:00:00",
			"data": {
				"data_header": {
					"report_name_eng": "Weighted-average Interbank Exchange Rate",
					"report_name_th": "อัตราแลกเปลี่ยนเฉลี่ยถ่วงน้ำหนักระหว่างธนาคาร",
					"report_uoq_name_eng": "Baht per 1 USD",
					"report_uoq_name_th": "บาทต่อ 1 ดอลลาร์สหรัฐ",
					"report_source_of_data": [],
					"report_remark": [],
					"last_updated": "2026-06-05 18:00:00"
				},
				"data_detail": [
					{"period": "2026-06-05", "rate": "33.77"}
				]
			}
		}
	}`

	var resp ReferenceRateResponse
	if err := json.NewDecoder(strings.NewReader(jsonData)).Decode(&resp); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}

	if len(resp.Result.Data.DataDetail) != 1 {
		t.Fatalf("expected 1 rate, got %d", len(resp.Result.Data.DataDetail))
	}
	if resp.Result.Data.DataDetail[0].Rate != "33.77" {
		t.Errorf("expected rate 33.77, got %s", resp.Result.Data.DataDetail[0].Rate)
	}
}

func TestPolicyRateResponseDecoding(t *testing.T) {
	jsonData := `{
		"result": {
			"api": "PolicyRate",
			"timestamp": "2026-06-06 12:00:00",
			"data": "2.50",
			"announcement_date": "2026-05-28",
			"news_text_en": "The Monetary Policy Committee decided...",
			"news_text_th": "คณะกรรมการนโยบายการเงินมีมติ...",
			"effective_datetime": "2026-05-28 14:00:00"
		}
	}`

	var resp PolicyRateResponse
	if err := json.NewDecoder(strings.NewReader(jsonData)).Decode(&resp); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}

	if resp.Result.Data != "2.50" {
		t.Errorf("expected rate 2.50, got %s", resp.Result.Data)
	}
	if resp.Result.AnnouncementDate != "2026-05-28" {
		t.Errorf("expected announcement date 2026-05-28, got %s", resp.Result.AnnouncementDate)
	}
}

func TestCategoryListResponseDecoding(t *testing.T) {
	jsonData := `{
		"result": {
			"api": "CategoryList",
			"timestamp": "2026-06-06 12:00:00",
			"category": [
				{"category": "EC_XT_077", "description_th": "อัตราแลกเปลี่ยน", "description_eng": "Exchange Rates"},
				{"category": "PF_SF_000", "description_th": "หนี้สาธารณะ", "description_eng": "Public Debt"}
			]
		}
	}`

	var resp CategoryListResponse
	if err := json.NewDecoder(strings.NewReader(jsonData)).Decode(&resp); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}

	if len(resp.Result.Category) != 2 {
		t.Fatalf("expected 2 categories, got %d", len(resp.Result.Category))
	}
	if resp.Result.Category[0].Category != "EC_XT_077" {
		t.Errorf("expected category EC_XT_077, got %s", resp.Result.Category[0].Category)
	}
}

func TestObservationsResponseDecoding(t *testing.T) {
	jsonData := `{
		"result": {
			"api": "Observations",
			"timestamp": "2026-06-06 12:00:00",
			"series": [
				{
					"series_code": "PF00000000Q00232",
					"series_name_th": "ตราสารหนี้ภาครัฐ_รวม",
					"series_name_eng": "Government debt securities_Total",
					"unit_th": "ล้านบาท",
					"unit_eng": "Million Baht",
					"series_type": "Stock",
					"frequency": "Quarterly",
					"last_update_date": "2026-05-15",
					"observations": [
						{"period_start": "2017-Q1", "value": "8648519.0000000"},
						{"period_start": "2017-Q2", "value": "8712345.0000000"}
					]
				}
			]
		}
	}`

	var resp ObservationsResponse
	if err := json.NewDecoder(strings.NewReader(jsonData)).Decode(&resp); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}

	if len(resp.Result.Series) != 1 {
		t.Fatalf("expected 1 series, got %d", len(resp.Result.Series))
	}
	series := resp.Result.Series[0]
	if series.SeriesCode != "PF00000000Q00232" {
		t.Errorf("expected series code PF00000000Q00232, got %s", series.SeriesCode)
	}
	if len(series.Observations) != 2 {
		t.Fatalf("expected 2 observations, got %d", len(series.Observations))
	}
	if series.Observations[0].Value != "8648519.0000000" {
		t.Errorf("expected value 8648519.0000000, got %s", series.Observations[0].Value)
	}
}

func TestBIBORResponseDecoding(t *testing.T) {
	jsonData := `{
		"result": {
			"api": "BIBOR",
			"timestamp": "2026-06-06 12:00:00",
			"data": {
				"data_detail": [
					{
						"period": "2026-06-05",
						"bankname_th": "ธนาคารกสิกรไทย",
						"bankname_eng": "KASIKORNBANK",
						"bibor_o_n": "2.45",
						"bibor_1_week": "2.50",
						"bibor_1_month": "2.60",
						"bibor_3_month": "2.75",
						"bibor_6_month": "2.90",
						"bibor_1_year": "3.10"
					}
				]
			}
		}
	}`

	var resp BIBORResponse
	if err := json.NewDecoder(strings.NewReader(jsonData)).Decode(&resp); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}

	if len(resp.Result.Data.DataDetail) != 1 {
		t.Fatalf("expected 1 record, got %d", len(resp.Result.Data.DataDetail))
	}
	bibor := resp.Result.Data.DataDetail[0]
	if bibor.BIBOR3M != "2.75" {
		t.Errorf("expected BIBOR 3M 2.75, got %s", bibor.BIBOR3M)
	}
}

func TestSpotRateResponseDecoding(t *testing.T) {
	jsonData := `{
		"result": {
			"api": "SpotRate",
			"timestamp": "2026-06-06 12:00:00",
			"data": {
				"data_header": {
					"report_name_eng": "Spot Rate",
					"report_name_th": "อัตราแลกเปลี่ยนสปอต",
					"report_uoq_name_eng": "Baht per 1 USD",
					"report_uoq_name_th": "บาทต่อ 1 ดอลลาร์สหรัฐ",
					"report_source_of_data": [],
					"report_remark": [],
					"last_updated": "2026-06-05 18:00:00"
				},
				"data_detail": [
					{"period": "2026-06-05", "bid_rate": "33.75", "offer_rate": "33.79"}
				]
			}
		}
	}`

	var resp SpotRateResponse
	if err := json.NewDecoder(strings.NewReader(jsonData)).Decode(&resp); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}

	if len(resp.Result.Data.DataDetail) != 1 {
		t.Fatalf("expected 1 rate, got %d", len(resp.Result.Data.DataDetail))
	}
	if resp.Result.Data.DataDetail[0].BidRate != "33.75" {
		t.Errorf("expected bid 33.75, got %s", resp.Result.Data.DataDetail[0].BidRate)
	}
}
