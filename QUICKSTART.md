# Quick Start Guide

Get started with MockingJay in 5 minutes.

## Prerequisites

- Go 1.26+ installed
- Terminal access

## Step 1: Clone the Repository

```bash
git clone https://github.com/ashczar77/mockingjay.git
cd mockingjay
```

## Step 2: Build the CLI

```bash
cd cli
go build -o mockingjay
```

## Step 3: Start the Example Voice AI Server

```bash
# In a new terminal
cd examples/voice-server
go run main.go
```

Server will start on http://localhost:9000

## Step 4: Run Your First Test

```bash
# Back in the first terminal
cd examples/voice-server
../../cli/mockingjay run
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
  Avg latency: 105ms

✨ All tests passed!
```

## Step 5: View Results in Dashboard (Optional)

```bash
# Terminal 1: Start backend
cd cloud/backend
go run main.go

# Terminal 2: Start dashboard
cd cloud/frontend
npm install
npm run dev

# Terminal 3: Run tests with reporting
cd examples/voice-server
../../cli/mockingjay run --api-url http://localhost:8080
```

Open http://localhost:3000 to see the dashboard.

## Next Steps

- Edit `mockingjay.yaml` to add your own test scenarios
- Point the CLI to your own voice AI endpoint
- Explore the [documentation](../docs) for advanced features

## Troubleshooting

**Server won't start:**
- Check if port 9000 is already in use
- Try: `lsof -ti:9000 | xargs kill`

**Tests fail:**
- Make sure the example server is running
- Check the endpoint in `mockingjay.yaml` is correct

**Need help?**
- Open an issue: https://github.com/ashczar77/mockingjay/issues
