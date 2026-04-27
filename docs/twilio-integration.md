# Real Call Testing (Twilio)

Make actual outbound phone calls to test your voice AI agent end-to-end.

## Prerequisites

1. A [Twilio account](https://www.twilio.com/try-twilio)
2. A Twilio phone number
3. A TwiML webhook URL that instructs the call (e.g. connects to your voice AI)

## Setup

```bash
export TWILIO_ACCOUNT_SID=ACxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
export TWILIO_AUTH_TOKEN=your_auth_token
export TWILIO_FROM_NUMBER=+15551234567
```

## Usage

```bash
mockingjay call \
  --to +15559876543 \
  --webhook https://your-server.com/twiml \
  --record
```

### Flags

| Flag | Env var | Description |
|------|---------|-------------|
| `--account-sid` | `TWILIO_ACCOUNT_SID` | Twilio Account SID |
| `--auth-token` | `TWILIO_AUTH_TOKEN` | Twilio Auth Token |
| `--from` | `TWILIO_FROM_NUMBER` | Caller phone number |
| `--to` | — | Phone number to call (required) |
| `--webhook` | — | TwiML webhook URL (required) |
| `--record` | — | Record the call |

## Output

```
🐦 MockingJay - Real Call Test

📞 Calling +15559876543 from +15551234567...
🔴 Recording enabled

📊 Call Result:
  SID:      CAxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
  Status:   completed
  Duration: 45s
  Recording: https://api.twilio.com/2010-04-01/Accounts/.../Recordings/...

✨ Call completed successfully!
```

## Transcribing the Recording

After a call with `--record`, transcribe the recording:

```bash
mockingjay transcribe --url https://api.twilio.com/2010-04-01/Accounts/.../Recordings/xxx.mp3
```
