# MockingJay 🐦

**Mock every conversation, catch every flaw.**

Open-source testing platform for voice AI agents. Catch bugs before your users do.

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.26+-00ADD8?logo=go)](https://go.dev/)

## Why MockingJay?

Voice AI agents fail in production because manual testing doesn't scale. At 10,000 calls per day, you can't listen to them all.

MockingJay automates voice AI testing with:
- 🎯 **Synthetic call generation** - Simulate thousands of test scenarios
- 📊 **Core metrics** - Latency, task completion, word error rate
- 🔧 **Developer-first** - CLI-first, code-first approach
- 🚀 **Fast execution** - Parallel testing with Go

## Status

🚧 **Early Development** - Week 1 of 10

Currently building:
- ✅ CLI framework
- ✅ Project initialization
- ✅ Test configuration structure
- 🚧 Voice simulation
- 🚧 Metrics calculation

## Quick Start

```bash
# Install (coming soon)
go install github.com/ashczar77/mockingjay@latest

# Initialize a new project
mockingjay init

# Edit mockingjay.yaml with your agent details

# Run tests
mockingjay run
```

## Example Configuration

```yaml
version: 1

agent:
  endpoint: "https://your-voice-agent.com/call"

scenarios:
  - name: "basic-greeting"
    description: "Test basic greeting flow"
    steps:
      - say: "Hello"
        expect: "greeting"

metrics:
  - latency
  - task_completion
  - word_error_rate

thresholds:
  latency_p95: 5000  # milliseconds
  task_completion: 85  # percentage
  word_error_rate: 10  # percentage
```

## Roadmap

- **Week 1-2:** CLI + Config parsing
- **Week 3-4:** Voice simulation
- **Week 5-6:** Metrics calculation
- **Week 7-8:** Cloud platform
- **Week 9-10:** Public launch

See [PROJECT-PLAN.md](PROJECT-PLAN.md) for details.

## Architecture

```
cli/
├── cmd/              # CLI commands (init, run, version)
├── internal/
│   ├── test/        # Test orchestration
│   ├── voice/       # Voice simulation
│   ├── metrics/     # Evaluation engine
│   └── config/      # Configuration
└── main.go
```

## Contributing

MockingJay is open source and contributions are welcome!

- Report bugs via [GitHub Issues](https://github.com/ashczar77/mockingjay/issues)
- Submit PRs for features or fixes
- Join discussions in [Issues](https://github.com/ashczar77/mockingjay/issues)

## License

MIT - See [LICENSE](LICENSE) file for details

## Links

- **Documentation:** (coming soon)
- **Cloud Platform:** (coming soon)
- **Discord:** (coming soon)

---

Built with ❤️ for the voice AI community
