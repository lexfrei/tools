# ow-exporter

Overwatch 2 statistics exporter for Prometheus monitoring.

## Overview

ow-exporter scrapes Overwatch 2 player profiles and exposes hero statistics via Prometheus metrics. It supports both PC and Console platforms, Quick Play and Competitive game modes.

## Features

- **Platform Support**: Separate statistics for PC (keyboard+mouse) and Console (controller)
- **Game Mode Support**: Quick Play and Competitive statistics
- **15 Core Metrics**: Time played, win rate, eliminations, accuracy, and more
- **REST API**: Add/remove players via HTTP API
- **Persistence**: SQLite storage for user management
- **Background Updates**: Automatic profile scraping every 5 minutes
- **Rate Limiting**: Protects against Blizzard rate limits

## Prometheus Metrics

All metrics include comprehensive labels:

```
ow_hero_time_played_seconds{username="player1", hero="mercy", platform="pc", gamemode="competitive"} 15540
ow_hero_win_percentage{username="player1", hero="mercy", platform="pc", gamemode="competitive"} 67.5
ow_hero_eliminations_per_life{username="player1", hero="mercy", platform="pc", gamemode="competitive"} 2.3
```

## API Endpoints

- `GET /api/users` - List all tracked users
- `POST /api/users` - Add new user to track
- `PUT /api/users/{username}` - Update user settings
- `DELETE /api/users/{username}` - Remove user
- `GET /metrics` - Prometheus metrics endpoint
- `GET /health` - Health check

## Configuration

```yaml
database:
  type: "sqlite"
  path: "/data/ow-exporter.db"

scraper:
  interval: "5m"
  timeout: "30s"
  user_agent: "ow-exporter/1.0"

server:
  port: 9420
  enable_api: true
```

## Usage

### Running with Docker:
```bash
# Run ow-exporter container
docker run -d \
  --name ow-exporter \
  -p 9420:9420 \
  -e PORT=9420 \
  ghcr.io/lexfrei/ow-exporter:latest

# Check health
curl http://localhost:9420/health
```

### Add a user to track:
```bash
curl -X POST http://localhost:9420/api/users \
  -H "Content-Type: application/json" \
  -d '{
    "username": "player1",
    "profile_id": "de5bb4aca17492e0bba120a1d1%7Ca92a11ef8d304356fccfff8df12e1dc6"
  }'
```

### View Prometheus metrics:
```bash
curl http://localhost:9420/metrics
```

## Development Status

🚧 **In Development** - See [Issue #439](https://github.com/lexfrei/tools/issues/439) for progress.

## Architecture

- **No headless browser required** - Uses HTTP + HTML parsing
- **Platform detection** - Automatic PC/Console separation
- **Hardcoded metrics** - Stable metric definitions for all heroes
- **Background scraping** - Non-blocking profile updates