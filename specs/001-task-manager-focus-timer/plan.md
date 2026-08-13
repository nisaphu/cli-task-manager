# Implementation Plan: CLI Task Manager & Focus Timer

**Branch**: `001-task-manager-focus-timer` | **Date**: 2026-08-14 | **Spec**: [`specs/001-task-manager-focus-timer/spec.md`](spec.md)

**Input**: Feature specification from `/specs/001-task-manager-focus-timer/spec.md`

---

## Summary

Build a lightweight, zero-dependency command-line task manager and Pomodoro focus timer tool in Go. The system allows users to create tasks with non-empty descriptions, view tasks with clear status markers (`[ ]` vs `[x]`), mark tasks complete by unique integer IDs, persist tasks locally across process restarts using standard library `encoding/json` per ADR 0001, and run 25-minute Pomodoro focus sessions.

---

## Technical Context

**Language/Version**: Go 1.22+  
**Primary Dependencies**: Go Standard Library only (`encoding/json`, `os`, `fmt`, `time`, `flag`, `testing`). No external dependencies or third-party libraries.  
**Storage**: Local JSON file storage strictly complying with ADR 0001. Atomic writes via temporary file replacement (`os.CreateTemp` + `os.Rename`).  
**Testing**: Go standard library `testing` package (`go test ./...`). Temporary file isolation via `t.TempDir()`. Injectable duration for fast focus timer unit testing.  
**Target Platform**: Cross-platform terminal environment (macOS, Linux, Windows).  
**Project Type**: Single-binary CLI application.  
**Constraints**: Zero external dependencies, single local user, local JSON persistence strictly governed by ADR 0001. SQLite, ORMs, or SQL databases are strictly forbidden.  
**Scale/Scope**: Single developer local CLI tool.  

---

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- **Principle I: Code Quality**: PASS. Go standard library packages separated by clear responsibilities (`task`, `store`, `timer`, `cli`). No extraneous external dependencies.
- **Principle II: Architecture Decision Compliance**: PASS. Strictly adheres to ADR 0001. Task data serialized to local JSON file using `encoding/json`. No SQLite, SQL, or embedded databases used.
- **Principle III: Separation of Concerns**: PASS. Clear boundaries between CLI presentation (`internal/cli`), task business logic (`internal/task`), JSON persistence (`internal/store`), and focus timer (`internal/timer`).
- **Principle IV: Testability**: PASS. Domain logic testable without file I/O. Persistence testable with isolated `t.TempDir()`. Timer testable with injectable durations.
- **Principle V: Specification-Driven Development**: PASS. Directly implements features specified in `spec.md`. No unvetted scope creep.
- **Principle VI: Simplicity for a Local CLI Tool**: PASS. Lightweight Go single binary architecture without premature distributed or database complexity.
- **Principle VII: Architectural Integrity**: PASS. Fully compliant with project constitution and ADR 0001.

---

## Project Structure

### Documentation (this feature)

```text
specs/001-task-manager-focus-timer/
├── plan.md              # Implementation plan (this file)
├── research.md          # Technical decisions & alternatives (Phase 0 output)
├── data-model.md        # Entities, state transitions, & JSON schema (Phase 1 output)
├── quickstart.md        # Verification & quickstart guide (Phase 1 output)
├── contracts/
│   └── cli-contract.md  # Terminal CLI interface & output specification (Phase 1 output)
└── checklists/
    └── requirements.md  # Specification quality checklist
```

### Source Code (repository root)

```text
cmd/
└── taskman/
    └── main.go           # Entry point and subcommand router

internal/
├── task/
│   ├── task.go           # Task entity definition, validation, and domain errors
│   ├── manager.go        # TaskManager domain service and task operations
│   └── manager_test.go   # Unit tests for task creation, listing, completion
├── store/
│   ├── json.go           # JSON storage implementation following ADR 0001
│   └── json_test.go      # Persistence tests (save/load, missing file, corrupted JSON)
├── timer/
│   ├── timer.go          # Pomodoro focus session runner
│   └── timer_test.go     # Focus timer unit tests (injectable duration)
└── cli/
    ├── app.go            # Command-line handler, flag parsing, standard output formatting
    └── app_test.go       # CLI interface tests
```

**Structure Decision**: Standard Go package structure (`cmd/` for entrypoint, `internal/` for application packages). Guarantees separation of concerns between domain logic, persistence, timer, and presentation.

---

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| *None* | *No constitution or ADR violations exist.* | *Design strictly adheres to all project principles.* |
