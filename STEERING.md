# MockingJay - Steering Document

**Last Updated:** March 7, 2026  
**Status:** Planning & Initial Development

## Project Vision

Build the "Postman for Voice AI" - an open-source developer tool that automatically tests voice AI agents for quality, reliability, and performance. Differentiate through developer-first approach, open source core, and lower pricing than enterprise competitors.

## Core Principles

1. **Open source first** - Build community moat, paid cloud platform for convenience
2. **Developer experience above all** - CLI-first, code-first, not dashboard-first
3. **Ship incrementally** - 10-week MVP, iterate based on community feedback
4. **Integration-first** - Make it easy to connect with existing voice AI platforms
5. **Quality over features** - Do core testing exceptionally well before expanding

## Key Decisions Log

### March 7, 2026

**Project Pivot:**
- ❌ Rejected: Competitor price tracker (too crowded, low differentiation)
- ❌ Rejected: Industry-specific AI templates (slower moat building)
- ✅ **Selected: Voice AI Testing Platform (MockingJay)**

**Rationale:** Voice AI testing has stronger fundamentals:
- Developer tools have proven exit potential ($7.5B GitHub, $5.6B Postman)
- Only 3-4 serious competitors vs dozens in other spaces
- Open source = defensible community moat
- Critical infrastructure = high switching costs
- Clear path to $100M+ company

**Brand Decision:**
- ✅ Name: **MockingJay**
- **Rationale:** 
  - "Mock" = testing term developers understand
  - "Jay" = bird that mimics (voice/sound theme)
  - Memorable without being gimmicky
  - Strong brand potential
- ⚠️ Trademark: Lionsgate owns "Mockingjay" for entertainment (Class 41), we're software (Class 9/42) - different industries, low risk

**Tech Stack Selection:**
- ✅ Backend: Go (perfect for concurrency, performance, CLI tools)
- ✅ Database: PostgreSQL + TimescaleDB (time-series metrics)
- ✅ Queue: Redis/NATS (job scheduling)
- ✅ CLI: Cobra framework
- ✅ Frontend: Next.js with TypeScript (cloud platform)
- ✅ Payments: Stripe (usage-based billing)
- ✅ Hosting: Railway/Fly.io (backend), Vercel (frontend)

**Rationale:** Go is ideal for this use case - excellent concurrency for parallel test execution, strong CLI ecosystem, and performance for voice processing. Open source core with paid cloud platform follows proven Playwright model.

**Project Structure:**
```
mockingjay/
├── cli/              # Open source CLI tool
│   ├── cmd/         # CLI commands
│   ├── internal/
│   │   ├── test/    # Test orchestration
│   │   ├── voice/   # Voice simulation
│   │   ├── metrics/ # Evaluation engine
│   │   └── config/  # Configuration
│   └── go.mod
├── cloud/           # Paid cloud platform
│   ├── backend/     # Go API
│   └── frontend/    # Next.js dashboard
└── docs/            # Documentation
```

## Open Questions

1. **Voice simulation approach:** Use existing TTS APIs or build custom voice generation?
2. **Metric calculation:** LLM-as-judge for all metrics or rule-based for some?
3. **Open source license:** MIT (permissive) vs Apache 2.0 (patent protection)?
4. **Cloud platform pricing:** Pure usage-based or tiered with usage limits?
5. **First integration:** Start with Twilio or build generic WebRTC interface?
6. **Community platform:** Discord, Slack, or GitHub Discussions?

## Success Criteria

### Week 2 (March 21)
- [ ] CLI can run basic voice test
- [ ] Latency measurement working
- [ ] Test scenario YAML parser functional

### Week 4 (April 4)
- [ ] Full synthetic call simulation
- [ ] Core metrics: latency, WER, task completion
- [ ] Can test real voice AI agent end-to-end

### Week 8 (May 2)
- [ ] Cloud platform live with authentication
- [ ] First 5 beta users testing
- [ ] Real-time metrics dashboard

### Week 10 (May 16)
- [ ] Open source launch (GitHub)
- [ ] Product Hunt launch
- [ ] 100+ GitHub stars
- [ ] 10+ community members

### Week 12 (May 30)
- [ ] 10 paying customers
- [ ] $1,000 MRR
- [ ] 500+ GitHub stars
- [ ] Active community (Discord/Slack)

## Risks & Mitigation Strategies

| Risk | Impact | Mitigation | Status |
|------|--------|------------|--------|
| Voice platform integration too complex | High | Start with one platform (Twilio), add more later | Planned |
| Synthetic call costs (TTS/ASR APIs) | Medium | Optimize usage, cache results, usage-based pricing | Planned |
| Competition from Hamming.ai | High | Differentiate with open source + developer-first + lower pricing | Active |
| Technical complexity delays launch | High | Start with MVP metrics only, add advanced features iteratively | Active |
| Low developer adoption | High | Build in public, engage communities early, focus on DX | Planned |
| Trademark conflict with Lionsgate | Low | Different industry (software vs entertainment), monitor for issues | Accepted |
| Open source cannibalizes paid product | Medium | Cloud platform offers convenience, scale, team features | Planned |

## Scope Management

### In Scope (MVP - Week 1-10)
- CLI tool for running voice tests
- Synthetic call generation (basic scenarios)
- Core metrics: latency, WER, task completion
- Test scenario configuration (YAML)
- Basic reporting (pass/fail, metrics)
- Integration with Twilio
- Cloud platform with authentication
- Basic dashboard (test results, metrics)
- Stripe integration (usage-based billing)

### Out of Scope (Post-MVP)
- Production monitoring (continuous)
- Advanced metrics (barge-in, sentiment, etc.)
- Multiple voice platform integrations
- Team collaboration features
- CI/CD native integration
- Real-time debugging UI
- Vertical-specific test suites
- Mobile app
- Enterprise features (SSO, audit logs)

### Scope Change Process
1. Document the proposed change in this file
2. Evaluate impact on timeline and MVP goals
3. Decide: ship with MVP, add to v2, or reject
4. Update plan.md if approved

## Communication & Collaboration

**Decision-making:**
- All technical decisions validated before implementation
- Document rationale for major choices in this file
- Prefer simple solutions over clever ones

**Progress tracking:**
- Update this doc weekly with status
- Mark completed items in plan.md
- Flag blockers immediately

## Next Immediate Actions

1. [ ] Set up GitHub repository (public for open source)
2. [ ] Initialize Go module structure
3. [ ] Set up Cobra CLI framework
4. [ ] Create first test scenario YAML schema
5. [ ] Implement basic latency measurement
6. [ ] Document architecture decisions
7. [ ] Set up project board (GitHub Projects)

---

## Notes & Learnings

*This section will be updated as we build and learn*

### March 7, 2026
- **Project pivot:** Moved from competitor price tracker to voice AI testing after market research
- **Market validation:** Voice AI testing has only 3-4 serious competitors, strong growth (23.6% CAGR)
- **Differentiation strategy:** Open source core + developer-first approach vs enterprise-focused competitors
- **Brand decision:** MockingJay chosen for double meaning (mock testing + bird that mimics)
- **Tech stack:** Go selected for concurrency, performance, and CLI ecosystem
- **Business model:** Following Playwright model (open source + paid cloud platform)
- **Target:** $100M+ company potential based on developer tools market comps
