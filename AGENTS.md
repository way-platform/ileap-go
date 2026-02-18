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
