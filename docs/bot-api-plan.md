# BOT API Go Client — Project Plan

**Status:** Draft  
**Date:** 2026-06-06  
**Author:** AI Assistant  
**Project:** Standalone Go client for Bank of Thailand (BOT) public APIs

---

## 1. Problem Statement

ThaiFA currently maintains Thailand public holiday data via manually extracted JSON files (`data/thailand-holidays-2026.json`). This is error-prone and requires manual updates each year.

The Bank of Thailand provides a public API with free endpoints including:
- Financial institution holidays
- Exchange rates (THB/USD, EUR, JPY, etc.)
- Interest rates (policy rate, bond yields, BIBOR)
- Economic statistics

No production-ready Go client exists for this API.

**Goal:** Build a standalone, reusable Go client library for the BOT public API.

---

## 2. Naming

**Rejected:** `bot-go` — ambiguous with chatbots/automation bots.

**Candidates:**
| Name | Pros | Cons |
|---|---|---|
| `bankofthailand-go` | Unambiguous, clear | Long |
| `thai-central-bank` | No "bot" | Less discoverable |
| `bot-or-th-go` | Uses actual domain | Still has "bot" |
| `thb-data` | Short, currency ref | Too generic |

**Recommendation:** `bankofthailand-go`
- Module: `github.com/jwitmann/bankofthailand-go`
- Clear, unambiguous, follows Go naming conventions

---

## 3. API Coverage (Phased)

### Phase 1: Financial Holidays (MVP)
**Endpoint:** `/api/holidays` (exact path TBD from docs)
**Rationale:** Immediate value — replaces manual extraction in ThaiFA
**Priority:** P0

**Data:**
- Date (YYYY-MM-DD)
- Description (Thai + English)
- Holiday type (public, bank, substitution)

### Phase 2: Exchange Rates
**Endpoints:**
- `/api/exchange-rate/avg` (average daily rates)
- `/api/exchange-rate/interbank` (weighted average)

**Rationale:** FX context for foreign equity funds
**Priority:** P1

**Data:**
- THB/USD, THB/EUR, THB/JPY, THB/CNY, THB/GBP
- Daily close
- Historical series

### Phase 3: Interest Rates
**Endpoints:**
- `/api/interest-rate/policy`
- `/api/interest-rate/bond-yields`
- `/api/interest-rate/bibor`

**Rationale:** Bond fund context
**Priority:** P2

**Data:**
- Policy rate (current + historical)
- Government bond yields (1Y, 2Y, 5Y, 10Y, 20Y)
- BIBOR (1M, 3M, 6M)

### Phase 4: Economic Statistics
**Endpoints:** TBD from `/api/statistics/...`

**Rationale:** Macro context
**Priority:** P3 (future)

---

## 4. Technical Design

### Architecture (mirrors `sec-go`)

```
bankofthailand-go/
├── client.go           # HTTP client, auth, rate limiting
├── options.go          # Functional options
├── error.go            # Error types
├── rate.go             # Rate limiter
├── retry.go            # Exponential backoff
├── models.go           # All API response structs
├── holidays.go         # Holiday-specific endpoints
├── exchange_rates.go   # FX endpoints
├── interest_rates.go   # Interest rate endpoints
│
├── client_test.go      # Unit tests (mock server)
├── holidays_test.go    # Holiday endpoint tests
│
├── cmd/
│   └── bot-holidays/   # CLI: fetch holidays for year(s)
│
├── docs/
│   └── API.md          # Endpoint documentation
│
├── go.mod
├── Makefile
└── README.md
```

### Key Features

1. **Zero external dependencies** (stdlib only, like `sec-go`)
2. **Built-in rate limiting** — respect BOT limits
3. **Exponential backoff retry** — handle transient failures
4. **Bilingual support** — Thai/English field mapping
5. **Optional TTL cache** — in-memory caching for repeated calls
6. **Pagination support** — for historical data series

### Auth

BOT API uses OAuth2 or API key (need to verify from docs).

**Primary:** API key from environment `BOT_API_KEY`  
**Secondary:** OAuth2 token flow (if required)

### Rate Limits

TBD from BOT API documentation. Estimate: ~1000 requests/day for public tier.

---

## 5. CLI Tools

### `bot-holidays`

```bash
# Fetch holidays for current year
bot-holidays

# Fetch for specific year
bot-holidays -year 2027

# Output as JSON (ThaiFA-compatible format)
bot-holidays -year 2026 -format thaifa

# Output as CSV
bot-holidays -year 2026 -format csv
```

**Output format (ThaiFA-compatible):**
```json
{
  "result": {
    "api": "API_V2.FIHolidays",
    "timestamp": "2025-05-11 16:41:07",
    "data": [
      {
        "HolidayWeekDay": "Thursday",
        "HolidayWeekDayThai": "วันพฤหัสบดี",
        "Date": "2026-01-01",
        "DateThai": "01/01/2569",
        "HolidayDescription": "New Year's Day",
        "HolidayDescriptionThai": "วันขึ้นปีใหม่"
      }
    ]
  }
}
```

---

## 6. Integration with ThaiFA

### Phase 1: Auto-fetch holidays
1. Add `bankofthailand-go` dependency to ThaiFA
2. Replace `LoadThailandHolidays()` with BOT API call
3. Cache result locally (same JSON format)
4. Auto-refresh when new year data is needed

### Phase 2: FX context
1. Fetch daily THB/USD rate
2. Display in fund info: "FX impact: +1.2% this week"

---

## 7. Development Plan

| Week | Task | Deliverable |
|---|---|---|
| 1 | Research BOT API docs, register for API key | API documentation, test credentials |
| 2 | Build core client (HTTP, auth, rate limiting, retry) | `client.go`, `options.go`, `error.go` |
| 3 | Implement holidays endpoint + tests | `holidays.go`, `client_test.go` |
| 4 | Build `bot-holidays` CLI | Working CLI tool |
| 5 | Integrate with ThaiFA | ThaiFA uses BOT API for holidays |
| 6 | Implement exchange rates | `exchange_rates.go` |
| 7 | Implement interest rates | `interest_rates.go` |
| 8 | Polish, docs, release v1.0 | README, docs, tagged release |

---

## 8. Open Questions

1. **API key registration:** Where and how to register for BOT API access?
2. **Rate limits:** What are the actual limits for public tier?
3. **Historical holidays:** Does BOT API provide past years, or only current/future?
4. **Holiday granularity:** Financial institution holidays only, or all public holidays?
5. **Data freshness:** When are next year's holidays published?

---

## 9. Success Criteria

- [ ] Can fetch 2026 holidays via API (matches manual JSON exactly)
- [ ] CLI tool produces ThaiFA-compatible output
- [ ] All endpoints have mock-based tests
- [ ] Zero external dependencies (stdlib only)
- [ ] ThaiFA can replace manual holiday files with API calls
- [ ] README with usage examples

---

## 10. Next Steps

1. **Register for BOT API key** at https://portal.api.bot.or.th/
2. **Read API docs** for holidays endpoint specifically
3. **Create repo** `github.com/jwitmann/bankofthailand-go`
4. **Initialize Go module** and project structure
5. **Build MVP** (holidays endpoint only)

---

**Notes:**
- This plan assumes the BOT holidays API returns data in a similar format to the current manual JSON files
- If the API format differs significantly, a translation adapter will be needed
- The project should follow the same patterns as `sec-go` for consistency
