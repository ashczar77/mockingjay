# Transcription

Transcribe call recordings using [Deepgram](https://deepgram.com) ASR.

## Prerequisites

A [Deepgram API key](https://console.deepgram.com/signup).

```bash
export DEEPGRAM_API_KEY=your_api_key
```

## Usage

### From a local file

```bash
mockingjay transcribe --file recording.wav
```

### From a URL

```bash
mockingjay transcribe --url https://api.twilio.com/recordings/xxx.mp3
```

The URL is downloaded to the current directory before transcription.

### Flags

| Flag | Env var | Description |
|------|---------|-------------|
| `--deepgram-key` | `DEEPGRAM_API_KEY` | Deepgram API key |
| `--file` | — | Local audio file path |
| `--url` | — | Remote audio URL |
| `--output-dir` | — | Directory to save downloaded audio (default: `.`) |

## Output

```
🐦 MockingJay - Transcribe

🎙️  Transcribing recording.wav...

📝 Transcript:
  Hello, I'd like to book an appointment for tomorrow at seven PM.

📊 Stats:
  Confidence: 97.3%
  Duration:   4.2s
  Words:      14
```

## Supported Formats

Deepgram supports MP3, WAV, FLAC, OGG, and most common audio formats.
