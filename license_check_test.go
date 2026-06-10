package bankofthailand

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestLicenseCheckResponseDecoding(t *testing.T) {
	jsonData := `{
		"ResultSet": [
			{"name": "Test Corp", "licenseType": "P-Loan"},
			{"name": "Another Co", "licenseType": "e-Money"}
		],
		"ResultSetInfo": {
			"QueryRecordPerPage": 10,
			"QueryTotalRecord": 2,
			"QueryCurrentPage": 1
		},
		"GroupInfo": [
			{"TypeCode": "", "TypeNameTH": "ทั้งหมด", "Count": 2},
			{"TypeCode": "j", "TypeNameTH": "นิติบุคคล", "Count": 2},
			{"TypeCode": "i", "TypeNameTH": "บุคคล", "Count": 0},
			{"TypeCode": "b", "TypeNameTH": "สถานประกอบการ", "Count": 0}
		]
	}`

	var resp LicenseCheckResponse
	if err := json.NewDecoder(strings.NewReader(jsonData)).Decode(&resp); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}

	if len(resp.ResultSet) != 2 {
		t.Fatalf("expected 2 results, got %d", len(resp.ResultSet))
	}
	if resp.ResultSetInfo.QueryTotalRecord != 2 {
		t.Errorf("expected total record 2, got %d", resp.ResultSetInfo.QueryTotalRecord)
	}
	if resp.ResultSetInfo.QueryRecordPerPage != 10 {
		t.Errorf("expected per page 10, got %d", resp.ResultSetInfo.QueryRecordPerPage)
	}
	if resp.ResultSetInfo.QueryCurrentPage != 1 {
		t.Errorf("expected current page 1, got %d", resp.ResultSetInfo.QueryCurrentPage)
	}
	if len(resp.GroupInfo) != 4 {
		t.Fatalf("expected 4 groups, got %d", len(resp.GroupInfo))
	}
	if resp.GroupInfo[1].TypeCode != "j" {
		t.Errorf("expected group type j, got %s", resp.GroupInfo[1].TypeCode)
	}
	if resp.GroupInfo[1].TypeNameTH != "นิติบุคคล" {
		t.Errorf("expected group name นิติบุคคล, got %s", resp.GroupInfo[1].TypeNameTH)
	}
	if resp.GroupInfo[1].Count != 2 {
		t.Errorf("expected group count 2, got %d", resp.GroupInfo[1].Count)
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

	if resp.ResultSet != nil {
		t.Errorf("expected nil ResultSet, got %v", resp.ResultSet)
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
