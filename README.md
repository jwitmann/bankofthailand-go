# bankofthailand-go

[![Go Reference](https://pkg.go.dev/badge/github.com/jwitmann/bankofthailand-go.svg)](https://pkg.go.dev/github.com/jwitmann/bankofthailand-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/jwitmann/bankofthailand-go)](https://goreportcard.com/report/github.com/jwitmann/bankofthailand-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/go-%3E%3D1.21-blue.svg)](https://golang.org)

Official Go client for [Bank of Thailand (BOT) Public APIs](https://bot-public-api.apigee.io/). Zero external dependencies.

## Features

- **Holidays** — Financial institution holidays (bilingual Thai/English)
- **Exchange Rates** — Daily/monthly/quarterly/annual THB/FX rates
- **Interest Rates** — Policy rate, BIBOR, deposit/loan rates
- **Economic Statistics** — Time-series observations, category search
- **Rate Limiting** — Built-in token bucket with per-API limits
- **Retry Logic** — Exponential backoff with configurable status codes
- **CLI Tool** — `bot-holidays` command-line utility

## Installation

```bash
go get github.com/jwitmann/bankofthailand-go
```

## Quick Start

### Authentication

Create `config/bot-keys.json`:

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

## Rate Limits

Per-API rate limits enforced by the client:

| API Category | Limit |
|-------------|-------|
| Holidays | 100 calls/hour |
| Exchange Rates | 200 calls/hour |
| Interest Rates | 200 calls/hour |
| Statistics | 2000 calls/hour |

```go
// Use a specific rate limiter
client, _ := bot.NewClient(
    bot.WithRateLimiter(bot.NewRateLimiterForExchangeRates()),
)

// Or check limits programmatically
info := bot.GetRateLimitInfo("exchange_rates")
fmt.Printf("Limit: %d calls/hour, Quota: %s\n", info.CallsPerHour, info.Quota)
```

## Configuration Options

```go
client, _ := bot.NewClient(
    bot.WithToken("your-token"),                              // Direct token
    bot.WithHTTPClient(&http.Client{Timeout: 60 * time.Second}), // Custom HTTP client
    bot.WithRateLimiter(bot.NewHourlyRateLimiter(500)),        // Custom rate limit
    bot.WithRetryPolicy(bot.DefaultRetryPolicy()),              // Custom retry
    bot.WithConfigPath("/path/to/config.json"),                // Custom config path
)
```

## CLI

```bash
# Holidays
go run ./cmd/bot-holidays -year 2026
go run ./cmd/bot-holidays -year 2026 -format thaifa
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
