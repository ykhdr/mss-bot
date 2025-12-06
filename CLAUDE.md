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

## Code Style

- Use squirrel for building SQL queries
- Storage layer uses interface pattern (`ServerStorage`)
- Graceful shutdown handling via signals
- Migrations registered programmatically with goose
