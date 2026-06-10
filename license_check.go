package bankofthailand

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"strconv"
)

const licenseCheckBaseURL = "https://gateway.api.bot.or.th/BotLicenseCheckAPI"

type LicenseCheckResponse struct {
	ResultSet     []map[string]interface{} `json:"ResultSet"`
	ResultSetInfo LicenseResultSetInfo     `json:"ResultSetInfo"`
	GroupInfo     []LicenseGroupInfo       `json:"GroupInfo"`
}

type LicenseResultSetInfo struct {
	QueryRecordPerPage int `json:"QueryRecordPerPage"`
	QueryTotalRecord   int `json:"QueryTotalRecord"`
	QueryCurrentPage   int `json:"QueryCurrentPage"`
}

type LicenseGroupInfo struct {
	TypeCode   string `json:"TypeCode"`
	TypeNameTH string `json:"TypeNameTH"`
	Count      int    `json:"Count"`
}

// TypeNameEnglish returns the English translation of the license type.
func (g LicenseGroupInfo) TypeNameEnglish() string {
	switch g.TypeCode {
	case "j":
		return "Legal Entity"
	case "i":
		return "Individual"
	case "b":
		return "Business Establishment"
	default:
		return "All"
	}
}

type AuthorizedDetailResponse struct {
	AuthorizationInfo struct {
		ID             string `json:"Id"`
		AuthorizedName string `json:"AuthorizedName"`
		BranchName     string `json:"BranchName"`
		TypeID         string `json:"TypeId"`
		TypeName       string `json:"TypeName"`
		LastUpdate     string `json:"LastUpdate"`
	} `json:"AuthorizationInfo"`
}

func (c *Client) SearchAuthorized(ctx context.Context, keyword string, page string, limit int) (*LicenseCheckResponse, error) {
	query := url.Values{}
	query.Set("keyword", keyword)
	if page != "" {
		query.Set("page", page)
	}
	if limit > 0 {
		query.Set("limit", strconv.Itoa(limit))
	}

	var result LicenseCheckResponse
	if err := c.requestJSON(ctx, licenseCheckBaseURL, "/SearchAuthorized", query, &result); err != nil {
		return nil, fmt.Errorf("failed to search authorized: %w", err)
	}
	return &result, nil
}

func (c *Client) GetLicense(ctx context.Context, authID, docID string) ([]byte, error) {
	query := url.Values{}
	query.Set("authId", authID)
	query.Set("docId", docID)

	u, _ := url.Parse(licenseCheckBaseURL + "/License")
	u.RawQuery = query.Encode()

	resp, err := c.GetURL(ctx, u.String())
	if err != nil {
		return nil, fmt.Errorf("failed to get license: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read license response: %w", err)
	}
	return data, nil
}

func (c *Client) GetAuthorizedDetail(ctx context.Context, id int) (*AuthorizedDetailResponse, error) {
	query := url.Values{}
	query.Set("id", strconv.Itoa(id))

	var result AuthorizedDetailResponse
	if err := c.requestJSON(ctx, licenseCheckBaseURL, "/AuthorizedDetail", query, &result); err != nil {
		return nil, fmt.Errorf("failed to get authorized detail: %w", err)
	}
	return &result, nil
}
