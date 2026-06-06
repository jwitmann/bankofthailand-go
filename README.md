# bankofthailand-go

Standalone Go client for Bank of Thailand (BOT) public APIs.

## Installation

```bash
go get github.com/jwitmann/bankofthailand-go
```

## Usage

```go
import "github.com/jwitmann/bankofthailand-go"
```

## Configuration

Create `config/bot-keys.json`:

```json
{
  "app_id": "your-app-id",
  "token": "your-static-token"
}
```

Or set environment variables:

```bash
export BOT_API_APP_ID=your-app-id
export BOT_API_TOKEN=your-static-token
```

## CLI

```bash
go run ./cmd/bot-holidays -year 2026
```

## License

MIT
