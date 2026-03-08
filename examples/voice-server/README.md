# Example Voice AI Server

Simple voice AI server for testing MockingJay.

## Features

- Responds to greetings
- Handles appointment booking
- Provides business hours
- Simulates realistic latency (~100ms)

## Usage

1. Start the server:
```bash
go run main.go
```

2. Server runs on http://localhost:9000

3. Test with MockingJay:
```bash
cd ../../cli
./mockingjay init
# Edit mockingjay.yaml - set endpoint to http://localhost:9000/call
./mockingjay run
```

## API

### POST /call
```json
{
  "text": "Hello"
}
```

Response:
```json
{
  "text": "Hello! How can I help you today?",
  "intent": "greeting",
  "success": true
}
```

### GET /health
Returns: `OK`

## Supported Intents

- `greeting` - "hello", "hi"
- `booking_intent` - "book", "appointment"
- `confirmation` - "tomorrow", "7pm"
- `cancellation` - "cancel"
- `hours_response` - "hours", "open"
- `help` - "help"
- `unknown` - anything else
