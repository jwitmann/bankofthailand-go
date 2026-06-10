package bankofthailand

import (
	"fmt"
	"strconv"
	"strings"
)

func lookup(s string, dict map[string]string) string {
	if v, ok := dict[s]; ok {
		return v
	}
	return s
}

func convertBuddhistToGregorian(yearStr string) (int, error) {
	be, err := strconv.Atoi(yearStr)
	if err != nil {
		return 0, err
	}
	return be - 543, nil
}

func translateAuctionName(name string) string {
	if name == "" {
		return ""
	}

	prefix := extractBondPrefix(name)
	translated := lookup(prefix, AuctionNameTranslation)
	if translated == prefix {
		return name
	}

	fiscalYear := extractFiscalYear(name)
	issueNum := extractIssueNumber(name)

	if fiscalYear != "" && issueNum != "" {
		gregYear, err := convertBuddhistToGregorian(fiscalYear)
		if err == nil {
			return fmt.Sprintf("%s in Fiscal Year %s (B.E. %s), Issue %s", translated, strconv.Itoa(gregYear), fiscalYear, issueNum)
		}
	}

	if fiscalYear != "" {
		gregYear, err := convertBuddhistToGregorian(fiscalYear)
		if err == nil {
			return fmt.Sprintf("%s in Fiscal Year %s (B.E. %s)", translated, strconv.Itoa(gregYear), fiscalYear)
		}
	}

	return name
}

func extractBondPrefix(name string) string {
	idx := strings.Index(name, "ในปีงบประมาณ")
	if idx == -1 {
		idx = strings.Index(name, "ปีงบประมาณ")
	}
	if idx > 0 {
		return strings.TrimSpace(name[:idx])
	}
	return ""
}

func extractFiscalYear(name string) string {
	idx := strings.Index(name, "พ.ศ.")
	if idx == -1 {
		return ""
	}
	start := idx + len("พ.ศ.")
	var end int
	for end = start; end < len(name); end++ {
		if name[end] < '0' || name[end] > '9' {
			break
		}
	}
	if end > start {
		return name[start:end]
	}
	return ""
}

func extractIssueNumber(name string) string {
	idx := strings.Index(name, "ครั้งที่")
	if idx == -1 {
		return ""
	}
	start := idx + len("ครั้งที่")
	name = strings.TrimLeft(name[start:], " ")
	var end int
	for end = 0; end < len(name); end++ {
		if name[end] >= '0' && name[end] <= '9' {
			continue
		}
		if name[end] == '-' || name[end] == ' ' {
			break
		}
	}
	if end > 0 {
		return strings.TrimSpace(name[:end])
	}
	return ""
}
