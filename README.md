# MockingJay 🐦

**Mock every conversation, catch every flaw.**

MockingJay is an open-source developer tool for testing voice AI agents. Catch bugs before your users do.

## Status

🚧 **Early Development** - Week 1 of 10-week build

Currently implementing:
- ✅ CLI framework (Cobra)
- ✅ Project initialization
- ✅ Basic test runner structure
- 🚧 Config parsing
- 🚧 Voice simulation
- 🚧 Metrics calculation

## Quick Start

```bash
# Install (coming soon)
go install github.com/yourusername/mockingjay@latest

# Initialize a new project
mockingjay init

# Edit mockingjay.yaml with your agent details

# Run tests
mockingjay run
```

## Features (Planned)

- 🎯 **Synthetic Call Generation** - Simulate thousands of test calls
- 📊 **Core Metrics** - Latency, task completion, word error rate
- 🔧 **Developer-First** - CLI-first, code-first approach
- 📝 **YAML Configuration** - Simple, readable test scenarios
- 🚀 **Fast Execution** - Parallel test execution with Go
- 🔌 **Integrations** - Twilio, ElevenLabs, Deepgram, and more

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
  latency_p95: 5000
  task_completion: 85
  word_error_rate: 10
```

## Architecture

```
cli/
├── cmd/              # CLI commands
├── internal/
│   ├── test/        # Test orchestration
│   ├── voice/       # Voice simulation
│   ├── metrics/     # Evaluation engine
│   └── config/      # Configuration
└── main.go
```

## Roadmap

**Week 1-2:** CLI + Config parsing  
**Week 3-4:** Voice simulation  
**Week 5-6:** Metrics calculation  
**Week 7-8:** Cloud platform  
**Week 9-10:** Launch  

## Contributing

MockingJay is open source! Contributions welcome.

## License

MIT (coming soon)

## Links

- Documentation: (coming soon)
- Cloud Platform: (coming soon)
- Discord: (coming soon)

---

Built with ❤️ for the voice AI community
