package bankofthailand

import "testing"

func TestConvertBuddhistToGregorian(t *testing.T) {
	got, err := convertBuddhistToGregorian("2560")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != 2017 {
		t.Errorf("convertBuddhistToGregorian(2560) = %d, want 2017", got)
	}

	got2, err := convertBuddhistToGregorian("2551")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got2 != 2008 {
		t.Errorf("convertBuddhistToGregorian(2551) = %d, want 2008", got2)
	}
}

func TestExtractBondPrefix(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"พันธบัตรรัฐบาลในปีงบประมาณ พ.ศ.2560 ครั้งที่ 1", "พันธบัตรรัฐบาล"},
		{"บัตรเงินคลังในปีงบประมาณ พ.ศ.2561 ครั้งที่ 2", "บัตรเงินคลัง"},
		{"no pattern here", ""},
	}

	for _, tt := range tests {
		got := extractBondPrefix(tt.input)
		if got != tt.expected {
			t.Errorf("extractBondPrefix(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

func TestExtractFiscalYear(t *testing.T) {
	got := extractFiscalYear("พันธบัตรรัฐบาลในปีงบประมาณ พ.ศ.2560 ครั้งที่ 1")
	if got != "2560" {
		t.Errorf("extractFiscalYear = %q, want 2560", got)
	}

	got2 := extractFiscalYear("no year")
	if got2 != "" {
		t.Errorf("extractFiscalYear = %q, want empty", got2)
	}
}

func TestExtractIssueNumber(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"ครั้งที่ 1", "1"},
		{"ครั้งที่ 12", "12"},
		{"ครั้งที่ 3-Test", "3"},
		{"no issue", ""},
	}

	for _, tt := range tests {
		got := extractIssueNumber(tt.input)
		if got != tt.expected {
			t.Errorf("extractIssueNumber(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

func TestTranslateAuctionName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"พันธบัตรรัฐบาลในปีงบประมาณ พ.ศ.2560 ครั้งที่ 1",
			"Government Bond in Fiscal Year 2017 (B.E. 2560), Issue 1",
		},
		{
			"พันธบัตรรัฐบาลในปีงบประมาณ พ.ศ.2551 ครั้งที่ 4",
			"Government Bond in Fiscal Year 2008 (B.E. 2551), Issue 4",
		},
		{"", ""},
		{"unknown prefix ในปีงบประมาณ พ.ศ.2560 ครั้งที่ 1", "unknown prefix ในปีงบประมาณ พ.ศ.2560 ครั้งที่ 1"},
	}

	for _, tt := range tests {
		got := translateAuctionName(tt.input)
		if got != tt.expected {
			t.Errorf("translateAuctionName(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

func TestDebtSecuritiesRecordAuctionName(t *testing.T) {
	rec := DebtSecuritiesRecord{
		AuctionNameTh: "พันธบัตรรัฐบาลในปีงบประมาณ พ.ศ.2560 ครั้งที่ 1",
	}

	if got := rec.AuctionName(LocaleThai); got != rec.AuctionNameTh {
		t.Errorf("AuctionName(Thai) = %q, want %q", got, rec.AuctionNameTh)
	}

	expectedEn := "Government Bond in Fiscal Year 2017 (B.E. 2560), Issue 1"
	if got := rec.AuctionName(LocaleEnglish); got != expectedEn {
		t.Errorf("AuctionName(English) = %q, want %q", got, expectedEn)
	}
}

func TestDebtSecuritiesRecordReOpenFrom(t *testing.T) {
	rec := DebtSecuritiesRecord{
		ReOpenFromTh: "พันธบัตรรัฐบาลในปีงบประมาณ พ.ศ.2551 ครั้งที่ 4",
	}

	if got := rec.ReOpenFrom(LocaleThai); got != rec.ReOpenFromTh {
		t.Errorf("ReOpenFrom(Thai) = %q, want %q", got, rec.ReOpenFromTh)
	}

	expectedEn := "Government Bond in Fiscal Year 2008 (B.E. 2551), Issue 4"
	if got := rec.ReOpenFrom(LocaleEnglish); got != expectedEn {
		t.Errorf("ReOpenFrom(English) = %q, want %q", got, expectedEn)
	}
}
