# Research & Architecture Decisions: CLI Task Manager & Focus Timer

**Feature**: CLI Task Manager & Focus Timer  
**Spec**: [`specs/001-task-manager-focus-timer/spec.md`](spec.md)  
**Date**: 2026-08-14  

---

## 1. JSON Storage Strategy & Atomic Persistence (ADR 0001 Compliance)

### Decision
Store task data using Go standard library `encoding/json` in a local JSON file (defaulting to `tasks.json` in the user's data directory or current working directory). Write operations will use atomic file replace (`os.CreateTemp` + `os.Rename`) to prevent file corruption in case of unexpected interrupts or process termination.

### Rationale
- **ADR 0001 Compliance**: Strictly adheres to the project constitution and ADR 0001 requiring local JSON persistence without external database drivers or ORMs.
- **Data Safety**: Direct writes to an existing file can lead to truncation or partial JSON corruptions if killed mid-write. Atomic write via temporary file replacement guarantees file integrity.
- **Zero External Dependencies**: Accomplished entirely using Go's built-in `os`, `io`, and `encoding/json` packages.

### Alternatives Considered
- Direct `os.WriteFile`: Simple, but vulnerable to zero-byte truncation if the process is forcefully killed (SIGINT/SIGTERM) mid-write.
- Embedded SQLite: Explicitly rejected by ADR 0001 and project constitution due to unnecessary setup and external driver overhead for single-user CLI.

---

## 2. CLI Command Routing & Standard Library Parsing

### Decision
Use a lightweight subcommand router built around Go's standard library `os.Args` and `flag.NewFlagSet`. Subcommands (`add`, `list`, `complete`, `focus`) are parsed cleanly with standard argument positioning.

### Rationale
- **Zero External Dependencies**: Avoids heavy third-party CLI frameworks (like Cobra or Urfave/CLI) keeping binary size under ~5MB and startup latency near instantaneous (~2ms).
- **Clear Error Guidance**: Custom flag sets per subcommand allow formatted help messages and exact usage syntax for invalid commands.

### Alternatives Considered
- Third-party CLI frameworks (Cobra/Viper): Rejected to comply with project principle I (Code Quality & Dependency Minimization).

---

## 3. Focus Timer Testability & Duration Injection

### Decision
Design `FocusSession` with an injectable duration parameter (`time.Duration`). Production default uses `25 * time.Minute`. Unit tests pass sub-second durations (e.g., `10 * time.Millisecond`) to verify start, ticker/sleep execution, and finish events synchronously without delaying test suites.

### Rationale
- **Fast Test Suite**: Tests run in milliseconds instead of waiting 25 real minutes.
- **Clean Separation**: Decouples time-stepping and duration parameters from the production CLI trigger.

### Alternatives Considered
- Real 25-minute sleep in tests: Unusable in continuous integration (adds 25 minutes to test run).
- Virtual clock mock interface: Unnecessary abstraction for a single-session timer; duration parameter injection provides identical testability with simpler code.

---

## 4. Task Identifier Generation & Sequence Strategy

### Decision
Use sequential auto-incrementing integer IDs (1, 2, 3, ...). The next available ID is calculated as `max(existing_ids) + 1` (or tracking an internal sequence counter in the task store file).

### Rationale
- **Developer Ergonomics**: Short integer IDs (`1`, `2`) are far easier to type in terminal commands (`taskman complete 1`) than UUIDs (`550e8400-e29b-41d4-a716-446655440000`).
- **Human Readability**: Matches standard developer expectations for CLI task managers.

### Alternatives Considered
- UUID strings: Unfriendly to CLI user experience.
- Random short hashes: Harder to type and remember than sequential integer IDs.
