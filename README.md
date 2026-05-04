# MockingJay 🐦

**A testing harness for voice AI agents.**

MockingJay is not a voice AI agent — it's the tool you use to test one. It sends scripted conversations to your agent, validates responses, records real calls, transcribes audio, and surfaces metrics in a dashboard. Think of it as CI/CD for voice AI.

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.26+-00ADD8?logo=go)](https://go.dev/)

---

## What MockingJay Does

Your voice AI agent lives separately — it could be a GPT-powered phone bot, an IVR system, a Twilio-based assistant, or any HTTP endpoint that accepts text and returns a response. MockingJay plugs into it and tests it:

| Command | What it does |
|---|---|
| `mockingjay run` | Sends scripted multi-turn conversations to your agent's HTTP endpoint. Validates that responses match expected intents, measures latency, tracks drop-off points, scores response quality. |
| `mockingjay ab` | Runs the same scenarios against two agent variants side-by-side and declares a winner based on latency, pass rate, and task completion. |
| `mockingjay call` | Makes a real outbound phone call via Twilio to your deployed agent. Records the call. |
| `mockingjay transcribe` | Converts a call recording (local file or URL) to text using Deepgram ASR. Saves the transcript to the dashboard. |
| Dashboard | Aggregates all results — pass rates, latency, intent accuracy, A/B comparisons, transcriptions — in a visual UI. |

---

## How It Fits Into Your Workflow

```
1. Develop your voice AI agent
2. Write test scenarios in mockingjay.yaml
3. Run mockingjay run on every deploy (in CI or locally)
4. Use mockingjay call + transcribe to test the real phone experience
5. Use mockingjay ab to compare agent versions before promoting
6. Monitor quality trends in the dashboard over time
```

The `examples/voice-server` in this repo is a minimal stand-in agent for local testing. Replace it with your real agent's endpoint when testing in production.

### Testing call quality end-to-end

`mockingjay call` only tells you whether the call connected and completed — not whether your agent said the right thing. To validate actual call quality, combine all three commands:

```bash
# Step 1: Make the call and record it
mockingjay call --to +15559876543 --webhook https://your-agent.com/voice --record

# Step 2: Transcribe the recording to get the text of what was said
mockingjay transcribe --file recording.wav --api-url http://localhost:8080

# Step 3: Review the transcript in the dashboard (http://localhost:3000)
# or feed the transcript back into your scenarios to validate intent accuracy
mockingjay run --api-url http://localhost:8080
```

The transcript gives you the actual content of the call. You can review it manually in the dashboard's Transcriptions tab, or use it to update your `mockingjay.yaml` scenarios and re-run `mockingjay run` to validate that your agent's responses matched expectations.

> **Roadmap:** A future version will close this loop automatically — call → transcribe → validate → report — in a single command.

---

## Quick Start

### 1. Clone and Build
```bash
git clone https://github.com/ashczar77/mockingjay.git
cd mockingjay/cli
go build -o mockingjay
```

### 2. Start the Example Voice Server (stand-in agent for local testing)
```bash
cd examples/voice-server
go run main.go
# Starts on http://localhost:9000
```

### 3. Run Test Scenarios
```bash
cd examples/voice-server
../../cli/mockingjay run
```

### 4. View Dashboard (Optional)
```bash
# Terminal 1: backend
cd cloud/backend && go run main.go

# Terminal 2: frontend
cd cloud/frontend && npm install && npm run dev

# Open http://localhost:3000
```

---

## Commands

### `mockingjay run` — Validate agent responses

Sends each scenario's conversation steps to your agent's HTTP endpoint and checks that the returned intent matches what you expect.

```bash
mockingjay run                          # uses mockingjay.yaml
mockingjay run -c my-config.yaml        # custom config
mockingjay run -s basic-greeting        # single scenario
mockingjay run --api-url http://localhost:8080  # report results to dashboard
```

What it measures:
- **Intent accuracy** — did the agent return the expected intent?
- **Latency** — how long did each response take?
- **Task completion** — did the full conversation flow succeed?
- **Drop-off points** — which step do users most often fail at?
- **Response quality** — completeness, sentiment, confidence scores
- **Multi-turn coherence** — does the agent retain context across turns?
- **Confusion patterns** — which inputs cause the agent to misfire?

### `mockingjay ab` — Compare two agent variants

```bash
mockingjay ab -c ab-test.yaml
```

Config:
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

### `mockingjay call` — Test the real phone experience

Makes an outbound call via Twilio to your deployed agent. The `--webhook` URL is your agent's TwiML endpoint — Twilio fetches it to get instructions for what to say/do during the call.

```bash
export TWILIO_ACCOUNT_SID=ACxxx
export TWILIO_AUTH_TOKEN=xxx
export TWILIO_FROM_NUMBER=+15551234567

mockingjay call \
  --to +15559876543 \
  --webhook https://your-agent.com/voice \
  --record
```

> **Note:** Twilio trial accounts can only call verified numbers and prepend a trial message. Upgrade your Twilio account to remove these restrictions.

### `mockingjay transcribe` — Convert recordings to text

Takes a call recording and returns a transcript with confidence score and duration. Use `--api-url` to save it to the dashboard.

```bash
export DEEPGRAM_API_KEY=xxx

# From local file
mockingjay transcribe --file recording.wav --api-url http://localhost:8080

# From URL (e.g. Twilio recording)
mockingjay transcribe --url https://api.twilio.com/recordings/xxx.mp3
```

---

## Configuration Reference

```yaml
version: 1

agent:
  endpoint: "http://localhost:9000/call"  # your agent's HTTP endpoint
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
```

---

## Voice AI API Contract

For `mockingjay run` and `mockingjay ab`, your agent must accept POST requests and return JSON:

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

For `mockingjay call`, your agent must expose a TwiML webhook endpoint that Twilio can fetch during the call.

---

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
│   └── frontend/           # Next.js dashboard
│
└── examples/
    └── voice-server/       # Stand-in agent for local testing only
```

---

## Environment Variables

| Variable | Description |
|---|---|
| `TWILIO_ACCOUNT_SID` | Twilio Account SID |
| `TWILIO_AUTH_TOKEN` | Twilio Auth Token |
| `TWILIO_FROM_NUMBER` | Twilio phone number to call from |
| `DEEPGRAM_API_KEY` | Deepgram API key for transcription |
| `DB_PATH` | Backend SQLite database path (default: `./mockingjay.db`) |
| `PORT` | Backend server port (default: `8080`) |

---

## Development Status

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
- [ ] User authentication
- [ ] Stripe integration
- [ ] Production monitoring
- [ ] Automated call quality loop (call → transcribe → validate → report in one command)

---

## Contributing

MockingJay is open source and contributions are welcome!

- Report bugs via [GitHub Issues](https://github.com/ashczar77/mockingjay/issues)
- Submit PRs for features or fixes

## License

MIT - See [LICENSE](LICENSE) file for details

---

Built with ❤️ for the voice AI community
