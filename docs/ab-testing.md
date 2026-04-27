# A/B Testing

Compare two versions of your voice AI agent side-by-side.

## Usage

```bash
mockingjay ab -c ab-test.yaml
```

## Configuration

Add an `ab_test` block to your config file. The top-level `agent` field is not required when using A/B testing.

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
  - name: "appointment-booking"
    steps:
      - say: "I want to book an appointment"
        expect: "booking_intent"
      - say: "Tomorrow at 7pm"
        expect: "confirmation"
```

## Output

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

- `✓` means the metric improved for variant B
- `✗` means the metric got worse for variant B
- The winner is determined by the number of improved metrics

## Metrics Compared

| Metric | Better when |
|--------|-------------|
| Pass Rate | Higher |
| Avg Latency | Lower |
| P95 Latency | Lower |
| Task Completion | Higher |
