# Contributing to MockingJay

Thanks for your interest in contributing! MockingJay is open source and welcomes PRs.

## Development Setup

```bash
git clone https://github.com/ashczar77/mockingjay.git
cd mockingjay

# Build CLI
cd cli && go build -o mockingjay

# Run example server
cd ../examples/voice-server && go run main.go

# Run tests
cd examples/voice-server && ../../cli/mockingjay run
```

## Project Structure

```
cli/          - Go CLI tool (open source core)
cloud/        - Backend API + Next.js dashboard
examples/     - Example voice AI server and configs
docs/         - Documentation
```

## Making Changes

1. Fork the repo
2. Create a branch: `git checkout -b feature/my-feature`
3. Make your changes
4. Verify the build: `cd cli && go build ./...`
5. Run `go vet ./...` to check for issues
6. Submit a PR

## Adding a New CLI Command

1. Create `cli/cmd/mycommand.go`
2. Register it in `cli/cmd/root.go`: `rootCmd.AddCommand(myCmd)`
3. Add any internal logic to `cli/internal/`

## Reporting Bugs

Open an issue at https://github.com/ashczar77/mockingjay/issues with:
- What you expected to happen
- What actually happened
- Steps to reproduce
- Go version and OS
