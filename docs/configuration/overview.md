# Configuration Overview

Learn how to configure MockingJay test scenarios.

---

## Config File Location

By default, MockingJay looks for `mockingjay.yaml` in the current directory.

Use a custom location:

```bash
mockingjay run --config /path/to/config.yaml
```

---

## Basic Structure

```yaml
version: 1                    # Config version (required)

agent:                        # Voice AI agent to test
  endpoint: "https://..."     # HTTP endpoint

scenarios:                    # Test scenarios
  - name: "test-name"
    steps:
      - say: "input"
        expect: "intent"

metrics:                      # Metrics to track
  - latency
  - task_completion

thresholds:                   # Performance thresholds
  latency_p95: 5000
  task_completion: 85
```

---

## Agent Configuration

Specify how to connect to your voice AI:

### HTTP Endpoint

```yaml
agent:
  endpoint: "https://api.example.com/call"
```

Your endpoint must accept POST requests with:

```json
{
  "text": "user input"
}
```

And respond with:

```json
{
  "text": "agent response",
  "intent": "detected_intent",
  "success": true
}
```

### Phone Number (Coming Soon)

```yaml
agent:
  phone: "+1-555-0100"
```

---

## Next Steps

- [Scenarios](scenarios.md) - Define test scenarios
- [Metrics](metrics.md) - Configure metrics and thresholds
