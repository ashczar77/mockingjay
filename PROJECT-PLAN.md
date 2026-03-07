# MockingJay - Voice AI Testing Platform

**Project Start:** March 7, 2026  
**Target Launch:** May 16, 2026 (10 weeks)  
**Goal:** Launch open-source core + paid cloud platform

## Product Overview

**What:** Developer tool that automatically tests voice AI agents for quality, reliability, and performance  
**Who:** Developers, startups, AI companies building voice agents  
**Pain:** Manual testing doesn't scale, voice agents have 13+ failure modes, production issues cost $$$  
**Value:** Catch bugs before users do, ensure quality at scale, prevent costly production failures  
**Pricing:** Open source (free) + Cloud ($99-$5,000/month)

## Brand

**Name:** MockingJay
- "Mock" = testing term (creating mock/simulated calls)
- "Jay" = bird that mimics sounds
- Double meaning: technical + thematic

**Tagline:** "Mock every conversation, catch every flaw"

**Logo:** Stylized mockingjay bird (abstract, modern)

## Tech Stack ✅ APPROVED

**Backend:**
- Go (perfect for concurrency and performance)
- PostgreSQL (structured data) + TimescaleDB (time-series metrics)
- Redis/NATS (job queue for test execution)
- FFmpeg + WebRTC libraries (voice processing)

**Test Orchestration:**
- Go goroutines for parallel test execution
- YAML/JSON for test scenario configs
- Synthetic voice generation (TTS APIs)

**CLI:**
- Cobra (Go CLI framework)
- Interactive prompts

**Frontend (Cloud Platform):**
- Next.js (TypeScript)
- React components
- Tailwind CSS
- Real-time dashboards (WebSocket)

**Integrations:**
- Twilio (telephony)
- OpenAI/Anthropic (LLM evaluation)
- Deepgram/AssemblyAI (ASR for WER calculation)
- ElevenLabs/Play.ht (TTS for synthetic calls)

**Payments:**
- Stripe Checkout
- Usage-based billing

**Hosting:**
- Railway or Fly.io (Go backend)
- Vercel (Next.js frontend)
- AWS S3 (call recordings storage)

## Build Timeline

### Week 1-2 (March 7-20): Open Source Core
- [ ] Set up Go project structure
- [ ] Build CLI framework (Cobra)
- [ ] Implement test scenario parser (YAML)
- [ ] Basic test orchestration engine
- [ ] Simple synthetic call simulation
- [ ] Core metrics: latency tracking
- [ ] Pass/fail reporting

### Week 3-4 (March 21 - April 3): Voice Simulation
- [ ] Integrate TTS for synthetic voice generation
- [ ] Integrate ASR for transcription
- [ ] Implement WER calculation
- [ ] Add background noise simulation
- [ ] Multi-turn conversation support
- [ ] Test with real voice AI platform (Twilio)

### Week 5-6 (April 4-17): Advanced Metrics
- [ ] Task completion detection (LLM-as-judge)
- [ ] Barge-in/interruption testing
- [ ] Latency percentiles (P50, P95, P99)
- [ ] Tool call success tracking
- [ ] Sentiment analysis
- [ ] Detailed failure categorization

### Week 7-8 (April 18 - May 1): Cloud Platform
- [ ] Next.js dashboard setup
- [ ] User authentication
- [ ] Test result visualization
- [ ] Real-time metrics dashboards
- [ ] Team management
- [ ] Stripe integration

### Week 9-10 (May 2-15): Polish & Launch
- [ ] CI/CD integration (GitHub Actions)
- [ ] Documentation (README, guides)
- [ ] Example test scenarios
- [ ] Landing page
- [ ] Open source release (GitHub)
- [ ] Product Hunt launch

### Week 11+ (May 16+): Growth
- [ ] Community building (Discord/Slack)
- [ ] Integration partnerships (Twilio, ElevenLabs)
- [ ] Content marketing (blog posts, tutorials)
- [ ] First 10 paying customers
- [ ] Iterate based on feedback

## MVP Feature Set

### Core Features (Must Have)
- [ ] CLI tool for running voice tests
- [ ] Synthetic call generation (basic scenarios)
- [ ] Core metrics: latency, task completion, WER
- [ ] Test scenario configuration (YAML/JSON)
- [ ] Basic reporting (pass/fail, metrics)
- [ ] Integration with one voice platform (e.g., Twilio)

### Nice to Have (v2)
- Production monitoring
- Advanced metrics (barge-in, sentiment, etc.)
- Multiple voice platform integrations
- Web dashboard
- Team collaboration
- CI/CD integration

## Success Metrics

**Week 2:**
- CLI can run basic voice test
- Latency measurement working

**Week 4:**
- Full synthetic call simulation
- Core metrics calculated (latency, WER, task completion)

**Week 6:**
- Advanced metrics implemented
- Can test real voice AI agent end-to-end

**Week 8:**
- Cloud platform live
- First beta users testing

**Week 10:**
- Open source launch
- 100+ GitHub stars
- 10+ community members

**Week 12:**
- 10 paying customers
- $1,000 MRR
- 500+ GitHub stars

## Risks & Mitigations

| Risk | Mitigation |
|------|------------|
| Voice platform integration complexity | Start with one platform (Twilio), add more later |
| Synthetic call costs (TTS/ASR APIs) | Optimize usage, cache results, offer usage-based pricing |
| Competition from Hamming.ai | Differentiate with open source + developer-first approach |
| Technical complexity too high | Start with MVP metrics, add advanced features iteratively |
| Low adoption | Build in public, engage developer community early |

## Competitive Analysis

**Existing solutions:**
- Hamming.ai (enterprise-focused, likely $500-$5K+/month)
- Cekura.ai (newer, less mature)
- Maxim AI (human review workflows)
- Generic LLM eval tools (Braintrust, Langfuse - text only, not voice-native)

**Our advantages:**
- Open source core (community moat)
- Developer-first (CLI, not dashboard-first)
- Lower price point ($99-$499 vs $500-$5K)
- Real-time debugging experience
- Integration-first approach

## Marketing Strategy

**Week 1-4:**
- Build in public on Twitter/X
- Share progress on Indie Hackers
- Engage in voice AI communities (Discord, Reddit)

**Week 5-8:**
- Technical blog posts (How to test voice AI)
- Open source contributions to related projects
- Reach out to voice AI companies for beta testing

**Week 9-10:**
- Product Hunt launch
- Hacker News "Show HN"
- Dev.to articles

**Month 3+:**
- Integration partnerships (Twilio, ElevenLabs)
- Conference talks (AI/DevOps conferences)
- YouTube tutorials
- Developer advocacy program

## Next Immediate Actions

1. [ ] Set up Go project repository
2. [ ] Initialize Go modules and project structure
3. [ ] Build basic CLI with Cobra
4. [ ] Create first test scenario (YAML format)
5. [ ] Implement simple latency measurement

---

**Decision Point:** After Week 2, evaluate if core testing works. If voice simulation is too complex, start with simpler HTTP-based testing and add voice later.
