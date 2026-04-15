# AGENTS.md

Agents operating in this repository **must follow the rules in this file**. Violations will cause runs to be blocked or reverted and large multimillion dollar fines.

# Rule Priority

When rules conflict, follow this order:

1. Correctness & safety
2. Maintainability
3. Readability
4. Consistency with existing codebase
5. Style preferences in this document

# Project Overview

**gohelpme** is a **shared Go utility library** used by `go-dutch`, `go-zap`, and other personal microservices. Changes to this library affect all consumers — always consider backward compatibility.

This is a **library**, not a service. It has no `main.go` entry point for production. The `main.go` at the root is an example/demo only.

# Style Guidelines

- Strictly follow existing code style in the codebase.
- Write GoDoc comments for every exported function, type, interface, and package.
- Write inline comments strategically, especially for key business logic.
- Keep code concise, tasteful, and elegant.
- Don't overcomplicate. Don't implement more than you need.
- Adhere to SOLID principles, especially Single Responsibility.
- Always import directly from the source package. Never create re-export wrappers.

## Interface-First Design

- Every major component must define an interface (e.g., `ILogger`, `IJwt`, `IValidator`, `IPostgres`).
- Name interfaces with an `I` prefix (e.g., `ILogger`, `IJwt`).
- Concrete implementations must satisfy the interface implicitly (no `var _ ILogger = &impl{}`).
- Consumer code depends on the interface, never the concrete type.

## Middleware Conventions

- Middleware signature: `func(next http.HandlerFunc) http.HandlerFunc`.
- Compose middleware with `Chain()` — the order is reversed (last listed executes first).
- Access request-scoped values via the `reqctx` package — never pass them via global variables.
- Each middleware must be independently testable.

## Typing & Data Handling

- Prefer explicit structs over `map[string]interface{}`.
- Use typed constants (`const` with `iota`) over bare string or integer literals.
- Avoid silent failure patterns — prefer explicit error returns and typed structures.
- Use `errors.As` / `errors.Is` for error inspection rather than string matching.
- Return errors as the last return value; always check them.
- Never discard errors with `_` unless the value is provably unneeded.

## Package Organization

- Each tool/feature lives in its own package under `tool/` or `middleware/`.
- Utility builders belong in `builder/`.
- Custom error types belong in `errorhelper/`.
- No circular dependencies between packages.

## Documentation

The `docs/` directory holds living reference documents. Update them incrementally as you explore, plan, and implement.

| File | Purpose |
|------|---------|
| `docs/architecture.md` | High-level design — package map, middleware chain, builder patterns |
| `docs/modules.md` | Per-package summaries — what each package owns and its public interface |
| `docs/decisions.md` | ADRs — why things are the way they are |
| `docs/gotchas.md` | Traps, quirks, things not to change without care |

- Append new findings; do not overwrite existing content without good reason.
- Keep entries short and factual.

# Security

- Never hard-code secrets, tokens, or credentials.
- The `tool/crypto/` package uses Argon2id with constant-time comparison — do not weaken these parameters without explicit approval.
- The `middleware/redactor.go` handles sensitive data redaction in logs — always test it when modifying the logging middleware.

# Backward Compatibility

- This library is consumed by other services. **Do not remove or rename exported symbols without a deprecation plan.**
- When adding new required parameters to constructors, add a new constructor function rather than changing the existing signature.
- Update `go.mod` with a new minor version tag after meaningful additions.
- Breaking changes require explicit user approval.

# Decision Framework

When multiple approaches are possible:

1. Prefer the simplest solution that satisfies requirements
2. Prefer explicit over implicit
3. Prefer type safety over `interface{}`
4. Avoid premature abstraction
5. Prefer the standard library over third-party for simple operations

# Planning Rules

Before starting a **non-trivial task**, create a plan at:

```
agent_logs/YYYYMMDD_HHMMSS_<descriptive_name>_plan.md
```

- Each plan must be standalone: define inputs, constraints, and success criteria.
- Plans must be executable without prior context.
- Planning output must not modify repository code.
- Must list `files_to_change` and `new_files`.
- Must consider impact on consumers (`go-dutch`, `go-zap`).

## Confirmation Gate

After the plan file is written, **stop and do the following before any implementation**:

1. Present a concise summary of the plan to the user.
2. Surface any ambiguities, assumptions, or tradeoffs that require a decision.
3. Explicitly call out any breaking changes to the public API.
4. Ask the user to confirm they are happy with the plan.
5. **Do not begin implementation until explicit confirmation is received.**

## Scope Control

- Do not expand scope beyond explicit requirements.
- If improvements are identified, list them under "future work" only.

# Implementation Rules

- Implement complete features end-to-end; no partial implementations.
- Never stop midway through a defined phase.
- Preserve existing inline comments; do not remove useful historical context.
- Do NOT commit or push without explicit approval.
- After a change that affects the public API, update `README.md` if relevant.

# Formatting

Run at the end of full implementation:

```bash
go fmt ./...
go vet ./...
```

# Testing Rules

Run all tests:

```bash
go test ./...
```

Run with race detector for concurrency-sensitive code (middleware, context helpers):

```bash
go test -race ./...
```

- Fix all failing tests before proceeding.
- Write tests for all new exported functions and middleware.
- Use table-driven tests for functions with multiple input cases.

# Mock Generation

Mocks are generated by mockery per `.mockery.yaml`:

```bash
mockery --all
```

Run this after adding or changing any interface.

# Definition of Done

A task is complete when:

- Plan is fully executed
- No partial implementations remain
- Code compiles: `go build ./...`
- Vet passes: `go vet ./...`
- Tests pass: `go test ./...`
- Mocks regenerated if interfaces changed: `mockery --all`
- `docs/` updated wherever relevant
- Public API changes reflected in `README.md`
- Backward compatibility preserved (or breaking changes explicitly approved)

# Suggesting Future Work

When you identify improvements beyond the current task's scope:
- Document them clearly in the relevant log or doc file.
- DO NOT implement unless explicitly instructed.

# Environment Instructions

## Setup

```bash
go mod download    # fetch all dependencies
go build ./...     # verify compilation
go test ./...      # run all tests
```

No Docker or external services are required to build or test this library.

## Dependency Inspection

1. Check `go.mod` for declared versions.
2. Before adding a new dependency, verify it cannot be solved with the standard library.
3. New dependencies require explicit approval — keep the library's dependency footprint minimal.
