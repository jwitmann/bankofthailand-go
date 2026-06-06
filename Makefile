.PHONY: build test lint fmt clean

build:
	go build ./...

test:
	go test ./...

lint:
	golangci-lint run ./...
	staticcheck ./...
	govulncheck ./...

fmt:
	gofumpt -w .

clean:
	rm -f bot-holidays
