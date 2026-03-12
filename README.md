# MockingJay 🐦

**Test voice AI agents with conversation intelligence.**

Open-source testing platform that catches bugs before your users do. Track conversation flows, validate intents, and measure response quality automatically.

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.26+-00ADD8?logo=go)](https://go.dev/)

## Why MockingJay?

Voice AI agents fail in production because manual testing doesn't scale. At 10,000 calls per day, you can't listen to them all.

MockingJay gives you:
- 💬 **Conversation Intelligence** - Track where users drop off, validate intent accuracy
- 🔄 **Multi-turn Analysis** - Measure context retention and dialogue coherence
- ✨ **Response Quality** - Score completeness, sentiment, and confidence automatically
- 📊 **Visual Dashboard** - See metrics at a glance with color-coded cards
- 🚀 **Fast Execution** - Parallel testing with Go
- 🔧 **Developer-first** - CLI-first, YAML config, Git-friendly

## Features

### ✅ Conversation Intelligence
```bash
💬 Conversation Intelligence:
  Success rate: 100.0%
  Intent accuracy: 100.0% (4/4 correct)
  Avg steps completed: 1.3
  Avg conversation duration: 142ms

🔄 Multi-turn Dialogue:
  Multi-turn conversations: 1/3
  Context retention: 100.0%
  Coherence score: 100.0%

✨ Response Quality:
  Completeness: 100.0%
  Positive sentiment: 75.0%
  Confidence: 100.0%
  Avg response length: 67 chars
```

### ✅ Performance Metrics
- P95/P99 latency tracking
- Task completion rates
- Parallel test execution
- Drop-off point detection

### ✅ Visual Dashboard
- Real-time metrics display
- Color-coded quality indicators
- Conversation flow visualization
- Historical test results

### 🔜 Coming Soon
- A/B testing framework (Week 3)
- Confusion pattern analysis (Week 3)
- Real phone call testing via Twilio (Week 4)
- Audio recording and transcription (Week 4)

## Quick Start

### 1. Clone and Build
```bash
git clone https://github.com/ashczar77/mockingjay.git
cd mockingjay/cli
go build -o mockingjay
```

### 2. Create Test Configuration
```bash
./mockingjay init
```

This creates `mockingjay.yaml`:
```yaml
version: 1

agent:
  endpoint: "http://localhost:9000/call"

scenarios:
  - name: "basic-greeting"
    description: "Test basic greeting flow"
    steps:
      - say: "Hello"
        expect: "greeting"
  
  - name: "appointment-booking"
    description: "Test multi-turn booking"
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
  latency_p95: 5000
  task_completion: 85
```

### 3. Run Tests
```bash
./mockingjay run
```

### 4. View Dashboard (Optional)
```bash
# Terminal 1: Start backend
cd cloud/backend
go run main.go

# Terminal 2: Start frontend
cd cloud/frontend
npm install
npm run dev

# Open http://localhost:3000
```

## Architecture

```
mockingjay/
├── cli/                          # Command-line interface
│   ├── cmd/                      # Commands (init, run, version)
│   ├── internal/
│   │   ├── config/              # YAML parsing
│   │   ├── test/                # Test execution
│   │   ├── flow/                # Conversation flow analysis
│   │   ├── dialogue/            # Multi-turn dialogue tracking
│   │   ├── quality/             # Response quality scoring
│   │   ├── voice/               # HTTP client for voice AI
│   │   └── reporter/            # Backend reporting
│   └── main.go
│
├── cloud/
│   ├── backend/                 # API server (Go + SQLite)
│   └── frontend/                # Dashboard (Next.js)
│
└── examples/
    └── voice-server/            # Example voice AI server
```

## How It Works

1. **Define scenarios** in YAML with expected conversation flows
2. **Run tests** - CLI executes scenarios in parallel
3. **Analyze results** - Tracks intent accuracy, context retention, response quality
4. **View metrics** - CLI output or visual dashboard
5. **Iterate** - Fix issues and re-test

## Example Output

```bash
🐦 MockingJay - Starting tests...

📋 Loading config from: mockingjay.yaml
🎯 Running 3 scenario(s)...

  [1/3] basic-greeting ✓ PASS (latency: 108ms)
  [2/3] appointment-booking ✓ PASS (latency: 108ms)
  [3/3] business-hours ✓ PASS (latency: 108ms)

📊 Results:
  Tests run: 3
  Passed: 3
  Failed: 0
  Pass rate: 100.0%

⚡ Performance:
  Avg latency: 108ms
  P95 latency: 108ms
  P99 latency: 108ms

💬 Conversation Intelligence:
  Success rate: 100.0%
  Intent accuracy: 100.0% (4/4 correct)
  Avg steps completed: 1.3

🔄 Multi-turn Dialogue:
  Context retention: 100.0%
  Coherence score: 100.0%

✨ Response Quality:
  Completeness: 100.0%
  Sentiment: 75.0%
  Confidence: 100.0%

✨ All tests passed!
```

## Development Status

**Week 2 Complete** ✅ - Conversation Intelligence shipped!

- [x] CLI framework with parallel execution
- [x] YAML configuration with validation
- [x] HTTP client for voice AI testing
- [x] Conversation flow tracking
- [x] Intent accuracy validation
- [x] Multi-turn dialogue analysis
- [x] Response quality metrics
- [x] Backend API (SQLite)
- [x] Visual dashboard (Next.js)
- [x] CI/CD with GitHub Actions
- [ ] A/B testing framework (Week 3)
- [ ] Confusion pattern analysis (Week 3)
- [ ] Twilio integration (Week 4)
- [ ] Audio recording (Week 4)

## Contributing

MockingJay is open source and contributions are welcome!

- Report bugs via [GitHub Issues](https://github.com/ashczar77/mockingjay/issues)
- Submit PRs for features or fixes
- Join discussions in [Issues](https://github.com/ashczar77/mockingjay/issues)

## License

MIT - See [LICENSE](LICENSE) file for details

---

Built with ❤️ for the voice AI community
