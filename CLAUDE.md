# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Service Startup

When the service starts, it will send a test message to confirm the Telegram Bot is working properly. This test message will include:
- Service startup confirmation
- Basic configuration status
- Version number in format: vYYYY.MM.DD.HHMM (e.g., v2025.07.09.1455)

The version number is read from `version.txt` file and should be manually updated each time code is modified. The version format is `YYYY.MM.DD.HHMM` based on the date and time when the code change was made.

## Project Overview

This is a Go-based Taiwan Railway monitoring service that tracks train schedules for Zhubei Station and sends notifications via Telegram Bot. The service monitors trains during configured hours (default 18:00-23:00) and sends updates every 30 minutes.

## Architecture

The project follows a clean architecture pattern with separate packages:

- `main.go` - Application entry point with signal handling
- `internal/config/` - Configuration management with environment variables
- `internal/tdx/` - TDX (Taiwan Data Exchange) API client for railway data
- `internal/telegram/` - Telegram Bot API client for notifications
- `internal/monitor/` - Scheduling and monitoring logic using cron jobs

## Key Components

### Configuration (`internal/config/config.go`)
- Loads environment variables from `.env` file
- Supports both authenticated and free-tier TDX API usage
- Required: `TELEGRAM_BOT_TOKEN`, `TELEGRAM_CHAT_ID`, `ZHUBEI_STATION_ID`
- Optional: `TDX_CLIENT_ID`, `TDX_CLIENT_SECRET` for higher API limits

### TDX Client (`internal/tdx/client.go`)
- Handles Taiwan Railway API authentication and requests
- Supports both authenticated and free-tier access
- Fetches train timetables and route information

### Telegram Bot (`internal/telegram/bot.go`)
- Sends formatted train information to configured chat
- Supports HTML formatting for messages
- Handles both basic and detailed train information

### Monitor Scheduler (`internal/monitor/scheduler.go`)
- Uses cron jobs for periodic checking (default every 30 minutes)
- Only monitors during configured hours (18:00-23:00)
- Filters upcoming trains and sends notifications

## Common Development Commands

```bash
# Install dependencies
go mod tidy

# Run the application
go run main.go

# Build the application
go build -o main .

# Run with Docker
docker build -t tg-rail-bot .
docker run -d --name tg-rail-bot -v $(pwd)/.env:/root/.env tg-rail-bot
```

## Environment Setup

Copy `.env.example` to `.env` and configure:
- `TELEGRAM_BOT_TOKEN` - Get from @BotFather
- `TELEGRAM_CHAT_ID` - Get from bot API `/getUpdates`
- `ZHUBEI_STATION_ID` - Station ID (default: 1180)
- `TDX_CLIENT_ID` / `TDX_CLIENT_SECRET` - Optional for higher API limits

## Testing Files

- `test_*.go` - Various API testing utilities
- `debug_api.go` - Debug utilities for API exploration
- `explore_api.go` - API exploration tools

## Dependencies

- `github.com/go-resty/resty/v2` - HTTP client
- `github.com/joho/godotenv` - Environment variable loading
- `github.com/robfig/cron/v3` - Cron job scheduling
- `github.com/sirupsen/logrus` - Structured logging

## Service Behavior

- Runs initial API test on startup
- Monitors only during configured hours (18:00-23:00)
- Sends notifications for upcoming trains
- Filters trains by direction (1=northbound, 0=southbound)
- Provides detailed route information to Fugang Station
- Graceful shutdown on SIGINT/SIGTERM

## Version Management

**IMPORTANT**: Every time you modify the code, you MUST update the version number in `version.txt` file:

1. Version format: `YYYY.MM.DD.HHMM` (e.g., `2025.07.10.0128`)
2. **MUST use Taiwan time (Asia/Taipei timezone)** - Use `TZ=Asia/Taipei date +"%Y.%m.%d.%H%M"` command
3. Based on the date and time when the code change was made
4. The service will display this version with `v` prefix (e.g., `v2025.07.10.0128`) in the startup message
5. This ensures each code change is properly tracked and visible in the Telegram notifications

**Example workflow:**
- Make code changes
- Get Taiwan time: `TZ=Asia/Taipei date +"%Y.%m.%d.%H%M"`
- Update `version.txt` with Taiwan timestamp: `2025.07.10.0128`
- The startup message will show: `📋 版本: v2025.07.10.0128`