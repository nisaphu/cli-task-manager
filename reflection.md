# Socio-Technical Reflection & Governance

## 1. Architectural Guardrails

ADR 0001 acted as the main architectural guardrail for persistence throughout the Spec-Driven Development workflow. The ADR established that task data must be persisted using local JSON files. This decision was then reinforced in the project constitution, implementation plan, and task breakdown.

During `/speckit.plan`, the persistence constraint was made explicit: the implementation had to use Go's standard library `encoding/json`, and SQLite, SQL databases, ORMs, embedded databases, and other persistence mechanisms were prohibited. The generated tasks also preserved this constraint by assigning persistence work specifically to the JSON storage layer and by including tests for JSON save/load behavior, missing files, and corrupted JSON.

During `/speckit.implement`, the AI agent did not attempt to deviate from the selected persistence strategy. The implementation used local JSON files in `internal/store/json.go` and used atomic file replacement with `os.CreateTemp` and `os.Rename`. No SQL or database dependency was introduced.

This demonstrates how the ADR constrained implementation before coding began. Instead of relying on the AI agent to independently choose a storage technology, the decision had already been made, documented, and propagated into the specification-driven workflow.

## 2. Closed-Loop Correction

The `/speckit.converge` step was effective as a final closed-loop governance check. It assessed the implementation against the active specification, implementation plan, completed task list, project constitution, and ADR 0001.

The convergence assessment reported:

- 17 of 17 functional requirements satisfied
- 5 of 5 success criteria satisfied
- 12 of 12 user-story acceptance scenarios satisfied
- 6 of 6 architecture-plan decisions satisfied
- 7 of 7 constitution principles satisfied
- ADR 0001 compliance verified
- 0 gaps found

It also searched the codebase for SQL/database usage and verified that persistence was implemented exclusively through local JSON files.

In this implementation, `/speckit.converge` did not identify architectural drift, so no corrective tasks were appended. This is still valuable because the command provided evidence that the final code converged with the original intent. If the implementation had introduced SQLite, omitted a required behavior, or violated a constitutional rule, the closed-loop process could have converted the discrepancy into additional work rather than allowing the implementation to be considered complete.

## 3. Human-in-the-Loop Governance

A tool such as Decision Guardian could assist a human reviewer by surfacing ADR 0001 directly when a Pull Request contains changes related to persistence or task storage.

For example, if a Pull Request modified `internal/store/` or introduced a database dependency, the reviewer could immediately see that ADR 0001 states that tasks must be persisted using local JSON files. This reduces the need for reviewers to manually search the repository for historical architectural decisions.

The human reviewer still remains responsible for judgment. The tool does not replace architectural review; instead, it provides the relevant decision context at the point where the reviewer needs it. This supports human-in-the-loop governance by helping the reviewer compare the proposed code change with the documented architectural intent before approving the Pull Request.

## 4. Automated Fitness Function

To treat the persistence decision as "Decisions as Code," an automated architecture test can be added to the Go test suite.

The proposed fitness function scans Go imports in the repository and fails if a forbidden alternative persistence technology is introduced. Examples include:

- `database/sql`
- SQLite drivers
- GORM
- Ent
- BoltDB
- BadgerDB

The test also verifies that the JSON storage implementation imports `encoding/json`.

Because this check is implemented as a Go test, it runs automatically with:

```bash
go test ./...
```

If a developer or AI coding agent introduces a forbidden persistence dependency, the test fails and therefore causes the build or CI quality gate to fail.

This provides an automated enforcement mechanism for ADR 0001 and helps prevent future architectural drift.

## Conclusion

This assignment demonstrated that architectural decisions can be made before implementation and then continuously governed throughout development. ADR 0001 captured the reason for selecting local JSON persistence, the constitution converted that decision into a project rule, Spec Kit propagated the rule through specification, planning, tasks, and implementation, and `/speckit.converge` verified that the final implementation remained aligned.

The result is a development process where source code is not the only source of architectural knowledge. Intent, decisions, specifications, and automated checks work together to preserve architectural integrity.
