---
applyTo: "**/*"
---

# Feature organization and planning guidance

## Feature grouping

- Group features by business functionality or user-facing capability, not by technical layer alone.
- Keep all code, handlers, services, models, tests, and supporting files for the same feature close together when the project structure allows it.
- Prefer adding to an existing functional area over scattering related changes across unrelated directories.
- Name folders and modules so the feature purpose is immediately clear.

## Feature design expectations

- Break larger work into small, focused feature pieces that can be implemented and reviewed independently.
- Keep each feature slice responsible for a clear outcome.
- Avoid mixing unrelated functionality into the same module, package, or change set.

## Required feature notes

- When introducing or planning a feature, include a short structured note with:
  - `Project Status`
  - `Summary`
  - `User Expectations`
  - `Acceptance Criteria`
  - `License Check`
- Keep these notes concise, specific, and directly tied to the feature being built.
- Write acceptance criteria in clear, testable statements.

## Writing guidance for feature notes

- `Project Status`: state whether the feature is proposed, planned, in progress, blocked, or complete.
- `Summary`: explain what the feature does and why it exists.
- `User Expectations`: describe what the user should be able to do or experience.
- `Acceptance Criteria`: define the observable conditions that must be true for the feature to be considered complete.
- `License Check`: record the licenses of new dependencies or assets and confirm they are acceptable for the project.

## License expectations

- Check the license of every new dependency, library, snippet source, and bundled asset before using it.
- MIT and similar permissive licenses can be used when they are compatible with the project.
- Prefer permissive licenses such as MIT, BSD, ISC, and Apache-2.0 unless told otherwise.
- Do not add dependencies with unclear, missing, or incompatible licenses.

## Delivery expectations

- Implement features in a way that preserves clear functional boundaries.
- Update related tests and documentation so they stay aligned with the feature summary and acceptance criteria.
