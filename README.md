# bankofthailand-go

[![Go Reference](https://pkg.go.dev/badge/github.com/jwitmann/bankofthailand-go.svg)](https://pkg.go.dev/github.com/jwitmann/bankofthailand-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/jwitmann/bankofthailand-go)](https://goreportcard.com/report/github.com/jwitmann/bankofthailand-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/go-%3E%3D1.21-blue.svg)](https://golang.org)

**Complete** Go client for **all** [Bank of Thailand (BOT) Public APIs](https://bot-public-api.apigee.io/). Zero external dependencies. Every published BOT API endpoint is implemented.

## Features

- **Holidays** — Financial institution holidays (bilingual Thai/English)
- **Exchange Rates** — Daily/monthly/quarterly/annual THB/FX rates
- **Interest Rates** — Policy rate, BIBOR, deposit/loan rates
- **Economic Statistics** — Time-series observations, category search
- **Debt Securities Auction** — Government bond and SOE bond auction results
- **License Check** — Search BOT-supervised business licenses (P-Loan, Nano Finance, e-Money)
- **Per-Endpoint Authentication** — Different API tokens per service category
- **Path-Aware Rate Limiting** — Automatic rate limit selection by endpoint
- **Retry Logic** — Exponential backoff with configurable status codes
- **CLI Tool** — `bot-holidays` command-line utility

> **Complete API Coverage** — This client implements every published BOT Public API endpoint, including the recently added Debt Securities Auction and License Check catalogues.

## Installation

```bash
go get github.com/jwitmann/bankofthailand-go
```

## Quick Start

### Authentication

Create `config/bot-keys.json` with per-endpoint tokens:

```json
{
  "app_id": "your-app-id",
  "tokens": {
    "others": "your-holidays-token",
    "exchange_rates": "your-exchange-rates-token",
    "interest_rates": "your-interest-rates-token",
    "statistics": "your-statistics-token",
    "debt_security_auction": "your-debt-securities-token",
    "license_check": "your-license-check-token"
  }
}
```

Or use a single token for all endpoints:

```json
{
  "app_id": "your-app-id",
  "token": "your-api-token"
}
```

Or set the environment variable:

```bash
export BOT_API_TOKEN="your-api-token"
```

The client automatically selects the appropriate token based on the endpoint URL.

### Holidays

```go
package main

import (
    "context"
    "fmt"
    "log"

    bot "github.com/jwitmann/bankofthailand-go"
)

func main() {
    client, err := bot.NewClient()
    if err != nil {
        log.Fatal(err)
    }

    holidays, err := client.GetHolidays(context.Background(), 2026)
    if err != nil {
        log.Fatal(err)
    }

    for _, h := range holidays {
        fmt.Printf("%s: %s / %s\n", h.Date, h.HolidayDescription, h.HolidayDescriptionThai)
    }
}
```

### Exchange Rates

```go
resp, err := client.GetDailyAverageExchangeRate(ctx, "2026-06-01", "2026-06-05", "USD")
if err != nil {
    log.Fatal(err)
}

for _, rate := range resp.Result.Data.DataDetail {
    fmt.Printf("%s: Buying=%s, Selling=%s, Mid=%s\n",
        rate.Period, rate.BuyingTransfer, rate.Selling, rate.MidRate)
}
```

### Policy Rate

```go
resp, err := client.GetPolicyRate(ctx)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Policy Rate: %s%% (announced %s)\n",
    resp.Result.Data, resp.Result.AnnouncementDate)
```

### Statistics

```go
// List categories
cats, err := client.GetCategoryList(ctx)

// Search series
search, err := client.SearchSeries(ctx, "GDP")

// Get observations
obs, err := client.GetObservations(ctx, "PF00000000Q00232", "2017-01-01", "2017-12-31", "")
```

### Debt Securities Auction

```go
resp, err := client.GetDebtSecuritiesAuction(ctx, "2017-09-01", "2017-09-30")
if err != nil {
    log.Fatal(err)
}
for _, rec := range resp.Result.Data.DataDetail {
    fmt.Printf("%s %s %s Yield=%s%% Status=%s\n",
        rec.AuctionDate, rec.DebtSecuritiesType, rec.ThaiBMASymbol,
        rec.WeightedAverageAcceptedYield, rec.AuctionStatus)
}
```

### License Check

```go
// Search for authorized entities
search, err := client.SearchAuthorized(ctx, "keyword", "", 10)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Found %d results\n", search.ResultSetInfo.QueryTotalRecord)
for _, g := range search.GroupInfo {
    fmt.Printf("  %s (%s): %d\n", g.TypeNameTH, g.TypeNameEnglish(), g.Count)
}

// Get license PDF document (returns raw PDF bytes)
pdfBytes, err := client.GetLicense(ctx, "authId", "docId")
if err != nil {
    log.Fatal(err)
}
os.WriteFile("license.pdf", pdfBytes, 0644)

// Get authorized entity details
detail, err := client.GetAuthorizedDetail(ctx, 123)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Name: %s\n", detail.AuthorizationInfo.AuthorizedName)
fmt.Printf("Type: %s\n", detail.AuthorizationInfo.TypeName)
```

## Rate Limits

The client automatically applies the correct rate limit based on the endpoint being called:

| API Category | Endpoints | Limit |
|-------------|-----------|-------|
| Holidays | `GetHolidays`, `GetHolidaysRaw` | 100 calls/hour |
| Exchange Rates | `GetDailyAverageExchangeRate`, `GetDailyReferenceRate`, `GetSpotRate`, `GetSwapPoint`, `GetImpliedInterestRate` | 200 calls/hour |
| Interest Rates | `GetPolicyRate`, `GetBIBOR`, `GetDepositRate`, `GetLoanRate`, `GetInterbankTransactionRate` | 200 calls/hour |
| Statistics | `GetCategoryList`, `GetSeriesList`, `GetObservations`, `SearchSeries` | 2000 calls/hour |
| Debt Securities | `GetDebtSecuritiesAuction` | 200 calls/hour |
| License Check | `SearchAuthorized`, `GetLicense`, `GetAuthorizedDetail` | 100 calls/hour |

No configuration needed — the client detects the endpoint from the URL and applies the appropriate limiter automatically.

```go
// Override with a single custom limiter for all endpoints
client, _ := bot.NewClient(
    bot.WithRateLimiter(bot.NewHourlyRateLimiter(500)),
)

// Disable rate limiting
client, _ := bot.NewClient(
    bot.WithRateLimiter(&bot.NoOpRateLimiter{}),
)

// Query limits programmatically
info := bot.GetRateLimitInfo("exchange_rates")
fmt.Printf("Limit: %d calls/hour, Quota: %s\n", info.CallsPerHour, info.Quota)
```

## Configuration Options

```go
client, _ := bot.NewClient(
    bot.WithToken("your-token"),                              // Single token for all endpoints
    bot.WithHTTPClient(&http.Client{Timeout: 60 * time.Second}), // Custom HTTP client
    bot.WithRateLimiter(bot.NewHourlyRateLimiter(500)),        // Override all limiters
    bot.WithRetryPolicy(bot.DefaultRetryPolicy()),              // Custom retry
    bot.WithConfigPath("/path/to/config.json"),                // Custom config path
)
```

## Error Handling

```go
holidays, err := client.GetHolidays(ctx, 2026)
if err != nil {
    if errors.Is(err, bot.ErrNoContent) {
        // HTTP 204 — data not yet available for this year
        log.Println("Holiday data not yet available")
        return
    }
    var apiErr *bot.APIError
    if errors.As(err, &apiErr) {
        log.Fatalf("API error: status=%d, message=%s", apiErr.StatusCode, apiErr.Message)
    }
    log.Fatalf("Request failed: %v", err)
}
```

## CLI

```bash
# Holidays (raw list)
go run ./cmd/bot-holidays -year 2026

# Holidays (ThaiFA-compatible format with API wrapper)
go run ./cmd/bot-holidays -year 2026 -format thaifa

# Holidays (CSV)
go run ./cmd/bot-holidays -year 2026 -format csv

# Install
make install
bot-holidays -year 2026
```

## API Documentation

See [docs/API.md](docs/API.md) for full endpoint documentation, response schemas, and examples.

## Development

```bash
# Format code
make fmt

# Run tests
make test

# Run linter
make lint

# Run all quality gates
make check
```

## License

MIT License — see [LICENSE](LICENSE) for details.

## Disclaimer

This is an unofficial client. Bank of Thailand APIs are subject to their terms of service. API access requires registration at the [BOT Developer Portal](https://bot-public-api.apigee.io/).
