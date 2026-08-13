# Tasks: CLI Task Manager & Focus Timer

**Feature**: CLI Task Manager & Focus Timer  
**Spec**: [`specs/001-task-manager-focus-timer/spec.md`](spec.md)  
**Plan**: [`specs/001-task-manager-focus-timer/plan.md`](plan.md)  

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Initialize Go module and directory layout according to the implementation plan.

- [x] T001 Initialize Go module `cli-task-manager` and directory structure (`cmd/taskman/`, `internal/task/`, `internal/store/`, `internal/timer/`, `internal/cli/`) in `go.mod`
- [x] T002 [P] Configure Go test runner and helper utilities in `internal/task/test_helpers.go`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core domain types, persistence interface, and domain errors required before implementing any user story.

**⚠️ CRITICAL**: No user story work can begin until this phase is complete.

- [x] T003 Implement `Task` struct, `TaskStatus` enum, and domain error definitions (`ErrEmptyDescription`, `ErrTaskNotFound`, `ErrTaskAlreadyCompleted`, `ErrCorruptedStorage`) in `internal/task/task.go`
- [x] T004 Define `Store` persistence interface and `TaskData` container schema in `internal/store/json.go`
- [x] T005 [P] Add unit tests for `Task` entity initialization and validation in `internal/task/task_test.go`

**Checkpoint**: Foundation ready - user story implementation can now begin.

---

## Phase 3: User Story 1 - Add and View Tasks (Priority: P1) 🎯 MVP

**Goal**: Allow users to create tasks with non-empty descriptions, assign unique integer IDs, and view task list with incomplete markers `[ ]`.

**Independent Test**: Add tasks via CLI, view list to verify ID, description, `[ ]` incomplete status, and confirm clear message when no tasks exist.

### Tests for User Story 1

- [x] T006 [P] [US1] Unit tests for `TaskManager.AddTask` with description validation and ID assignment in `internal/task/manager_test.go`
- [x] T007 [P] [US1] Unit tests for task list formatting (`[ ]` status tag) and empty list notice in `internal/cli/app_test.go`

### Implementation for User Story 1

- [x] T008 [US1] Implement `TaskManager.AddTask` and `TaskManager.ListTasks` domain logic in `internal/task/manager.go`
- [x] T009 [US1] Implement CLI handlers for `add` and `list` subcommands in `internal/cli/app.go`

**Checkpoint**: At this point, User Story 1 (MVP) is fully functional and testable independently.

---

## Phase 4: User Story 2 - Complete Tasks and State Persistence (Priority: P2)

**Goal**: Mark tasks complete by ID, safely handle already-completed or non-existent IDs, and persist tasks locally to JSON per ADR 0001 across CLI restarts.

**Independent Test**: Add a task, complete task by ID, restart CLI process, verify completed `[x]` status remains intact, and verify error responses for non-existent IDs or corrupted JSON files.

### Tests for User Story 2

- [x] T010 [P] [US2] Unit tests for JSON persistence saving, reloading, missing file handling, and corrupted JSON handling using `t.TempDir()` in `internal/store/json_test.go`
- [x] T011 [P] [US2] Unit tests for `TaskManager.CompleteTask`, non-existent ID error, and idempotent completion in `internal/task/manager_test.go`

### Implementation for User Story 2

- [x] T012 [US2] Implement JSON storage loading, saving, and atomic file replace (`os.CreateTemp` + `os.Rename`) in `internal/store/json.go`
- [x] T013 [US2] Implement `TaskManager.CompleteTask` domain business logic in `internal/task/manager.go`
- [x] T014 [US2] Integrate `Store` persistence into `TaskManager` lifecycle in `internal/task/manager.go`
- [x] T015 [US2] Implement CLI handler for `complete` subcommand in `internal/cli/app.go`

**Checkpoint**: At this point, User Stories 1 AND 2 work independently with full persistence.

---

## Phase 5: User Story 3 - Pomodoro Focus Timer (Priority: P3)

**Goal**: Start a 25-minute Pomodoro focus session from CLI with start notification, completion notification, and terminal cancellation handling.

**Independent Test**: Execute `focus` command, observe start notification, test session execution using injectable duration without waiting 25 real minutes, and verify Ctrl+C cancellation.

### Tests for User Story 3

- [x] T016 [P] [US3] Unit tests for `FocusSession` execution using injectable sub-second durations in `internal/timer/timer_test.go`

### Implementation for User Story 3

- [x] T017 [US3] Implement `FocusSession` component with Go `time` utilities and 25-minute production default in `internal/timer/timer.go`
- [x] T018 [US3] Implement CLI handler for `focus` subcommand and terminal signal cancellation (SIGINT/Ctrl+C) in `internal/cli/app.go`

**Checkpoint**: User Story 3 is fully operational and testable.

---

## Phase 6: User Story 4 - CLI Usability and Error Guidance (Priority: P4)

**Goal**: Provide a clean subcommand router (`add`, `list`, `complete`, `focus`, `help`) with actionable error messages for unknown commands or missing parameters.

**Independent Test**: Execute unrecognized subcommands or omit required parameters and verify exit code 1 with clear usage instructions.

### Tests for User Story 4

- [x] T019 [P] [US4] Unit tests for CLI subcommand router, unknown command error formatting, and help menu output in `internal/cli/app_test.go`

### Implementation for User Story 4

- [x] T020 [US4] Implement CLI subcommand router, argument parsing, usage help guide, and entrypoint in `cmd/taskman/main.go`

---

## Phase 7: Polish & Cross-Cutting Concerns

**Purpose**: Integration verification, documentation audit, and compliance checks.

- [x] T021 [P] Integration tests for complete CLI task management and persistence workflow in `internal/cli/app_test.go`
- [x] T022 [P] Perform quickstart validation scenarios and verify strict compliance with ADR 0001 and project constitution in `specs/001-task-manager-focus-timer/quickstart.md`

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately.
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories.
- **User Stories (Phases 3 - 6)**: All depend on Foundational phase completion.
  - User Stories proceed in priority order (P1 → P2 → P3 → P4) or in parallel where files do not conflict.
- **Polish (Phase 7)**: Depends on all user story phases being complete.

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational (Phase 2). No dependencies on other stories.
- **User Story 2 (P2)**: Can start after Foundational (Phase 2). Integrates persistence with TaskManager.
- **User Story 3 (P3)**: Can start after Foundational (Phase 2). Independent timer module.
- **User Story 4 (P4)**: Can start after Foundational (Phase 2). Connects subcommand routing across all handlers.

---

## Parallel Execution Opportunities

- All tasks marked `[P]` can run in parallel within their respective phases:
  - Setup: `T002`
  - Foundational: `T005`
  - User Story 1: `T006`, `T007`
  - User Story 2: `T010`, `T011`
  - User Story 3: `T016`
  - User Story 4: `T019`
  - Polish: `T021`, `T022`

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational
3. Complete Phase 3: User Story 1 (Add and View Tasks)
4. **VALIDATE**: Verify task creation and listing independently via unit tests.

### Incremental Delivery

1. Setup + Foundational → Core domain ready
2. Add User Story 1 → Test independently → MVP!
3. Add User Story 2 → Add JSON persistence & completion → Complete task storage lifecycle
4. Add User Story 3 → Add Pomodoro focus timer
5. Add User Story 4 → Finalize CLI subcommand routing & polish
