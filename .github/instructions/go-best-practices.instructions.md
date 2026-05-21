---
applyTo: "**/*.go,go.mod,go.sum"
---

# Go best practices and coding standards

## Core principles

- Prefer the standard library before adding third-party dependencies.
- Write clear, idiomatic Go that follows `gofmt` formatting and common Go naming conventions.
- Keep implementations simple, explicit, and easy to reason about.
- Break code into small, focused pieces that each do one job well.
- Favor composition over deep abstractions or inheritance-style patterns.
- Follow SOLID principles where they improve clarity, maintainability, and changeability.

## Coding standards

- Keep packages small and focused around one responsibility.
- Use short, descriptive names; avoid abbreviations unless they are standard in Go.
- Apply proper whitespace consistently and follow idiomatic Go formatting, including spaces in control statements such as `if err != nil {` and `for i := 0; i < n; i++ {`.
- Accept interfaces where behavior is consumed, and return concrete types where practical.
- Pass `context.Context` as the first parameter for request-scoped or cancelable work.
- All errors MUST be handled explicitly; never ignore returned errors.
- Return wrapped errors with useful context when they cross boundaries.
- Avoid panics in application code except for unrecoverable startup or programmer errors.
- Keep functions small enough to read quickly; extract helpers when a function starts doing multiple jobs.
- Prefer small files, small functions, and small types over large multipurpose units.
- Write table-driven tests for logic-heavy behavior where it improves clarity and coverage.
- Add GoDoc comments for exported packages, types, functions, methods, and other exported identifiers.

## Project structure

- Organize code by domain or feature before technical layer when possible.
- Keep the application entry point thin in `cmd/<app-name>`.
- Put reusable internal application code in `internal/` to avoid exposing private implementation details.
- Put public libraries in `pkg/` only when the package is intentionally meant for external reuse.
- Keep transport, service, and storage concerns separated when the project grows, but avoid premature layering.
- Store configuration schemas, migrations, and other operational assets in clearly named top-level directories only when needed.

## Dependency and API design

- Keep module dependencies minimal and actively used.
- Prefer small interfaces that model the exact behavior required by the caller.
- Avoid creating interfaces only for testing; use concrete types unless an abstraction has production value.
- Design constructors that validate required dependencies up front.
- Apply SOLID pragmatically: single responsibility, extension without risky modification, substitutable implementations, narrow interfaces, and dependency inversion where it reduces coupling.

## Concurrency and reliability

- Use goroutines only when concurrency provides a clear benefit.
- Coordinate goroutines with contexts, channels, and `sync` primitives in straightforward ways.
- Ensure goroutines can stop cleanly and do not leak on cancellation or shutdown.
- Be explicit about ownership of shared state and protect it correctly.

## Logging and observability

- Produce structured, actionable log messages.
- Do not log secrets, tokens, or sensitive user data.
- Return errors to the caller when they need to decide how to handle failure; log at the application boundary.

## Testing and performance

- Comprehensive tests MUST be created for new and changed behavior, and those tests MUST pass.
- Cover success cases, failure cases, edge cases, and regression-prone paths.
- Add benchmarks for performance-critical code paths or when a change may affect runtime or allocation behavior.
- Use benchmark results to compare approaches and avoid speculative performance claims.

## Maintenance expectations

- Keep package documentation and exported identifiers clear when code is intended for reuse.
- Update tests and relevant documentation when behavior, public APIs, or setup steps change.
- Prefer incremental improvements that preserve readability over clever optimizations.
