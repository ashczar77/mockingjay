# MockingJay 🐦

**Test voice AI agents with conversation intelligence.**

Open-source testing platform that catches bugs before your users do. Track conversation flows, validate intents, measure response quality, run A/B tests, and transcribe real calls.

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.26+-00ADD8?logo=go)](https://go.dev/)

## Why MockingJay?

Voice AI agents fail in production because manual testing doesn't scale. At 10,000 calls per day, you can't listen to them all.

MockingJay gives you:
- 💬 **Conversation Intelligence** - Track where users drop off, validate intent accuracy
- 🔄 **Multi-turn Analysis** - Measure context retention and dialogue coherence
- ✨ **Response Quality** - Score completeness, sentiment, and confidence automatically
- 🔬 **A/B Testing** - Compare two agent variants side-by-side
- 📞 **Real Call Testing** - Make actual phone calls via Twilio
- 🎙️ **Transcription** - Transcribe call recordings with Deepgram ASR
- 📊 **Visual Dashboard** - See metrics at a glance with color-coded cards
- 🚀 **Fast Execution** - Parallel testing with Go
- 🔧 **Developer-first** - CLI-first, YAML config, Git-friendly

## Quick Start

### 1. Clone and Build
```bash
git clone https://github.com/ashczar77/mockingjay.git
cd mockingjay/cli
go build -o mockingjay
```

### 2. Start the Example Voice Server
```bash
cd examples/voice-server
go run main.go
# Server starts on http://localhost:9000
```

### 3. Run Tests
```bash
cd examples/voice-server
../../cli/mockingjay run
```

### 4. View Dashboard (Optional)
```bash
# Terminal 1: Start backend
cd cloud/backend && go run main.go

# Terminal 2: Start frontend
cd cloud/frontend && npm install && npm run dev

# Open http://localhost:3000
```

## Commands

### `mockingjay run` — Run test scenarios
```bash
mockingjay run                          # uses mockingjay.yaml
mockingjay run -c my-config.yaml        # custom config
mockingjay run -s basic-greeting        # single scenario
mockingjay run --api-url http://localhost:8080  # report to dashboard
```

### `mockingjay ab` — A/B test two agent variants
```bash
mockingjay ab -c ab-test.yaml
```

Config with `ab_test` block:
```yaml
version: 1

ab_test:
  variant_a:
    name: "v1-baseline"
    endpoint: "http://localhost:9000/call"
  variant_b:
    name: "v2-new-model"
    endpoint: "http://localhost:9001/call"

scenarios:
  - name: "basic-greeting"
    steps:
      - say: "Hello"
        expect: "greeting"
```

Output:
```
🔬 A/B Test: v1-baseline vs v2-new-model

  Metric                    v1-baseline  v2-new-model        Delta
  ------                       --------      --------        -----
  Avg Latency (ms)              105.0         89.0  ✓       -16.0
  P95 Latency (ms)              108.0         92.0  ✓       -16.0
  Pass Rate (%)                 100.0        100.0            +0.0
  Task Completion (%)           100.0        100.0            +0.0

  🏆 Winner: v2-new-model
```

### `mockingjay call` — Make a real phone call via Twilio
```bash
export TWILIO_ACCOUNT_SID=ACxxx
export TWILIO_AUTH_TOKEN=xxx
export TWILIO_FROM_NUMBER=+15551234567

mockingjay call \
  --to +15559876543 \
  --webhook https://your-twiml-server.com/voice \
  --record
```

### `mockingjay transcribe` — Transcribe a call recording
```bash
# From local file
mockingjay transcribe --file recording.wav

# From URL (e.g. Twilio recording)
mockingjay transcribe --url https://api.twilio.com/recordings/xxx.mp3

# Set API key via env
export DEEPGRAM_API_KEY=xxx
```

## Configuration Reference

```yaml
version: 1

agent:
  endpoint: "http://localhost:9000/call"  # HTTP voice AI endpoint
  phone: "+15551234567"                   # or phone number for Twilio

scenarios:
  - name: "basic-greeting"
    description: "Test basic greeting flow"
    steps:
      - say: "Hello"
        expect: "greeting"

  - name: "appointment-booking"
    description: "Multi-turn booking flow"
    steps:
      - say: "I want to book an appointment"
        expect: "booking_intent"
      - say: "Tomorrow at 7pm"
        expect: "confirmation"

metrics:
  - latency
  - task_completion
  - intent_accuracy

thresholds:
  latency_p95: 5000    # ms
  task_completion: 85  # percent

# Optional: A/B test configuration
ab_test:
  variant_a:
    name: "baseline"
    endpoint: "http://localhost:9000/call"
  variant_b:
    name: "new-model"
    endpoint: "http://localhost:9001/call"
```

## Voice AI API Contract

Your agent must accept POST requests and return JSON:

**Request:**
```json
{ "text": "Hello" }
```

**Response:**
```json
{
  "text": "Hello! How can I help you today?",
  "intent": "greeting",
  "success": true
}
```

## Architecture

```
mockingjay/
├── cli/
│   ├── cmd/
│   │   ├── run.go          # mockingjay run
│   │   ├── ab.go           # mockingjay ab
│   │   ├── call.go         # mockingjay call (Twilio)
│   │   ├── transcribe.go   # mockingjay transcribe (Deepgram)
│   │   └── init.go         # mockingjay init
│   └── internal/
│       ├── ab/             # A/B test comparison
│       ├── audio/          # Deepgram transcription
│       ├── config/         # YAML config parsing
│       ├── confusion/      # Confusion pattern detection
│       ├── dialogue/       # Multi-turn dialogue analysis
│       ├── dropoff/        # Drop-off point detection
│       ├── flow/           # Conversation flow analysis
│       ├── quality/        # Response quality scoring
│       ├── reporter/       # Backend reporting client
│       ├── test/           # Test execution engine
│       ├── twilio/         # Twilio phone call client
│       └── voice/          # HTTP voice AI client
│
├── cloud/
│   ├── backend/            # Go API server (SQLite)
│   │   └── main.go         # REST API: results, metrics, ab-tests, transcriptions
│   └── frontend/           # Next.js dashboard
│       └── app/page.tsx    # Tabbed dashboard
│
└── examples/
    └── voice-server/       # Example voice AI server for testing
```

## Environment Variables

| Variable | Description |
|---|---|
| `TWILIO_ACCOUNT_SID` | Twilio Account SID |
| `TWILIO_AUTH_TOKEN` | Twilio Auth Token |
| `TWILIO_FROM_NUMBER` | Twilio phone number to call from |
| `DEEPGRAM_API_KEY` | Deepgram API key for transcription |
| `DB_PATH` | Backend SQLite database path (default: `./mockingjay.db`) |
| `PORT` | Backend server port (default: `8080`) |

## Development Status

**Week 3-4 Complete** ✅

- [x] CLI framework with parallel execution
- [x] YAML configuration with validation
- [x] HTTP client for voice AI testing
- [x] Conversation flow tracking
- [x] Intent accuracy validation
- [x] Multi-turn dialogue analysis
- [x] Response quality metrics
- [x] Drop-off point detection
- [x] Confusion pattern analysis
- [x] A/B testing framework
- [x] Twilio integration (real phone calls)
- [x] Audio recording & transcription (Deepgram)
- [x] Backend API (SQLite)
- [x] Visual dashboard (Next.js) with 4 tabs
- [ ] User authentication (Week 7-8)
- [ ] Stripe integration (Week 7-8)
- [ ] Production monitoring (Post-MVP)

## Contributing

MockingJay is open source and contributions are welcome!

- Report bugs via [GitHub Issues](https://github.com/ashczar77/mockingjay/issues)
- Submit PRs for features or fixes

## License

MIT - See [LICENSE](LICENSE) file for details

---

Built with ❤️ for the voice AI community
