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

---

## Client

### Creating a Client

```go
client, err := bot.NewClient(options...)
```

**Options:**

| Option | Description |
|--------|-------------|
| `WithToken(token)` | Set API token directly |
| `WithHTTPClient(client)` | Use custom `*http.Client` |
| `WithBaseURL(url)` | Override base URL (for testing) |
| `WithRateLimiter(limiter)` | Set custom rate limiter |
| `WithRetryPolicy(policy)` | Set custom retry policy |
| `WithTimeout(duration)` | Set HTTP timeout |
| `WithConfigPath(path)` | Set config file path |

### Default Behavior

- Loads token from `BOT_API_TOKEN` env var or `config/bot-keys.json`
- Rate limiter: 100 calls/hour (holidays default)
- HTTP timeout: 30 seconds
- Retry: 3 retries with exponential backoff on 5xx errors

---

## Authentication

The client supports two authentication methods:

### Environment Variable

```bash
export BOT_API_TOKEN="your-api-token"
```

### Config File

Create `config/bot-keys.json`:

```json
{
  "app_id": "your-app-id",
  "token": "your-api-token"
}
```

The token is sent as an `Authorization` header on every request:

```
Authorization: your-api-token
```

---

## Rate Limiting

Built-in token bucket rate limiter with per-API limits:

### Limits

| Endpoint Category | Calls/Hour | Quota |
|------------------|-----------|-------|
| Holidays | 100 | unlimited |
| Exchange Rates | 200 | unlimited |
| Interest Rates | 200 | unlimited |
| Statistics | 2000 | unlimited |

### Usage

```go
// Use default limiter (100/hour)
client, _ := bot.NewClient()

// Use category-specific limiter
client, _ := bot.NewClient(
    bot.WithRateLimiter(bot.NewRateLimiterForExchangeRates()), // 200/hour
)

// Custom limit
client, _ := bot.NewClient(
    bot.WithRateLimiter(bot.NewHourlyRateLimiter(500)),
)

// Disable limiter
client, _ := bot.NewClient(
    bot.WithRateLimiter(&bot.NoOpRateLimiter{}),
)

// Query limits
info := bot.GetRateLimitInfo("statistics")
fmt.Printf("%d calls/hour, %s quota\n", info.CallsPerHour, info.Quota)
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
    var apiErr *bot.APIError
    if errors.As(err, &apiErr) {
        fmt.Printf("API Error: status=%d, message=%s\n",
            apiErr.StatusCode, apiErr.Message)
    }
}
```

**APIError Fields:**

| Field | Type | Description |
|-------|------|-------------|
| `StatusCode` | `int` | HTTP status code |
| `Message` | `string` | Response body or status text |
| `URL` | `string` | Request URL |

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
