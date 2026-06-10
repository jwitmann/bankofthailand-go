package bankofthailand

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestLicenseCheckResponseDecoding(t *testing.T) {
	jsonData := `{
		"ResultSet": [
			{"Id": "1", "AuthorizedName": " ธนาคารออมสิน ", "BranchName": " ธนาคารออมสิน ", "TypeId": "ธนาคารออมสิน", "TypeName": "ธนาคารออมสิน", "LastUpdate": "", "Address": "470 ถนนพหลโยธิน", "Telephone": "0-2299-8000", "DepositFlag": "T", "LoanFlag": "T"}
		],
		"ResultSetInfo": {
			"QueryRecordPerPage": 10,
			"QueryTotalRecord": 1,
			"QueryCurrentPage": 1
		},
		"GroupInfo": [
			{"TypeCode": "", "TypeNameTH": "ทั้งหมด", "Count": 1},
			{"TypeCode": "j", "TypeNameTH": "นิติบุคคล", "Count": 1}
		]
	}`

	var resp LicenseCheckResponse
	if err := json.NewDecoder(strings.NewReader(jsonData)).Decode(&resp); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}

	if len(resp.ResultSet) != 1 {
		t.Fatalf("expected 1 result, got %d", len(resp.ResultSet))
	}
	rec := resp.ResultSet[0]
	if rec.ID != "1" {
		t.Errorf("expected ID 1, got %s", rec.ID)
	}
	if rec.AuthorizedName != " ธนาคารออมสิน " {
		t.Errorf("unexpected AuthorizedName: %s", rec.AuthorizedName)
	}
	if rec.TypeName != "ธนาคารออมสิน" {
		t.Errorf("unexpected TypeName: %s", rec.TypeName)
	}
	if rec.DepositFlag != "T" {
		t.Errorf("expected DepositFlag T, got %s", rec.DepositFlag)
	}
	if resp.ResultSetInfo.QueryTotalRecord != 1 {
		t.Errorf("expected total record 1, got %d", resp.ResultSetInfo.QueryTotalRecord)
	}
	if resp.ResultSetInfo.QueryRecordPerPage != 10 {
		t.Errorf("expected per page 10, got %d", resp.ResultSetInfo.QueryRecordPerPage)
	}
	if len(resp.GroupInfo) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(resp.GroupInfo))
	}
	if resp.GroupInfo[1].TypeCode != "j" {
		t.Errorf("expected group type j, got %s", resp.GroupInfo[1].TypeCode)
	}
	if resp.GroupInfo[1].TypeNameTH != "นิติบุคคล" {
		t.Errorf("expected group name นิติบุคคล, got %s", resp.GroupInfo[1].TypeNameTH)
	}
	if resp.GroupInfo[1].Count != 1 {
		t.Errorf("expected group count 1, got %d", resp.GroupInfo[1].Count)
	}

	if rec.TypeNameLocalized(LocaleEnglish) != "Government Savings Bank" {
		t.Errorf("unexpected English TypeName: %s", rec.TypeNameLocalized(LocaleEnglish))
	}
	if rec.TypeNameLocalized(LocaleThai) != "ธนาคารออมสิน" {
		t.Errorf("unexpected Thai TypeName: %s", rec.TypeNameLocalized(LocaleThai))
	}
}

func TestLicenseCheckResponseDecoding_NoResults(t *testing.T) {
	jsonData := `{
		"ResultSet": null,
		"ResultSetInfo": {
			"QueryRecordPerPage": 1,
			"QueryTotalRecord": 0,
			"QueryCurrentPage": 1
		},
		"GroupInfo": [
			{"TypeCode": "", "TypeNameTH": "ทั้งหมด", "Count": 0},
			{"TypeCode": "j", "TypeNameTH": "นิติบุคคล", "Count": 0}
		]
	}`

	var resp LicenseCheckResponse
	if err := json.NewDecoder(strings.NewReader(jsonData)).Decode(&resp); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}

	if len(resp.ResultSet) > 0 {
		t.Errorf("expected empty ResultSet, got %v", resp.ResultSet)
	}
	if resp.ResultSetInfo.QueryTotalRecord != 0 {
		t.Errorf("expected 0 total records, got %d", resp.ResultSetInfo.QueryTotalRecord)
	}
}

func TestAuthorizedDetailResponseDecoding(t *testing.T) {
	jsonData := `{
		"AuthorizationInfo": {
			"Id": "123",
			"AuthorizedName": "บริษัท เครดีโว่ (ไทยแลนด์) จำกัด",
			"BranchName": "",
			"TypeId": "ผู้ประกอบธุรกิจการให้สินเชื่อที่มิใช่สถาบันการเงิน",
			"TypeName": "ผู้ประกอบธุรกิจการให้สินเชื่อที่มิใช่สถาบันการเงิน",
			"LastUpdate": "10/06/2026"
		}
	}`

	var resp AuthorizedDetailResponse
	if err := json.NewDecoder(strings.NewReader(jsonData)).Decode(&resp); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}

	if resp.AuthorizationInfo.ID != "123" {
		t.Errorf("expected ID 123, got %s", resp.AuthorizationInfo.ID)
	}
	if resp.AuthorizationInfo.AuthorizedName != "บริษัท เครดีโว่ (ไทยแลนด์) จำกัด" {
		t.Errorf("unexpected AuthorizedName: %s", resp.AuthorizationInfo.AuthorizedName)
	}
	if resp.AuthorizationInfo.TypeID != "ผู้ประกอบธุรกิจการให้สินเชื่อที่มิใช่สถาบันการเงิน" {
		t.Errorf("unexpected TypeId: %s", resp.AuthorizationInfo.TypeID)
	}
	if resp.AuthorizationInfo.LastUpdate != "10/06/2026" {
		t.Errorf("expected LastUpdate 10/06/2026, got %s", resp.AuthorizationInfo.LastUpdate)
	}
}
