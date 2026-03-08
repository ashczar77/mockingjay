# MockingJay 🐦

**Mock every conversation, catch every flaw.**

Open-source testing platform for voice AI agents. Catch bugs before your users do.

---

## Why MockingJay?

Voice AI agents fail in production because manual testing doesn't scale. At 10,000 calls per day, you can't listen to them all.

MockingJay automates voice AI testing with:

- 🎯 **Synthetic call generation** - Simulate thousands of test scenarios
- 📊 **Core metrics** - Latency, task completion, word error rate
- 🔧 **Developer-first** - CLI-first, code-first approach
- 🚀 **Fast execution** - Parallel testing with Go

---

## Quick Example

```yaml
# mockingjay.yaml
version: 1

agent:
  endpoint: "https://your-voice-agent.com/call"

scenarios:
  - name: "greeting"
    steps:
      - say: "Hello"
        expect: "greeting"
```

```bash
mockingjay run
```

```
🐦 MockingJay - Starting tests...

  [1/1] greeting ✓ PASS (latency: 150ms)

📊 Results:
  Tests run: 1
  Passed: 1
  Pass rate: 100.0%

⚡ Performance:
  Avg latency: 150ms
  P95 latency: 150ms
  Task completion: 100.0%

✨ All tests passed!
```

---

## Features

### 🎯 Scenario Testing
Define realistic conversation flows and test them automatically.

### 📊 Performance Metrics
Track latency (avg, P95, P99), task completion, and error rates.

### 🔄 Parallel Execution
Run multiple scenarios simultaneously for fast feedback.

### 📈 Dashboard
View test results and trends over time with the built-in dashboard.

### 🔌 API Integration
Store results in the backend and integrate with your CI/CD pipeline.

---

## Get Started

Choose your path:

<div class="grid cards" markdown>

-   :material-clock-fast:{ .lg .middle } __Quick Start__

    ---

    Get up and running in 5 minutes

    [:octicons-arrow-right-24: Quick Start](getting-started/quickstart.md)

-   :material-download:{ .lg .middle } __Installation__

    ---

    Install MockingJay on your system

    [:octicons-arrow-right-24: Installation](getting-started/installation.md)

-   :material-file-document:{ .lg .middle } __Configuration__

    ---

    Learn how to configure test scenarios

    [:octicons-arrow-right-24: Configuration](configuration/overview.md)

-   :material-code-braces:{ .lg .middle } __Examples__

    ---

    See real-world examples and use cases

    [:octicons-arrow-right-24: Examples](examples.md)

</div>

---

## Open Source

MockingJay is MIT licensed and open source.

- **GitHub:** [ashczar77/mockingjay](https://github.com/ashczar77/mockingjay)
- **Issues:** [Report bugs](https://github.com/ashczar77/mockingjay/issues)
- **Discussions:** [Ask questions](https://github.com/ashczar77/mockingjay/discussions)
