# API Documentation

Complete reference for the Bank of Thailand Go client.

## Table of Contents

- [Client](#client)
- [Authentication](#authentication)
- [Rate Limiting](#rate-limiting)
- [Retry Logic](#retry-logic)
- [Error Handling](#error-handling)
- [Holidays](#holidays)
- [Exchange Rates](#exchange-rates)
- [Interest Rates](#interest-rates)
- [Statistics](#statistics)
- [Debt Securities Auction](#debt-securities-auction)
- [License Check](#license-check)

---

## Client

### Creating a Client

```go
client, err := bot.NewClient(options...)
```

**Options:**

| Option | Description |
|--------|-------------|
| `WithToken(token)` | Set single API token for all endpoints |
| `WithHTTPClient(client)` | Use custom `*http.Client` |
| `WithBaseURL(url)` | Override base URL (for testing) |
| `WithRateLimiter(limiter)` | Override all endpoint limiters with one |
| `WithRetryPolicy(policy)` | Set custom retry policy |
| `WithTimeout(duration)` | Set HTTP timeout |
| `WithConfigPath(path)` | Set config file path |

### Default Behavior

- Loads credentials from `BOT_API_TOKEN` env var or `config/bot-keys.json`
- Supports per-endpoint tokens via `tokens` map in config
- Automatically selects rate limiter based on endpoint URL:
  - Holidays: 100 calls/hour
  - Exchange Rates: 200 calls/hour
  - Interest Rates: 200 calls/hour
  - Statistics: 2000 calls/hour
- HTTP timeout: 30 seconds
- Retry: 3 retries with exponential backoff on 5xx errors

---

## Authentication

The client supports three authentication methods:

### Per-Endpoint Tokens (Recommended)

Create `config/bot-keys.json` with separate tokens per service:

```json
{
  "app_id": "your-app-id",
  "tokens": {
    "others": "your-holidays-token",
    "exchange_rates": "your-exchange-rates-token",
    "interest_rates": "your-interest-rates-token",
    "statistics": "your-statistics-token"
  }
}
```

The client automatically selects the appropriate token based on the endpoint URL:

| Endpoint Key | Matching URL Patterns |
|-------------|----------------------|
| `others` | `financial-institutions-holidays` |
| `exchange_rates` | `Stat-ExchangeRate`, `Stat-ReferenceRate`, `Stat-SpotRate`, `Stat-SwapPoint`, `Stat-ThaiBahtImpliedInterestRate` |
| `interest_rates` | `PolicyRate`, `BIBOR`, `DepositRate`, `LoanRate`, `Stat-InterbankTransactionRate` |
| `statistics` | `categorylist`, `serieslist`, `observations`, `search-series` |
| `debt_security_auction` | `BondAuction` |
| `license_check` | `BotLicenseCheckAPI` |

### Single Token

```json
{
  "app_id": "your-app-id",
  "token": "your-api-token"
}
```

### Environment Variable

```bash
export BOT_API_TOKEN="your-api-token"
```

The token is sent as an `Authorization` header on every request:

```
Authorization: your-api-token
```

### Token Inspection

```go
// Check if a per-endpoint token is configured
if client.HasEndpointToken("exchange_rates") {
    // Uses separate token for exchange rate APIs
}

// List configured endpoint token keys
keys := client.EndpointTokenKeys()
```

---

## Rate Limiting

Built-in token bucket rate limiter with **automatic per-API selection**:

### Limits

| Endpoint Category | Calls/Hour | Matching Endpoints |
|------------------|-----------|-------------------|
| Holidays | 100 | `GetHolidays`, `GetHolidaysRaw` |
| Exchange Rates | 200 | `GetDailyAverageExchangeRate`, `GetDailyReferenceRate`, `GetSpotRate`, `GetSwapPoint`, `GetImpliedInterestRate` |
| Interest Rates | 200 | `GetPolicyRate`, `GetBIBOR`, `GetDepositRate`, `GetLoanRate`, `GetInterbankTransactionRate` |
| Statistics | 2000 | `GetCategoryList`, `GetSeriesList`, `GetObservations`, `SearchSeries` |
| Debt Securities | 200 | `GetDebtSecuritiesAuction` |
| License Check | 100 | `SearchAuthorized`, `GetLicense`, `GetAuthorizedDetail` |

### Usage

```go
// Default: automatic per-endpoint rate limiting
client, _ := bot.NewClient()

// Override with a single limiter for all endpoints
client, _ := bot.NewClient(
    bot.WithRateLimiter(bot.NewHourlyRateLimiter(500)),
)

// Disable limiting
client, _ := bot.NewClient(
    bot.WithRateLimiter(&bot.NoOpRateLimiter{}),
)

// Query limits
info := bot.GetRateLimitInfo("statistics")
fmt.Printf("%d calls/hour, %s quota\n", info.CallsPerHour, info.Quota)
```

### Custom Limiter

```go
limiter := bot.NewHourlyRateLimiter(1000) // 1000 calls/hour
client, _ := bot.NewClient(bot.WithRateLimiter(limiter))
```

---

## Retry Logic

Exponential backoff with configurable retry conditions.

### Default Policy

- Max retries: 3
- Retry on: 5xx status codes, network errors
- Backoff: exponential (base delay × attempt)

### Custom Policy

```go
policy := &bot.RetryPolicy{
    MaxRetries: 5,
    ShouldRetry: func(err error) bool {
        return true // retry all errors
    },
    ShouldRetryStatus: func(code int) bool {
        return code >= 500 || code == 429
    },
    Backoff: func(attempt int) time.Duration {
        return time.Duration(attempt) * time.Second
    },
}

client, _ := bot.NewClient(bot.WithRetryPolicy(policy))
```

---

## Error Handling

All API errors return `*APIError`:

```go
resp, err := client.GetHolidays(ctx, 2026)
if err != nil {
    if errors.Is(err, bot.ErrNoContent) {
        // HTTP 204 — data not yet available
        fmt.Println("Data not yet available for this period")
        return
    }
    var apiErr *bot.APIError
    if errors.As(err, &apiErr) {
        fmt.Printf("API Error: status=%d, message=%s\n",
            apiErr.StatusCode, apiErr.Message)
    }
}
```

**Error Types:**

| Error | Description |
|-------|-------------|
| `*APIError` | HTTP 4xx/5xx errors |
| `ErrNoContent` | HTTP 204 — requested data not yet available |

**APIError Fields:**

| Field | Type | Description |
|-------|------|-------------|
| `StatusCode` | `int` | HTTP status code |
| `Message` | `string` | Response body or status text |

---

## Holidays

### GetHolidays

```go
func (c *Client) GetHolidays(ctx context.Context, year int) ([]Holiday, error)
```

Fetch financial institution holidays for a given year.

**Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `year` | `int` | Yes | Year (2000–2100) |

**Response:** `[]Holiday`

| Field | Type | Description |
|-------|------|-------------|
| `Date` | `string` | Date in YYYY-MM-DD format |
| `DateThai` | `string` | Date in DD/MM/YYYY (Buddhist era) |
| `HolidayWeekDay` | `string` | Day of week (English) |
| `HolidayWeekDayThai` | `string` | Day of week (Thai) |
| `HolidayDescription` | `string` | Holiday name (English) |
| `HolidayDescriptionThai` | `string` | Holiday name (Thai) |

**Example:**

```go
holidays, err := client.GetHolidays(ctx, 2026)
for _, h := range holidays {
    fmt.Printf("%s: %s\n", h.Date, h.HolidayDescription)
}
// Output:
// 2026-01-01: New Year's Day
// 2026-04-13: Songkran Festival
```

### GetHolidaysRaw

```go
func (c *Client) GetHolidaysRaw(ctx context.Context, year int) (*HolidaysResponse, error)
```

Fetch holidays with the full API response wrapper (includes `api` name and `timestamp`).

**Response:** `*HolidaysResponse`

```go
type HolidaysResult struct {
    API       string    `json:"api"`       // "API_V2.FIHolidays"
    Timestamp string    `json:"timestamp"` // "2026-06-06 10:30:15"
    Data      []Holiday `json:"data"`
}
```

**Example:**

```go
resp, err := client.GetHolidaysRaw(ctx, 2026)
if err != nil {
    if errors.Is(err, bot.ErrNoContent) {
        fmt.Println("Holiday data not yet available for 2026")
        return
    }
    log.Fatal(err)
}
fmt.Printf("API: %s, Timestamp: %s, Holidays: %d\n",
    resp.Result.API, resp.Result.Timestamp, len(resp.Result.Data))
```

---

## Exchange Rates

### GetDailyAverageExchangeRate

```go
func (c *Client) GetDailyAverageExchangeRate(ctx context.Context, startPeriod, endPeriod, currency string) (*ExchangeRateResponse, error)
```

Daily average exchange rate (THB vs 19 currencies).

**Parameters:**

| Parameter | Format | Required | Description |
|-----------|--------|----------|-------------|
| `start_period` | `YYYY-MM-DD` | Yes | Start date |
| `end_period` | `YYYY-MM-DD` | Yes | End date |
| `currency` | string | No | Currency code (e.g., `USD`, `EUR`) |

**Response:** `*ExchangeRateResponse`

```go
type ExchangeRateData struct {
    Period          string // "2026-06-05"
    CurrencyID      string // "USD"
    CurrencyNameTh  string // "ดอลลาร์สหรัฐ"
    CurrencyNameEng string // "US DOLLAR"
    BuyingSight     string // "33.50"
    BuyingTransfer  string // "33.55"
    Selling         string // "34.00"
    MidRate         string // "33.775"
}
```

### GetDailyReferenceRate

```go
func (c *Client) GetDailyReferenceRate(ctx context.Context, startPeriod, endPeriod string) (*ReferenceRateResponse, error)
```

Weighted-average interbank exchange rate (THB/USD only).

**Parameters:**

| Parameter | Format | Required |
|-----------|--------|----------|
| `start_period` | `YYYY-MM-DD` | Yes |
| `end_period` | `YYYY-MM-DD` | Yes |

**Response:** `*ReferenceRateResponse`

```go
type ReferenceRateData struct {
    Period string // "2026-06-05"
    Rate   string // "33.77"
}
```

### GetSpotRate

```go
func (c *Client) GetSpotRate(ctx context.Context, startPeriod, endPeriod string) (*SpotRateResponse, error)
```

Spot rate bid/offer for USD/THB.

**Response:** `*SpotRateResponse`

```go
type SpotRateData struct {
    Period    string // "2026-06-05"
    BidRate   string // "33.75"
    OfferRate string // "33.79"
}
```

### GetSwapPoint

```go
func (c *Client) GetSwapPoint(ctx context.Context, startPeriod, endPeriod, termType string) (*SwapPointResponse, error)
```

Onshore swap points in satangs.

**Parameters:**

| Parameter | Required | Description |
|-----------|----------|-------------|
| `term_type` | No | Term type (e.g., `1M`, `3M`, `6M`) |

### GetImpliedInterestRate

```go
func (c *Client) GetImpliedInterestRate(ctx context.Context, startPeriod, endPeriod, rateType string) (*ImpliedRateResponse, error)
```

Thai Baht implied interest rates from swap market.

**Parameters:**

| Parameter | Required | Description |
|-----------|----------|-------------|
| `rate_type` | No | Rate type |

---

## Interest Rates

### GetPolicyRate

```go
func (c *Client) GetPolicyRate(ctx context.Context) (*PolicyRateResponse, error)
```

Current monetary policy rate.

**Response:** `*PolicyRateResponse`

| Field | Type | Description |
|-------|------|-------------|
| `Data` | `string` | Rate (e.g., `"2.50"`) |
| `AnnouncementDate` | `string` | MPC meeting date |
| `NewsTextEn` | `string` | English announcement text |
| `NewsTextTh` | `string` | Thai announcement text |
| `EffectiveDateTime` | `string` | Effective date/time |

### GetBIBOR

```go
func (c *Client) GetBIBOR(ctx context.Context, startPeriod, endPeriod, bank string) (*BIBORResponse, error)
```

Bangkok Interbank Offered Rate by bank.

**Parameters:**

| Parameter | Required | Description |
|-----------|----------|-------------|
| `bank` | No | Bank name filter |

**Response:** `*BIBORResponse`

```go
type BIBORData struct {
    Period      string // "2026-06-05"
    BankNameTh  string // "ธนาคารกสิกรไทย"
    BankNameEng string // "KASIKORNBANK"
    BIBORON     string // Overnight
    BIBOR1W     string // 1 week
    BIBOR1M     string // 1 month
    BIBOR2M     string // 2 months
    BIBOR3M     string // 3 months
    BIBOR6M     string // 6 months
    BIBOR9M     string // 9 months
    BIBOR1Y     string // 1 year
}
```

### GetBIBORAverage

```go
func (c *Client) GetBIBORAverage(ctx context.Context, startPeriod, endPeriod string) (*BIBORResponse, error)
```

Average BIBOR across all contributing banks.

### GetDepositRate

```go
func (c *Client) GetDepositRate(ctx context.Context, startPeriod, endPeriod string) (*DepositRateResponse, error)
```

Deposit interest rates for individuals (saving, fixed term).

**Response:** `*DepositRateResponse`

```go
type DepositRateData struct {
    Period       string // "2026-06-05"
    BankNameEng  string
    SavingMin    string // Min saving rate
    SavingMax    string // Max saving rate
    Fix3MthsMin  string // 3-month fixed min
    Fix3MthsMax  string // 3-month fixed max
    Fix6MthsMin  string // 6-month fixed min
    Fix6MthsMax  string // 6-month fixed max
    Fix12MthsMin string // 12-month fixed min
    Fix12MthsMax string // 12-month fixed max
    Fix24MthsMin string // 24-month fixed min
    Fix24MthsMax string // 24-month fixed max
}
```

### GetLoanRate

```go
func (c *Client) GetLoanRate(ctx context.Context, startPeriod, endPeriod string) (*LoanRateResponse, error)
```

Loan interest rates (MOR, MLR, MRR, ceiling, default).

**Response:** `*LoanRateResponse`

```go
type LoanRateData struct {
    Period        string // "2026-06-05"
    BankNameEng   string
    MOR           string // Minimum Overdraft Rate
    MLR           string // Minimum Loan Rate
    MRR           string // Minimum Retail Rate
    CeilingRate   string
    DefaultRate   string
    CreditCardMin string
    CreditCardMax string
}
```

### GetInterbankTransactionRate

```go
func (c *Client) GetInterbankTransactionRate(ctx context.Context, startPeriod, endPeriod, termType string) (*InterbankTransactionRateResponse, error)
```

Interbank transaction rates by tenor.

**Parameters:**

| Parameter | Required | Description |
|-----------|----------|-------------|
| `term_type` | No | O/N, T/N, fixed term, etc. |

**Response:** `*InterbankTransactionRateResponse`

```go
type InterbankTransactionRateData struct {
    Period                      string
    TermTypeNameEng             string
    MinInterestRate             string
    MaxInterestRate             string
    ModeInterestRate            string
    WeightedAverageInterestRate string
}
```

---

## Statistics

### GetCategoryList

```go
func (c *Client) GetCategoryList(ctx context.Context) (*CategoryListResponse, error)
```

List all available statistical categories.

**Response:** `*CategoryListResponse`

```go
type Category struct {
    Category       string // "EC_XT_077"
    DescriptionTh  string // "อัตราแลกเปลี่ยน"
    DescriptionEng string // "Exchange Rates"
}
```

### GetSeriesList

```go
func (c *Client) GetSeriesList(ctx context.Context, category string) (*SeriesListResponse, error)
```

List series within a category.

**Parameters:**

| Parameter | Required | Description |
|-----------|----------|-------------|
| `category` | Yes | Category code |

### GetObservations

```go
func (c *Client) GetObservations(ctx context.Context, seriesCode, startPeriod, endPeriod, sortBy string) (*ObservationsResponse, error)
```

Fetch time-series observations.

**Parameters:**

| Parameter | Format | Required | Description |
|-----------|--------|----------|-------------|
| `series_code` | string | Yes | Series identifier |
| `start_period` | `YYYY-MM-DD` | Yes | Start date |
| `end_period` | `YYYY-MM-DD` | No | End date |
| `sort_by` | `asc`/`desc` | No | Sort order |

**Response:** `*ObservationsResponse`

```go
type ObservationSeries struct {
    SeriesCode     string
    SeriesNameEng  string
    UnitEng        string // "Million Baht"
    Frequency      string // "Quarterly"
    Observations   []Observation
}

type Observation struct {
    PeriodStart string // "2017-Q1"
    Value       string // "8648519.0000000"
}
```

### SearchSeries

```go
func (c *Client) SearchSeries(ctx context.Context, keyword string) (*SearchResponse, error)
```

Search across all statistical series.

**Parameters:**

| Parameter | Required | Description |
|-----------|----------|-------------|
| `keyword` | Yes | Search term |

**Response:** `*SearchResponse`

```go
type SeriesDetail struct {
    SeriesCode         string
    SeriesNameEng      string
    SeriesCategories   string
    Frequency          string
    UnitTh             string
    LastUpdatedDate    string
    SourceOfDataEng    string
    ReleaseScheduleEng string
    DescriptionEng     string
}
```

---

## Debt Securities Auction

### GetDebtSecuritiesAuction

```go
func (c *Client) GetDebtSecuritiesAuction(ctx context.Context, startPeriod, endPeriod string) (*DebtSecuritiesResponse, error)
```

Fetch government and SOE bond auction results.

**Parameters:**

| Parameter | Format | Required | Description |
|-----------|--------|----------|-------------|
| `start_period` | `YYYY-MM-DD` | Yes | Start auction date |
| `end_period` | `YYYY-MM-DD` | Yes | End auction date |

**Response:** `*DebtSecuritiesResponse`

```go
type DebtSecuritiesRecord struct {
    AuctionDate                     string // "2017-09-26"
    DebtSecuritiesType              string // "Government Bonds"
    ThaiBMASymbol                   string // "LB233A"
    ISINCode                        string // "TH0623033303"
    AuctionNameTh                   string // Thai auction name
    CFICode                         string // "DBFTFR"
    CouponRate                      string // "5.5"
    TimeToMaturity                  string // "5.46 Yrs"
    PaymentDate                     string // "2017-09-28"
    StartDateOfInterestEarningPeriod string // "2017-09-13"
    MaturityDate                    string // "2023-03-13"
    IssueAmountNCB_CB               string // "2000.0000000"
    AcceptedAmountNCB_CB            string // "2000.0000000"
    AcceptedAmountNCB               string // ""
    AcceptedAmountCB                string // "2000.0000000"
    GreenshoeOptionAmount           string // "400.0000000"
    PAOAmount                       string // ""
    OverAllotmentAmount             string // ""
    GrandTotalAmount                string // "2400.0000000"
    AcceptedLowestYield             string // "1.7070000"
    AcceptedHighestYield            string // "1.7090000"
    WeightedAverageAcceptedYield    string // "1.7077000"
    BidCoverageRatio                string // "2.2000000"
    AuctionStatus                   string // "Approve"
}
```

---

## License Check

### SearchAuthorized

```go
func (c *Client) SearchAuthorized(ctx context.Context, keyword string, page string, limit int) (*LicenseCheckResponse, error)
```

Search for BOT-supervised business licenses (P-Loan, Nano Finance, e-Money, etc.).

**Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `keyword` | `string` | Yes | Search keyword |
| `page` | `string` | No | Page position |
| `limit` | `int` | No | Results per page |

**Response:** `*LicenseCheckResponse`

```go
type LicenseCheckResponse struct {
    ResultSet     []map[string]interface{} // Varying license record fields
    ResultSetInfo LicenseResultSetInfo     // Pagination info
    GroupInfo     []LicenseGroupInfo       // Category breakdown
}

type LicenseGroupInfo struct {
    TypeCode   string // "j", "i", "b", or ""
    TypeNameTH string // "นิติบุคคล", "บุคคล", "สถานประกอบการ", "ทั้งหมด"
    Count      int
}
```

**Translation:**

```go
for _, g := range resp.GroupInfo {
    fmt.Printf("%s (%s): %d\n", g.TypeNameTH, g.TypeNameEnglish(), g.Count)
}
// Output:
// นิติบุคคล (Legal Entity): 42
// บุคคล (Individual): 5
// สถานประกอบการ (Business Establishment): 12
```

### GetLicense

```go
func (c *Client) GetLicense(ctx context.Context, authID, docID string) ([]byte, error)
```

Download a license document as a **PDF**.

**Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `authId` | `string` | Yes | Authorized entity ID |
| `docId` | `string` | Yes | Document reference number |

**Returns:** `[]byte` containing the PDF document.

**Example:**

```go
pdfBytes, err := client.GetLicense(ctx, "12345", "DOC-2024-001")
if err != nil {
    log.Fatal(err)
}
os.WriteFile("license.pdf", pdfBytes, 0644)
```

### GetAuthorizedDetail

```go
func (c *Client) GetAuthorizedDetail(ctx context.Context, id int) (*AuthorizedDetailResponse, error)
```

Fetch detailed information about an authorized entity.

**Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `id` | `int` | Yes | Entity ID |

**Response:** `*AuthorizedDetailResponse`

```go
type AuthorizedDetailResponse struct {
    AuthorizationInfo struct {
        ID             string // "123"
        AuthorizedName string // "บริษัท ... จำกัด"
        BranchName     string // ""
        TypeID         string // "ผู้ประกอบธุรกิจ..."
        TypeName       string // "ผู้ประกอบธุรกิจ..."
        LastUpdate     string // "10/06/2026"
    }
}
```

---

## Common Response Patterns

### Result Wrapper

Most responses follow this pattern:

```json
{
  "result": {
    "api": "API_NAME",
    "timestamp": "2026-06-06 12:00:00",
    "data": { ... }
  }
}
```

### Data Header

Exchange rate and statistics APIs include metadata:

```go
type DataHeader struct {
    ReportNameEng    string         // Report name (English)
    ReportNameTh     string         // Report name (Thai)
    ReportUOQNameEng string         // Unit of quantity (English)
    ReportUOQNameTh  string         // Unit of quantity (Thai)
    SourceOfData     []SourceOfData // Data sources
    Remarks          []Remark       // Report remarks
    LastUpdated      string         // Last update timestamp
}
```

---

## Response Types Summary

| Method | Return Type | Base URL |
|--------|-------------|----------|
| `GetHolidays` | `[]Holiday` | `gateway.api.bot.or.th/financial-institutions-holidays` |
| `GetHolidaysRaw` | `*HolidaysResponse` | `gateway.api.bot.or.th/financial-institutions-holidays` |
| `GetDailyAverageExchangeRate` | `*ExchangeRateResponse` | `gateway.api.bot.or.th/Stat-ExchangeRate/v2` |
| `GetDailyReferenceRate` | `*ReferenceRateResponse` | `gateway.api.bot.or.th/Stat-ReferenceRate/v2` |
| `GetSpotRate` | `*SpotRateResponse` | `gateway.api.bot.or.th/Stat-SpotRate/v2/SPOTRATE` |
| `GetSwapPoint` | `*SwapPointResponse` | `gateway.api.bot.or.th/Stat-SwapPoint/v2/SWAPPOINT` |
| `GetImpliedInterestRate` | `*ImpliedRateResponse` | `gateway.api.bot.or.th/Stat-ThaiBahtImpliedInterestRate/v2` |
| `GetPolicyRate` | `*PolicyRateResponse` | `gateway.api.bot.or.th/PolicyRate/v3/policy_rate` |
| `GetBIBOR` | `*BIBORResponse` | `gateway.api.bot.or.th/BIBOR/v2` |
| `GetDepositRate` | `*DepositRateResponse` | `gateway.api.bot.or.th/DepositRate/v2` |
| `GetLoanRate` | `*LoanRateResponse` | `gateway.api.bot.or.th/LoanRate/v2` |
| `GetInterbankTransactionRate` | `*InterbankTransactionRateResponse` | `gateway.api.bot.or.th/Stat-InterbankTransactionRate/v2` |
| `GetCategoryList` | `*CategoryListResponse` | `gateway.api.bot.or.th/categorylist` |
| `GetObservations` | `*ObservationsResponse` | `gateway.api.bot.or.th/observations` |
| `SearchSeries` | `*SearchResponse` | `gateway.api.bot.or.th/search-series` |
| `GetDebtSecuritiesAuction` | `*DebtSecuritiesResponse` | `gateway.api.bot.or.th/BondAuction/bond_auction_v2` |
| `SearchAuthorized` | `*LicenseCheckResponse` | `gateway.api.bot.or.th/BotLicenseCheckAPI` |
| `GetLicense` | `[]byte` (PDF) | `gateway.api.bot.or.th/BotLicenseCheckAPI` |
| `GetAuthorizedDetail` | `*AuthorizedDetailResponse` | `gateway.api.bot.or.th/BotLicenseCheckAPI` |
