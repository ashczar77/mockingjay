# Quick Start

Get started with MockingJay in 5 minutes.

---

## Step 1: Start the Example Server

MockingJay includes an example voice AI server for testing:

```bash
cd examples/voice-server
go run main.go
```

Server starts on `http://localhost:9000`

---

## Step 2: Initialize a Project

In a new terminal:

```bash
cd examples/voice-server
../../cli/mockingjay init
```

This creates `mockingjay.yaml` with example configuration.

---

## Step 3: Run Your First Test

```bash
mockingjay run
```

You should see:

```
🐦 MockingJay - Starting tests...

📋 Loading config from: mockingjay.yaml
🎯 Running 3 scenario(s)...

  [1/3] basic-greeting ✓ PASS (latency: 105ms)
  [2/3] appointment-booking ✓ PASS (latency: 105ms)
  [3/3] business-hours ✓ PASS (latency: 105ms)

📊 Results:
  Tests run: 3
  Passed: 3
  Failed: 0
  Pass rate: 100.0%

⚡ Performance:
  Avg latency: 105ms
  P95 latency: 105ms
  P99 latency: 105ms
  Task completion: 100.0%

✨ All tests passed!
```

---

## Step 4: Explore the Config

Open `mockingjay.yaml`:

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
    description: "Test appointment booking"
    steps:
      - say: "I want to book an appointment"
        expect: "booking_intent"
      - say: "Tomorrow at 7pm"
        expect: "confirmation"

  - name: "business-hours"
    description: "Ask about business hours"
    steps:
      - say: "What are your hours?"
        expect: "hours_response"

metrics:
  - latency
  - task_completion

thresholds:
  latency_p95: 5000
  task_completion: 85
```

---

## Step 5: Add Your Own Scenario

Edit `mockingjay.yaml` and add a new scenario:

```yaml
scenarios:
  # ... existing scenarios ...
  
  - name: "cancellation"
    description: "Test cancellation flow"
    steps:
      - say: "I need to cancel my appointment"
        expect: "cancellation"
```

Run tests again:

```bash
mockingjay run
```

---

## Step 6: Test Your Own Voice AI

Update the endpoint in `mockingjay.yaml`:

```yaml
agent:
  endpoint: "https://your-voice-agent.com/call"
```

Your voice AI must respond to POST requests:

**Request:**
```json
{
  "text": "Hello"
}
```

**Response:**
```json
{
  "text": "Hello! How can I help?",
  "intent": "greeting",
  "success": true
}
```

---

## Next Steps

- [Configuration](../configuration/overview.md) - Learn how to configure tests
