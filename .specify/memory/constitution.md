<!--
Sync Impact Report:
- Version change: Unversioned Initial Scaffold → v1.0.0
- List of modified principles:
  - Added Principle I: Code Quality
  - Added Principle II: Architecture Decision Compliance
  - Added Principle III: Separation of Concerns
  - Added Principle IV: Testability
  - Added Principle V: Specification-Driven Development
  - Added Principle VI: Simplicity for a Local CLI Tool
  - Added Principle VII: Architectural Integrity
- Added sections:
  - Technical & Operational Constraints
  - Development Workflow & Quality Gates
- Removed sections: None
- Follow-up TODOs: None
-->

# CLI Task Manager & Focus Timer Constitution

## Core Principles

### I. Code Quality
- Code MUST be kept simple, readable, and maintainable.
- Implementation MUST prefer small, focused modules and functions adhering to single responsibility.
- Developers MUST avoid unnecessary complexity, over-engineering, and extraneous external dependencies.

*Rationale*: A clean, maintainable codebase ensures long-term readability, ease of debugging, and minimal cognitive overhead during development and maintenance.

### II. Architecture Decision Compliance
- All implementation decisions MUST follow documented Architecture Decision Records (ADRs).
- ADRs serve as authoritative architectural guardrails and MUST be respected unless superseded by a newer accepted ADR.
- ADR 0001 requires task persistence using local JSON files.
- Developers MUST NOT introduce SQLite, external databases, ORMs, or alternative persistence mechanisms unless a new ADR explicitly changes this decision.

*Rationale*: Strict compliance with accepted ADRs prevents architectural drift and ensures structural choices remain documented, deliberate, and aligned with project constraints.

### III. Separation of Concerns
- System components MUST maintain clear boundaries between task management logic, persistence logic, CLI interaction, and focus timer functionality.
- The persistence implementation MUST be decoupled and replaceable without requiring changes to core task-management business logic.

*Rationale*: Decoupling core business domain logic from I/O, CLI presentation, and persistence mechanisms simplifies independent testing, maintenance, and future interface adapters.

### IV. Testability
- Core business logic MUST be testable in isolation, completely independent from CLI input/output streams and file persistence mechanisms.
- All critical domain behaviors, state transitions, and edge cases MUST be covered by automated unit and integration tests.

*Rationale*: Automated unit tests for domain logic provide rapid feedback, guarantee behavioral correctness, and prevent regressions without relying on disk I/O or terminal interactions.

### V. Specification-Driven Development
- Specifications, implementation plans, and task breakdowns serve as the explicit source of truth for all implementation work.
- Implementation MUST NOT introduce scope creep, unapproved features, or behaviors not defined in the active specification.

*Rationale*: Strictly adhering to approved specifications ensures predictable project execution, eliminates wasted effort on unvetted features, and aligns delivery with requirements.

### VI. Simplicity for a Local CLI Tool
- System design MUST optimize specifically for a single-user, local command-line interface environment.
- Developers MUST avoid premature scalability, distributed systems patterns, complex caching layers, or unnecessary infrastructure complexity.

*Rationale*: Premature optimization and complex infrastructure introduce operational friction without adding value for a local single-user developer tool.

### VII. Architectural Integrity
- Any implementation that conflicts with an accepted ADR, constitution principle, or specification MUST be corrected before task completion or code merge.
- Continuous architectural review is mandatory throughout the development lifecycle.

*Rationale*: Enforcing strict completion gates prevents technical debt accumulation and guarantees long-term compliance with governance standards.

## Technical & Operational Constraints

- **Execution Environment**: Operates entirely as a local single-user CLI application.
- **Persistence Strategy**: Task data is serialized and stored using local JSON files as defined in ADR 0001.
- **I/O Protocols**: Interacts with the terminal using standard I/O conventions (stdin/arguments for input, stdout for formatted output, stderr for errors).
- **Dependency Minimization**: Prefers standard library utilities and lightweight dependencies to keep the installation light and fast.

## Development Workflow & Quality Gates

- **Specification Gate**: Feature development MUST NOT begin without an approved feature specification and task breakdown.
- **Test Gate**: All existing and new automated tests MUST pass cleanly before code changes are merged or finalized.
- **ADR Audit Gate**: Any change impacting system architecture or persistence MUST be verified against active ADRs before completion.

## Governance

- **Amendment Procedure**: Proposed amendments to this constitution MUST be documented, submitted with explicit rationales, and formally accepted before taking effect.
- **Versioning Policy**: Semantic versioning (`MAJOR.MINOR.PATCH`) governs this constitution:
  - `MAJOR`: Backward-incompatible principle removals or fundamental governance shifts.
  - `MINOR`: Addition of new core principles or major structural guidance sections.
  - `PATCH`: Clarifications, typo fixes, and non-semantic formatting updates.
- **Compliance Enforcement**: Every pull request, task implementation, and code audit MUST be evaluated against these principles. Non-compliant contributions MUST be resolved prior to task completion.

**Version**: 1.0.0 | **Ratified**: 2026-08-13 | **Last Amended**: 2026-08-13
