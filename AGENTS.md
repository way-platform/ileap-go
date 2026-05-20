# Agent Instructions

This document provides project-specific context and instructions for AI agents.

## Local Skills

- `way-go-style`: Use `.agents/skills/way-go-style/SKILL.md`
- `mise`: Use `.agents/skills/mise/SKILL.md`
- `ileap`: Use `.agents/skills/ileap/SKILL.md`

## Key Conventions

### Build

- **Tool manager**: mise — `mise install` then `mise run build`
- **Testing**: Use standard `testing` and `github.com/google/go-cmp/cmp` **only**. No frameworks (Testify, Ginkgo, etc.).
- **Linting**: Run `GolangCI-Lint` v2. Configure via project-specific `.golangci.yml`.

### Conformance Tests

- **`ileaptest` package**: Reusable conformance suite via `ileaptest.RunConformanceTests(t, cfg)`
- **Run against remote**: `ILEAP_SERVER_URL=<url> ILEAP_USERNAME=<u> ILEAP_PASSWORD=<p> go test -v ./ileaptest/...`
- **Env vars**: `ILEAP_SERVER_URL`, `ILEAP_USERNAME`, `ILEAP_PASSWORD`

### ACT Conformance Tests

- ACT binary panics on local URLs, but server logs show results
  - **Success**: `POST /auth/token 200`, multiple `GET /2/footprints 200`, multiple `GET /2/ileap/tad 200`, one `GET /2/ileap/tad 403`
- Our server: `https://demo.ileap.way.cloud` / `hello` / `pathfinder`
- Sine Foundation: `https://api.ileap.sine.dev` / `hello` / `pathfinder`
- **Debug remote**:
  ```bash
  gcloud beta run services logs read ileap-demo-server --project way-ileap-demo-prod --region europe-north1 --freshness='10m'
  ```
