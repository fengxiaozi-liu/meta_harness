# Meta Harness

Code review and testing harness for AI-assisted development.

## Overview

Meta Harness is a Go-based harness framework that integrates with code review tools to validate AI-generated code changes.

## Structure

```
meta_harness/
├── cmd/run/              # Entry point
├── internal/
│   ├── adapter/          # External tool adapters (Claude CLI, Codex CLI)
│   ├── controller/       # Main orchestration logic
│   ├── executor/         # Test execution
│   ├── gate/             # Quality gates (spec validation)
│   ├── review/           # Code review integration
│   ├── schema/           # Data schemas
│   ├── state/            # State machine
│   ├── testutil/         # Test utilities
│   └── validate/         # Validation logic
├── unit_test/            # Unit test configuration
└── artifacts/            # Output artifacts
```

## Quick Start

```bash
go run ./cmd/run/main.go
```

## Configuration

Unit test thresholds are configured in `unit_test/unit_test.yml`.
