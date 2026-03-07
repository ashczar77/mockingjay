# MockingJay - Project Plan

**Status:** Week 1 of 10  
**Target Launch:** May 16, 2026

## Vision

Build the "Postman for Voice AI" - an open-source testing platform that helps developers catch voice AI bugs before production.

## Product

**What:** CLI tool + cloud platform for testing voice AI agents  
**Model:** Open source core + paid cloud (Playwright model)  
**Pricing:** Free (self-hosted) | $99-$5,000/month (cloud)

## 10-Week Roadmap

### Week 1-2: CLI Foundation
- [x] CLI framework (Cobra)
- [x] Project initialization
- [x] Config structure (YAML)
- [ ] Config file parsing
- [ ] Basic test orchestration

### Week 3-4: Voice Simulation
- [ ] TTS integration (synthetic voices)
- [ ] ASR integration (transcription)
- [ ] Multi-turn conversations
- [ ] Background noise simulation
- [ ] Test with Twilio

### Week 5-6: Metrics Engine
- [ ] Latency tracking (P50, P95, P99)
- [ ] Word Error Rate (WER) calculation
- [ ] Task completion detection
- [ ] Barge-in/interruption testing
- [ ] Report generation

### Week 7-8: Cloud Platform
- [ ] Next.js dashboard
- [ ] User authentication
- [ ] Test result visualization
- [ ] Team management
- [ ] Stripe integration

### Week 9-10: Launch
- [ ] Documentation
- [ ] Example scenarios
- [ ] CI/CD integration
- [ ] Landing page
- [ ] Product Hunt launch

## Tech Stack

**CLI (Open Source):**
- Go 1.26+
- Cobra (CLI framework)
- YAML config

**Cloud Platform:**
- Backend: Go + PostgreSQL + TimescaleDB
- Frontend: Next.js + TypeScript
- Queue: Redis/NATS
- Hosting: Railway/Fly.io + Vercel

**Integrations:**
- Twilio (telephony)
- OpenAI/Anthropic (LLM evaluation)
- Deepgram/AssemblyAI (ASR)
- ElevenLabs/Play.ht (TTS)

## Success Metrics

**Week 2:** CLI can run basic tests  
**Week 4:** Voice simulation working  
**Week 8:** Cloud platform live, 5 beta users  
**Week 10:** Public launch, 100+ GitHub stars  
**Week 12:** 10 paying customers, $1K MRR

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) (coming soon)
