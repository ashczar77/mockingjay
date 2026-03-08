# MockingJay Backend API

REST API for storing and retrieving test results.

## Setup

1. Install dependencies:
```bash
go mod download
```

2. Run server:
```bash
go run main.go
```

Database file (`mockingjay.db`) is created automatically.

## API Endpoints

### POST /api/results
Store test result
```json
{
  "scenario": "basic-greeting",
  "passed": true,
  "latency_ms": 150,
  "error": ""
}
```

### GET /api/results?limit=100
Get recent test results

### GET /api/health
Health check

## Environment Variables

- `DB_PATH` - SQLite database file path (default: `./mockingjay.db`)
- `PORT` - Server port (default: `8080`)
