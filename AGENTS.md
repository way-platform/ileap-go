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

### ACT Conformance Tests

To run the ACT test suite:
- `ACTLocal`: Run `./tools/mage actlocal` to test against a local server.
  - The ACT binary currently crashes on local URLs with a `panic`, but the server logs indicate if the tests succeeded prior to the crash.
  - **Success Criteria:** Inspect the output for `level=INFO msg=request` lines. A successful run will show:
    - `POST path=/auth/token status=200`
    - Multiple `GET path=/2/footprints status=200`
    - Multiple `GET path=/2/ileap/tad status=200`
    - One `GET path=/2/ileap/tad status=403` (testing invalid token, expected to be 403)
  - If any of the expected `200` requests return `4xx` (other than the deliberate `403`/`401` tests) or `5xx`, the test failed.
- `ACT`: Run `./tools/mage act <baseURL> <username> <password>` to test against a specific remote server.
  - The Sine Foundation's reference implementation uses:
    - Base URL: `https://api.ileap.sine.dev`
    - User: `hello`
    - Password: `pathfinder`
  - Our implementation uses:
    - Base URL: `https://ileap.wayplatform.com`
    - User: `ileap-demo@way.cloud`
    - Password: `HelloPrimaryData`
  - **Debugging Remote Tests:** To investigate failures on our remote demo server, view the Cloud Run logs:
    ```bash
    gcloud beta run services logs read ileap-demo-server --project way-ileap-demo-prod --region europe-north1 --freshness='10m'
    ```
