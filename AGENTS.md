# Agent Instructions

This document provides project-specific context and instructions for AI agents.

## Local Skills

- `way-go-style`: Use `.agents/skills/way-go-style/SKILL.md`
- `way-magefile`: Use `.agents/skills/way-magefile/SKILL.md`
- `ileap`: Use `.agents/skills/ileap/SKILL.md`

## Key Conventions

### Way Go Style

- **Testing**: Use standard `testing` and `github.com/google/go-cmp/cmp` **only**. No frameworks (Testify, Ginkgo, etc.).
- **Linting**: Run `GolangCI-Lint` v2. Configure via project-specific `.golangci.yml`.
- **Build**: Use `way-magefile` skill.

### Conformance Tests

- **`ileaptest` package**: Reusable conformance suite via `ileaptest.RunConformanceTests(t, cfg)`
- **Run against remote**: `./tools/mage conformancetest <baseURL> <username> <password>`
- **Env vars**: `ILEAP_SERVER_URL`, `ILEAP_USERNAME`, `ILEAP_PASSWORD`

### ACT Conformance Tests

- `ACTLocal`: `./tools/mage actlocal` — tests against a local server
  - ACT binary panics on local URLs, but server logs show results
  - **Success**: `POST /auth/token 200`, multiple `GET /2/footprints 200`, multiple `GET /2/ileap/tad 200`, one `GET /2/ileap/tad 403`
- `ACT`: `./tools/mage act <baseURL> <username> <password>` — tests against remote
  - Our server: `https://demo.ileap.way.cloud` / `hello` / `pathfinder`
  - Sine Foundation: `https://api.ileap.sine.dev` / `hello` / `pathfinder`
- **Debug remote**:
    ```bash
    gcloud beta run services logs read ileap-demo-server --project way-ileap-demo-prod --region europe-north1 --freshness='10m'
    ```
