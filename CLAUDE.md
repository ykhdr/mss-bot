# MSS-Bot

Telegram bot for monitoring Minecraft server status.

## Tech Stack

- **Language**: Go 1.23
- **Bot Framework**: telegram-bot-api/v5
- **Database**: SQLite with go-sqlite3 driver
- **Query Builder**: squirrel
- **Migrations**: goose/v3
- **Minecraft Query**: minequery/v2
- **Config Format**: KDL (using kdl-go)
- **Testing**: testify

## Project Structure

```
cmd/bot/          - Application entrypoint
configs/          - Configuration files (KDL format)
internal/
  app/            - Application initialization (not yet implemented)
  config/         - Configuration parsing
  storage/        - Storage interface and models
    models/       - Domain models (Server)
    sqlite/       - SQLite implementation with migrations
```

## Commands

The project includes a Makefile for common tasks:

```bash
# Build and run
make all            # Run deps, test, lint, and build
make build          # Compile the project
make run            # Run the compiled binary
make dev            # Build and run with local config

# Development
make deps           # Install dependencies
make test           # Run tests
make coverage       # Run tests with coverage report

# Linting
make lint           # Run golangci-lint
make lint-install   # Install golangci-lint
make lint-fix       # Fix linting issues automatically

# Docker
make docker-build   # Build Docker image
make docker-run     # Run with docker-compose

# CI pipeline
make ci             # Run full CI pipeline (deps, test, lint, build)
```

Manual commands:

```bash
# Run the bot
go run ./cmd/bot -config configs/config.kdl

# Run tests
go test ./...

# Build
go build -o mss-bot ./cmd/bot
```

## Configuration

Config file: `configs/config.kdl`

```kdl
bot {
    token "YOUR_TELEGRAM_BOT_TOKEN"
}

database {
    path "./data/mss-bot.db"
}

minecraft {
    timeout 5  // seconds
}
```

For local development, create `configs/config.local.kdl` (gitignored).

## Code Style and Quality

- Use squirrel for building SQL queries
- Storage layer uses interface pattern (`ServerStorage`)
- Graceful shutdown handling via signals
- Migrations registered programmatically with goose

### Linting

The project uses [golangci-lint](https://golangci-lint.run/) for code quality enforcement with a comprehensive set of linters:

- **Core linters**: errcheck, gosimple, govet, ineffassign, staticcheck, unused
- **Style linters**: gofmt, goimports, whitespace, misspell
- **Security linters**: gosec
- **Performance linters**: prealloc, bodyclose
- **Complexity linters**: gocyclo, funlen, gocognit

Configuration is stored in [`.golangci.yml`](fleet-file://catp5i93i9f126r5rb1q/Users/ykhdr/Library/Caches/JetBrains/Air/agents/air/add-golangci-linter-to-project-7e2f0f5a-5/mss-bot/.golangci.yml?type=file&root=%252F).

Run linting:
```bash
make lint           # Run all linters
make lint-fix       # Automatically fix issues where possible
```

### CI/CD

GitHub Actions workflow (`.github/workflows/ci.yml`) runs on every push and pull request:
1. **Test**: Run unit tests with coverage
2. **Lint**: Run golangci-lint with strict checks
3. **Build**: Compile binaries for artifact storage

The CI pipeline ensures code quality and prevents regressions.
