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
  - `Implementation Checklist`
- Keep these notes concise, specific, and directly tied to the feature being built.
- Write acceptance criteria in clear, testable statements.
- Store feature notes in the `docs` folder.
- Use one file per feature so each feature has its own focused specification.
- Name feature files clearly in kebab-case, for example `docs/features/kanban-board.md`.
- Use markdown checkboxes in the user expectations so it is clear which expectations are already covered.
- Use markdown checkboxes in the acceptance criteria so it is clear which criteria have been met.
- Use markdown checkboxes in the implementation checklist so it is clear what has been completed.

## Writing guidance for feature notes

- `Project Status`: state whether the feature is proposed, planned, in progress, blocked, or complete.
- `Summary`: explain what the feature does and why it exists.
- `User Expectations`: describe what the user should be able to do or experience, using markdown checkboxes such as `- [x]` and `- [ ]`.
- `Acceptance Criteria`: define the observable conditions that must be true for the feature to be considered complete, using markdown checkboxes such as `- [x]` and `- [ ]`.
- `Implementation Checklist`: list the main deliverables with markdown checkboxes such as `- [x]` and `- [ ]`.

## Delivery expectations

- Implement features in a way that preserves clear functional boundaries.
- Update related tests and documentation so they stay aligned with the feature summary and acceptance criteria.
- When a feature is implemented, ensure its dedicated file in `docs` is added or updated.
- Keep the user expectations current so the document shows which expectations are already covered.
- Keep the acceptance criteria current so the document shows which completion conditions are already satisfied.
- Keep the implementation checklist current so the document shows what is already done and what is still pending.
